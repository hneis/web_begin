// Package server provides ...
package server

import (
	"github.com/go-chi/chi"
)

func (serv *Server) bindRoutes(r *chi.Mux) {
	r.Route("/", func(r chi.Router) {
		r.Get("/", serv.getTemplateHandler)
		r.Get("/post/{id}", serv.getPostHandler)
		r.Get("/post/{id}/edit", serv.getPostEditHandler)
		r.Route("/api/v1", func(r chi.Router) {
			r.Post("/post", serv.postPostHandler)
			r.Delete("/post/{id}", serv.deletePostHandler)
			r.Put("/post/{id}", serv.putPostHandler)
		})
	})
}
