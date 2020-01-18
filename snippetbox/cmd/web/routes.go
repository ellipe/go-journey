package main

import (
	"net/http"
)

func configRouter(mux *http.ServeMux, app *application) {
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
}
