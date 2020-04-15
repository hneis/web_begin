package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

/* usage
go run main.go -url="https://ru.wikipedia.org/wiki/%D0%9F%D0%B5%D0%BD%D0%B5%D0%BB%D0%BE%D0%BF%D0%B0","https://www.cisco.com/c/ru_ru/index.html","https://ast.ru/authors/duglas-penelopa-986618/","https://www.velosklad.ru/"
*/

type urls []string

func (u urls) String() string {
	return strings.Join(u, ",")
}

func (u *urls) Set(value string) error {
	for _, url := range strings.Split(value, ",") {
		*u = append(*u, url)
	}

	return nil
}

type data struct {
	Url   string
	Body  []byte
	Error error
}

func findInContent(pattern string, urls []string) []string {
	chanels := []chanData{}

	for _, url := range urlsFlag {
		chanels = append(chanels, getContent(url))
	}

	result := []string{}
	for data := range merge(chanels) {
		if bytes.Contains(data.Body, []byte(pattern)) {
			result = append(result, data.Url)
		}
	}

	return result
}

func getContent(url string) <-chan data {
	c := make(chan data)

	go func() {
		var body []byte
		var err error

		resp, err := http.Get(url)
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)

		c <- data{Body: body, Error: err, Url: url}

		close(c)
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

var urlsFlag urls
var pattern string

func init() {
	flag.Var(&urlsFlag, "url", "comma-separeated list of urls to use search pattern")
	flag.StringVar(&pattern, "pattern", "", "text to search in web content")
}

func main() {
	flag.Parse()
	if len(urlsFlag) == 0 {
		flag.Usage()
	}

	if pattern == "" {
		flag.Usage()
	}

	fmt.Println(findInContent(pattern, urlsFlag))
}
