package main

import (
	"net/http"

	"github.com/justinas/alice" // Allow composing middleware
)

func (app *application) routes() http.Handler {
	// Compose a common middleware to be used in every request.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/count", app.count)

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets", neuter(fileServer)))

	return standardMiddleware.Then(mux)
}
