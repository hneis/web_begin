// Package main provides ...
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

const (
	COOKIE_KEY = "cookie"
)

func main() {
	stopchan := make(chan os.Signal)

	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		router.Put("/cookie/{key}/{value}", PutCookieHandler)
		router.Get("/cookie/{key}", GetCookieHandler)
	})

	logrus.SetReportCaller(true)

	go func() {
		err := http.ListenAndServe(":8080", router)
		log.Fatal(err)
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
}

// PutCookieHandler TODO: NEEDS COMMENT INFO
func PutCookieHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	value := chi.URLParam(r, "value")

	cookie, _ := r.Cookie(key)
	if cookie == nil {
		cookie = &http.Cookie{
			Name: key,
			Path: "/",
		}
	}

	cookie.Value = value
	http.SetCookie(w, cookie)
}

// GetCookieHandler TODO: NEEDS COMMENT INFO
func GetCookieHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	cookie, _ := r.Cookie(key)
	if cookie != nil {
		w.Write([]byte(cookie.Value))
	} else {
		w.Write([]byte(fmt.Sprintf("Cookie \"%s\" not set", key)))
	}
}
