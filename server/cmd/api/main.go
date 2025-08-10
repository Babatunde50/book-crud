package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime/debug"
	"sync"

	"github.com/Babatunde50/book-crud/server/business/book"
	"github.com/Babatunde50/book-crud/server/business/book/bookdb"
	"github.com/Babatunde50/book-crud/server/business/urlprocessor"
	"github.com/Babatunde50/book-crud/server/internal/database"
	"github.com/Babatunde50/book-crud/server/internal/version"
	"github.com/lmittmann/tint"
)

// @title           Book Crud API
// @version         1.0
// @description     Backend assessment API for books and URL processing
// @termsOfService  http://swagger.io/terms/

// @contact.name   Babatunde Ololade
// @contact.email  babatundeololade50@gmail.com

// @host      localhost:4748
// @BasePath  /
// @schemes   http
func main() {
	logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

type config struct {
	baseURL  string
	httpPort int
	db       struct {
		dsn         string
		automigrate bool
	}
}

type application struct {
	config           config
	logger           *slog.Logger
	wg               sync.WaitGroup
	db               *database.DB
	bookCore         *book.Core
	urlProcessorCore *urlprocessor.URLProcessor
}

func run(logger *slog.Logger) error {
	var cfg config

	flag.StringVar(&cfg.baseURL, "base-url", "http://localhost:4748", "base URL for the application")
	flag.IntVar(&cfg.httpPort, "http-port", 4748, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "user:pass@localhost:5432/db", "postgreSQL DSN")
	flag.BoolVar(&cfg.db.automigrate, "db-automigrate", true, "run migrations on startup")

	showVersion := flag.Bool("version", false, "display version and exit")

	flag.Parse()

	if *showVersion {
		fmt.Printf("version: %s\n", version.Get())
		return nil
	}

	db, err := database.New(cfg.db.dsn, cfg.db.automigrate)
	if err != nil {
		return err
	}
	defer db.Close()

	bookStore := bookdb.New(db)
	bookCore := book.NewCore(bookStore)

	urlProcessorCore := urlprocessor.New()

	app := &application{
		config:           cfg,
		logger:           logger,
		db:               db,
		bookCore:         bookCore,
		urlProcessorCore: urlProcessorCore,
	}

	return app.serveHTTP()
}
