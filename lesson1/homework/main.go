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
		"https://www.cisco.com/c/ru_ru/index.html",
		"https://ast.ru/authors/duglas-penelopa-986618/",
		"https://www.velosklad.ru/",
	}
)

func findContent(pattern string, urls []string) (result []string) {
	for _, url := range urls {
		conn, err := http.Get(url)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		defer conn.Body.Close()

		data, err := ioutil.ReadAll(conn.Body)
		if err != nil {
			log.Println(err.Error())
		}

		if bytes.Contains(data, []byte(pattern)) {
			result = append(result, url)
		}
	}

	return
}

func main() {
	pattern := "Пенелопа"
	fmt.Println(findContent(pattern, urls))
}
