package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/hneis/web_begin/lesson8/homework/config"
	"github.com/hneis/web_begin/lesson8/homework/server"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MONGO_URI = "mongodb://localhost:27017"

func newLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)

	return lg
}

func main() {
	flagConfiPath := flag.String("config", "./config.yaml", "Default yaml config path")
	flag.Parse()

	conf, err := config.ReadConfig(*flagConfiPath)
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}
	lg, err := config.NewLogger(&conf.Logger)

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		lg.Fatal(err)
	}

	db := client.Database("blog")

	serv := server.New(&conf.Server, lg, db, ctx)

	go func() {
		err := serv.Start()
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill)
	<-stopSig

}
