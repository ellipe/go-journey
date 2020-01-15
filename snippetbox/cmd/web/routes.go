package main

import (
	"net/http"
)

func configRouter(mux *http.ServeMux) {
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
}
