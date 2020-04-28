// Package server provides ...
package server

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/go-chi/chi"
	"github.com/hneis/web_begin/lesson5/homework/models"
	uuid "github.com/satori/go.uuid"
)

// getTemplateHandler - возвращает шаблон
func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")

	if templateName == "" {
		templateName = serv.indexTemplate
	}

	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, templateName))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("Page").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	posts, err := models.GetAllPostItems(serv.db)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	serv.Page.Posts = posts

	if err := templ.Execute(w, serv.Page); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func (serv *Server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, "/blogs/post.html"))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("post").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post := models.PostItem{ID: postID}
	if err := post.Get(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if err := templ.Execute(w, post); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func (serv *Server) getPostEditHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	file, err := os.Open(path.Join(serv.rootDir, serv.templatesDir, "/blogs/post_edit.html"))
	if err != nil {
		if err == os.ErrNotExist {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		serv.SendInternalErr(w, err)
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	templ, err := template.New("post").Parse(string(data))
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post := models.PostItem{ID: postID}
	if err := post.Get(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if err := templ.Execute(w, post); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func (serv *Server) postPostHandler(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)

	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)

	post.ID = uuid.NewV4().String()
	post.Created = time.Now().String()

	if err := post.Insert(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ = json.Marshal(post)
	w.Write(data)
}

func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	post := models.PostItem{ID: postID}
	if err := post.Delete(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func (serv *Server) putPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	data, _ := ioutil.ReadAll(r.Body)

	serv.lg.Println(data)
	post := models.PostItem{}
	_ = json.Unmarshal(data, &post)
	post.ID = postID
	serv.lg.Println(post)

	if err := post.Update(serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
