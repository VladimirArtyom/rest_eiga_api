package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const version string = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn string
	}
}

type application struct {
	cfg    config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8000, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:123456@localhost:5432/greenlight?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	// Init logger
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Init database
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	logger.Printf("database connection pool is established")

	var app *application = &application{
		cfg:    cfg,
		logger: logger,
	}

	var router *httprouter.Router = app.routes()

	serve := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	logger.Printf("Starting server in %s mode on port %d", app.cfg.env, app.cfg.port)
	err = serve.ListenAndServe()
	logger.Fatal(err)

}

func openDB(cfg config) (*sql.DB, error) {
	// make a connection
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// validate the connection, if the connection cannot be mad in 5 seconds, return an error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // executes after the function return

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
