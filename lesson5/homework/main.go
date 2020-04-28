package main

import (
	"database/sql"
	"flag"
	"os"
	"os/signal"

	"github.com/hneis/web_begin/lesson5/homework/server"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/sqlboiler/boil"

	_ "github.com/go-sql-driver/mysql"
)

func newLogger() *logrus.Logger {
	lg := logrus.New()
	lg.SetReportCaller(false)
	lg.SetFormatter(&logrus.TextFormatter{})
	lg.SetLevel(logrus.DebugLevel)

	return lg
}

func main() {
	flagRootDir := flag.String("rootdir", "./www", "root dir of the server")
	flagServAddr := flag.String("addr", "localhost:8888", "server address")
	flag.Parse()

	lg := newLogger()
	db, err := sql.Open("mysql", "root:root@tcp(172.22.0.2)/blog")
	if err != nil {
		lg.WithError(err).Fatal("can't connect to db")
	}
	defer db.Close()
	serv := server.New(lg, *flagRootDir, db)

	boil.DebugMode = true
	go func() {
		err := serv.Start(*flagServAddr)
		if err != nil {
			lg.WithError(err).Fatal("can't run the server")
		}
	}()

	stopSig := make(chan os.Signal)
	signal.Notify(stopSig, os.Interrupt, os.Kill)
	<-stopSig

}
