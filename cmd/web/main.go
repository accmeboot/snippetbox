package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

type config struct {
	P string // PORT
	H string // HOST
}

var cfg config = config{
	P: "4000",
	H: "localhost",
}

type application struct {
	errorlog *log.Logger
	infolog  *log.Logger
	cfg      config
}

func main() {
	flag.StringVar(&cfg.P, "p", cfg.P, "HTTP network port")
	flag.StringVar(&cfg.H, "h", cfg.H, "HTTP network host")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", cfg.H, cfg.P)

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	app := &application{
		errorlog: errorlog,
		infolog:  infolog,
		cfg:      cfg,
	}

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	infolog.Printf(fmt.Sprintf("Starting server on http://%s", address))

	srv := &http.Server{
		Addr:     address,
		ErrorLog: errorlog,
		Handler:  mux,
	}

	err := srv.ListenAndServe()

	errorlog.Fatal(err)
}
