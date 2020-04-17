package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)

type data struct {
	Url   string
	Body  []byte
	Error error
}

func SearchInSites(pattern string, urls []string) ([]string, error) {
	chanels := []chanData{}

	for _, u := range urls {
		_url, err := url.ParseRequestURI(u)
		if err != nil {
			return []string{}, err
		}
		chanels = append(chanels, getContent(_url.String()))
	}

	result := []string{}
	for data := range merge(chanels) {
		if bytes.Contains(data.Body, []byte(pattern)) {
			result = append(result, data.Url)
		}
	}

	return result, nil
}

func getContent(url string) <-chan data {
	c := make(chan data)

	go func() {
		defer close(c)
		var body []byte
		var err error

		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)

		c <- data{Body: body, Error: err, Url: url}
	}()

	return c
}

type chanData <-chan data

func merge(cs []chanData) <-chan data {
	var wg sync.WaitGroup
	out := make(chan data)

	output := func(c <-chan data) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
