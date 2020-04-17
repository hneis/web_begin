// Package main provides ...
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

type Query struct {
	Search string
	Sites  []string
}

func main() {
	stopchan := make(chan os.Signal)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Get("/search", SearchHnadler)
	})

	logrus.SetReportCaller(true)

	go func() {
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
}

func SearchHnadler(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}

	var query Query

	if err := json.Unmarshal(data, &query); err != nil {
		fmt.Fprintln(w, err.Error())
	}

	feature := make(chan bool)

	go func() {
		query.Sites, err = SearchInSites(query.Search, query.Sites)
		if err != nil {
			logrus.Error(err)
			// Добавить ошибку в json?
		}

		feature <- true
	}()

	<-feature

	buff, err := json.Marshal(query)
	if err != nil {
		fmt.Fprintln(w, err.Error())
		logrus.Error(err)
		// Добавить ошибку в json?
		// return если не смогли получить json?
	}

	w.Write(buff)
}
