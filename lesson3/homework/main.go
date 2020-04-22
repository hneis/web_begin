package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/go-chi/chi"
)

type Server struct {
	lg    *logrus.Logger
	Title string
	Posts PostItems
}

type PostItems []PostItem
type PostItem struct {
	Id      int64
	Title   string
	Body    string // Добавить short and long text?
	Created string // переделать на time.Time
	Author  Author
}
type Author struct {
	Avatar string
	Name   string
}

func main() {
	stopchan := make(chan os.Signal)

	r := chi.NewRouter()
	lg := logrus.New()
	server := Server{
		lg:    lg,
		Title: "Megablog",
		Posts: PostItems{
			{1, "Title1", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2018",
				Author{"/www/static/image/avatar.jpg", "Elon Musk"},
			},
			{2, "Title2", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2016",
				Author{"./www/static/image/avatar.jpg", "Elon Musk"},
			},
			{3, "Title3", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2016",
				Author{"./www/static/image/avatar.jpg", "Elon Musk"},
			},
			{4, "Title4", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2016",
				Author{"./www/static/image/avatar.jpg", "Elon Musk"},
			},
			{5, "Title5", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2016",
				Author{"./www/static/image/avatar.jpg", "Elon Musk"},
			},
			{6, "Title6", "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet.",
				"April 01, 2016",
				Author{"./www/static/image/avatar.jpg", "Elon Musk"},
			},
		},
	}

	r.Route("/", func(r chi.Router) {
		r.Get("/blog", server.HandleGetBlog)
		r.Get("/blog/post/{id}", server.HandleGetBlogPost)
	})

	fmt.Println("start")
	go func() {
		err := http.ListenAndServe(":9090", r)
		if err != nil {
			lg.WithError(err)
		}
	}()

	signal.Notify(stopchan, os.Interrupt, os.Kill)
	<-stopchan
	fmt.Println("\nstop")
}

func (s *Server) HandleGetBlog(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("./www/static/blogs/index.html")
	data, err := ioutil.ReadAll(file)

	templ := template.Must(template.New("blog").Parse(string(data)))

	err = templ.ExecuteTemplate(w, "blog", s)

	if err != nil {
		s.lg.WithError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) HandleGetBlogPost(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	file, _ := os.Open("./www/static/blogs/post.html")
	data, _ := ioutil.ReadAll(file)

	templ := template.Must(template.New("post").Parse(string(data)))

	p := PostItem{}

	s.lg.Println(id)
	for _, cp := range s.Posts {
		if cp.Id == id {
			p = cp
			break
		}
	}

	if p.Id == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	s.lg.Println(p)

	err := templ.ExecuteTemplate(w, "post", p)
	if err != nil {
		s.lg.WithError(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
