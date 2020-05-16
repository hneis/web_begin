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
	"github.com/hneis/web_begin/lesson8/homework/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getTemplateHandler - возвращает шаблон
func (serv *Server) getTemplateHandler(w http.ResponseWriter, r *http.Request) {
	templateName := chi.URLParam(r, "template")

	if templateName == "" {
		templateName = serv.conf.IndexTemplate
	}

	file, err := os.Open(path.Join(serv.conf.RootDir, serv.conf.TemplatesDir, templateName))
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

	posts, err := models.GetPosts(serv.ctx, serv.db)
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

	file, err := os.Open(path.Join(serv.rootDir, serv.conf.TemplatesDir, "/blogs/post.html"))
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

	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post, err := models.GetPost(serv.ctx, serv.db, objectID)
	if err != nil {
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
	serv.lg.Infoln("edit")

	file, err := os.Open(path.Join(serv.rootDir, serv.conf.TemplatesDir, "/blogs/post_edit.html"))
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
	serv.lg.Info(postID)
	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post, err := models.GetPost(serv.ctx, serv.db, objectID)
	serv.lg.Infoln(post)
	if err != nil {
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

	dataPost := struct {
		Author string
		Title  string
		Text   string
	}{}
	_ = json.Unmarshal(data, &dataPost)

	post := models.Post{
		Title:   dataPost.Title,
		Content: dataPost.Text,
		Created: time.Now().String(),
		Author: models.Author{
			Name: dataPost.Author,
		},
	}

	if err := post.Insert(serv.ctx, serv.db); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ = json.Marshal(post)
	w.Write(data)
}

func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post := models.Post{
		Mongo: models.Mongo{
			ID: objectID,
		},
	}

	err = post.Delete(serv.ctx, serv.db)

	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}

func (serv *Server) putPostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	data, _ := ioutil.ReadAll(r.Body)

	serv.lg.Println(data)
	dataPost := struct {
		ID     string
		Author string
		Title  string
		Text   string
	}{}
	_ = json.Unmarshal(data, &dataPost)

	objectID, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post, err := models.GetPost(serv.ctx, serv.db, objectID)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	post.Title = dataPost.Title
	post.Content = dataPost.Text
	post.Author.Name = dataPost.Author

	err = post.Update(serv.ctx, serv.db)

	serv.lg.Println(post)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}
}
