package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	app.infolog.Printf("[%s] %s:%s%s", r.Method, app.cfg.H, app.cfg.P, r.URL)

	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	files := []string{
		"./ui/html/base.html",
		"./ui/html/partials/nav.html",
		"./ui/html/pages/home.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)

		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	app.infolog.Printf("[%s] %s:%s%s", r.Method, app.cfg.H, app.cfg.P, r.URL)
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil || id < 1 {
		app.notFound(w)

		return
	}

	fmt.Fprintf(w, "The snippet with ID %d is comming", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	app.infolog.Printf("[%s] %s:%s%s", r.Method, app.cfg.H, app.cfg.P, r.URL)
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)

		return
	}

	w.Write([]byte("snippetCreate"))
}
