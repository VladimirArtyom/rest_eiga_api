package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/jsonlog"
	"github.com/VladimirArtyom/rest_eiga_api/internal/mailer"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version string = "1.0.0"

type config struct {
	port int
	env  string
	cors struct {
		origins map[string]bool 
	}
	db   struct {
		dsn                string
		maxOpenConnections int
		maxIdleConnections int
		maxIdleTimeout     int
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	smtp struct {
		host string
		port int
		username string
		password string
		sender string
	}
}

type application struct {
	cfg    config
	logger *jsonlog.Logger
	models data.Models
	mailer *mailer.Mailer
	wg sync.WaitGroup
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

	// Flag set 始める
	var corsFlagSet bool = false
	// Flag Set 終わり

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

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// SMTPの変数
	smtp_port := convertStrToInt(os.Getenv("SMTP_PORT"), logger)
	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "The host of the mail server")
	flag.IntVar(&cfg.smtp.port, "smtp-port", smtp_port, "The port of the mail server")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME") , "The username of the mail server")
	flag.StringVar(&cfg.smtp.password, "smpt-password", os.Getenv("SMTP_PASSWORD"), "The password of the mail server")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "The sender of the mail")

	// Allowed origins
	flag.Func("cors-trusted-origin", "Trusted origins, separated by COMMA", func(origins string) error {
		
		for _, origin := range strings.Split(origins, ",") {
			cfg.cors.origins[origin] = true
		}
		corsFlagSet = true
		return nil
	},)
	
	if !corsFlagSet {
		cfg.cors.origins = parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
	}
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
		mailer: mailer.New(cfg.smtp.host,
						cfg.smtp.port,
						cfg.smtp.username,
						cfg.smtp.password,
						cfg.smtp.sender),
	}


	// サーバーオブジェクトからのすべてのERRORが処理されています。 (All error from server objects are handled)
	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, map[string]string{"message": "Unable to start the server."})
	}
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
