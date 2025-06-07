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

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/jsonlog"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

const version string = "1.0.0"

type config struct {
	port int
	env  string
	db   struct {
		dsn                string
		maxOpenConnections int
		maxIdleConnections int
		maxIdleTimeout     int
	}
}

type application struct {
	cfg    config
	logger *jsonlog.Logger
	models data.Models
}

func main() {
	var cfg config

	//　古いlogger
	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	// 新しいLogger
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)
	// load .env file
	err := godotenv.Load()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	port := convertStrToInt(os.Getenv("PORT"), logger)
	db_max_open_conns := convertStrToInt(os.Getenv("DATABASE_MAX_OPEN_CONNECTIONS"), logger)
	db_max_idle_conns := convertStrToInt(os.Getenv("DATABASE_MAX_IDLE_CONNECTIONS"), logger)
	db_max_idle_timeout := convertStrToInt(os.Getenv("DATABASE_MAX_IDLE_TIMEOUT"), logger)

	flag.IntVar(&cfg.port, "port", port, "Api server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("DATABASE_URL"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConnections, "db-max-open-conns", db_max_open_conns, "PosgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConnections, "db-max-idle-conns", db_max_idle_conns, "PosgreSQL max idle connections")
	flag.IntVar(&cfg.db.maxIdleTimeout, "db-max-idle-timeout", db_max_idle_timeout, "PosgreSQL max idle timeout")
	flag.Parse()

	// Init database
	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()

	logger.PrintInfo("database connection pool is established", nil)

	var app *application = &application{
		cfg:    cfg,
		logger: logger,
		models: data.NewModels(db),
	}

	var router *httprouter.Router = app.routes()

	// サーバーオブジェクトからのすべてのERRORが処理されています。 (All error from server objects are handled)
	serve := http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      router,
		ErrorLog:     log.New(app.logger, "", 0), // 新しい実証 ( New implementation  )
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// 正しい形式に変更しますた。 (I changed it to the correct format)
	logger.PrintInfo("Starting server", map[string]string{
		"addr": serve.Addr,
		"env":  app.cfg.env,
	})
	err = serve.ListenAndServe()
	logger.PrintFatal(err, map[string]string{"message": "Unable to start the server."})
}

func openDB(cfg config) (*sql.DB, error) {
	// make a connection
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}
	// set the db settings
	db.SetMaxOpenConns(cfg.db.maxOpenConnections)                             // the max connection for idle and active
	db.SetMaxIdleConns(cfg.db.maxIdleConnections)                             // the max idle limit that the db has.
	db.SetConnMaxIdleTime(time.Minute * time.Duration(cfg.db.maxIdleTimeout)) // timeout before the idle connection releasing its resources

	// validate the connection, if the connection cannot be mad in 5 seconds, return an error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // executes after the function return

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
