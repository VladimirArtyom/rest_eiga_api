package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
)

const version string = "1.0.0"

type config struct  {
	port int
	env string
}

type application struct {
	cfg config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Init logger
	logger := log.New(os.Stdout, "", log.Ldate | log.Ltime)

	var app *application = &application{
		cfg: cfg,
		logger: logger,
	}

	var router *httprouter.Router = app.routes()
	

  serve := http.Server{
		Addr: fmt.Sprintf(":%d", app.cfg.port),
		Handler: router,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Starting server in %s mode on port %d", app.cfg.env, app.cfg.port)
	err := serve.ListenAndServe()
	logger.Fatal(err)

}
