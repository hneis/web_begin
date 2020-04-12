package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	urls = []string{
		"https://ru.wikipedia.org/wiki/%D0%9F%D0%B5%D0%BD%D0%B5%D0%BB%D0%BE%D0%BF%D0%B0",
		"https://www.kinopoisk.ru/film/189651/",
		"https://www.cisco.com/c/ru_ru/index.html",
		"https://ast.ru/authors/duglas-penelopa-986618/",
		"https://www.velosklad.ru/",
	}
)

func findContent(pattern string, urls []string) {
	for _, url := range urls {
		conn, err := http.Get(url)
		defer conn.Body.Close()

		if err != nil {
			log.Println(err.Error())
		}

		data, err := ioutil.ReadAll(conn.Body)
		if err != nil {
			log.Println(err.Error())
		}

		if bytes.Contains(data, []byte(pattern)) {
			fmt.Println(url)
		}
	}
}

func main() {
	pattern := "Пенелопа"
	findContent(pattern, urls)
}
