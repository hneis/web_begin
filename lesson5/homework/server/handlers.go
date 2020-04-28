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
	"github.com/gofrs/uuid"
	"github.com/hneis/web_begin/lesson5/homework/blog/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
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

	posts, err := models.Posts(qm.Load(models.PostRels.Author)).All(serv.ctx, serv.db)
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

	post, err := models.Posts(qm.Load(models.PostRels.Author), models.PostWhere.ID.EQ(postID)).One(serv.ctx, serv.db)
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

	post, err := models.Posts(qm.Load(models.PostRels.Author), models.PostWhere.ID.EQ(postID)).One(serv.ctx, serv.db)
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

	author, err := models.Users(models.UserWhere.Name.EQ(dataPost.Author)).One(serv.ctx, serv.db)
	_ = err
	serv.lg.Println(author)
	insert := false
	if author == nil {
		insert = true
		author = &models.User{
			ID:   uuid.Must(uuid.NewV4()).String(),
			Name: dataPost.Author,
			Hoby: "",
		}
		if err := author.Insert(serv.ctx, serv.db, boil.Infer()); err != nil {
			serv.SendInternalErr(w, err)
			return
		}
	}

	post := models.Post{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Title:   dataPost.Title,
		Content: dataPost.Text,
		Created: time.Now().String(),
	}

	err = post.SetAuthor(serv.ctx, serv.db, insert, author)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	if err := post.Insert(serv.ctx, serv.db, boil.Infer()); err != nil {
		serv.SendInternalErr(w, err)
		return
	}

	data, _ = json.Marshal(post)
	w.Write(data)
}

func (serv *Server) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	postID := chi.URLParam(r, "id")

	post, _ := models.FindPost(serv.ctx, serv.db, postID)
	rowsAff, err := post.Delete(serv.ctx, serv.db)
	_ = rowsAff

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

	author, err := models.Users(models.UserWhere.Name.EQ(dataPost.Author)).One(serv.ctx, serv.db)
	_ = err
	serv.lg.Println(author)
	insert := false
	if author == nil {
		insert = true
		serv.lg.Println("Add new author")
		author = &models.User{
			ID:   uuid.Must(uuid.NewV4()).String(),
			Name: dataPost.Author,
			Hoby: "",
		}
		if err := author.Insert(serv.ctx, serv.db, boil.Infer()); err != nil {
			serv.SendInternalErr(w, err)
			return
		}
	}

	post := models.Post{
		ID:      postID,
		Title:   dataPost.Title,
		Content: dataPost.Text,
	}

	err = post.SetAuthor(serv.ctx, serv.db, insert, author)
	if err != nil {
		serv.SendInternalErr(w, err)
		return
	}
	serv.lg.Println(post)

	if rowAff, err := post.Update(serv.ctx, serv.db, boil.Infer()); err != nil {
		_ = rowAff
		serv.SendInternalErr(w, err)
		return
	}
}
