package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.accme.com/cmd/internal/models"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type config struct {
	P   string // PORT
	H   string // HOST
	DSN string // MySQL data source name
}

var cfg config = config{
	P:   "4000",
	H:   "localhost",
	DSN: "web:34896728@/snippetbox?parseTime=true",
}

type application struct {
	errorlog      *log.Logger
	infolog       *log.Logger
	cfg           config
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	flag.StringVar(&cfg.P, "p", cfg.P, "HTTP network port")
	flag.StringVar(&cfg.H, "h", cfg.H, "HTTP network host")
	flag.StringVar(&cfg.DSN, "dsn", cfg.DSN, "MySQL data source name")

	flag.Parse()

	address := fmt.Sprintf("%s:%s", cfg.H, cfg.P)

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	infolog.Printf("%sDB connection to %s ...%s", Yellow, cfg.DSN, Reset)

	db, err := openDB(cfg.DSN)

	if err != nil {
		errorlog.Fatal(err)
	}

	infolog.Printf("%sConnected to %s%s", Yellow, cfg.DSN, Reset)

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorlog.Fatal(err)
	}

	app := &application{
		errorlog:      errorlog,
		infolog:       infolog,
		cfg:           cfg,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	infolog.Printf(fmt.Sprintf("%sStarting server on http://%s%s", Green, address, Reset))

	srv := &http.Server{
		Addr:     address,
		ErrorLog: errorlog,
		Handler:  app.routes(),
	}

	err = srv.ListenAndServe()

	errorlog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
