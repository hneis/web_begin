package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/hneis/web_begin/lesson8/homework/models"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServerConfig struct {
	Addr          string `yaml:"addr"`
	TemplatesDir  string `yaml:"templatesDir"`
	IndexTemplate string `yaml:"indexTemplate"`
	RootDir       string `yaml:"rootDir"`
}

type Server struct {
	lg      *logrus.Logger
	db      *mongo.Database
	rootDir string
	Page    models.Page
	ctx     context.Context
	conf    *ServerConfig
}

func New(conf *ServerConfig, lg *logrus.Logger, db *mongo.Database, ctx context.Context) *Server {
	return &Server{
		lg:   lg,
		db:   db,
		Page: models.Page{},
		ctx:  context.Background(),
		conf: conf,
	}
}

func (serv *Server) Start() error {
	r := chi.NewRouter()
	serv.bindRoutes(r)
	serv.lg.Infof("Server is started on %s", serv.conf.Addr)
	serv.lg.Debugf("%v", serv.conf)

	return http.ListenAndServe(serv.conf.Addr, r)
}

func (serv *Server) SendErr(w http.ResponseWriter, err error, code int, obj ...interface{}) {
	serv.lg.WithField("data", obj).WithError(err).Error("server error")
	w.WriteHeader(code)
	errModel := models.ErrorModel{
		Code:     code,
		Err:      err.Error(),
		Desc:     "server error",
		Internal: obj,
	}
	data, _ := json.Marshal(errModel)
	w.Write(data)
}

func (serv *Server) SendInternalErr(w http.ResponseWriter, err error, obj ...interface{}) {
	serv.SendErr(w, err, http.StatusInternalServerError, obj)
}
