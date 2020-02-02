package main

import (
	"net/http"

	"github.com/bmizerany/pat"
	"github.com/justinas/alice" // Allow composing middleware
)

func (app *application) routes() http.Handler {
	// Compose a common middleware to be used in every request.
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	mux := pat.New()
	mux.Get("/", http.HandlerFunc(app.home))
	mux.Get("/snippet/create", http.HandlerFunc(app.createSnippetForm))
	mux.Post("/snippet/create", http.HandlerFunc(app.createSnippet))
	mux.Get("/snippet/:id", http.HandlerFunc(app.showSnippet))
	mux.Get("/count", http.HandlerFunc(app.count))

	fileServer := http.FileServer(http.Dir("./ui/assets/"))
	mux.Get("/assets/", http.StripPrefix("/assets", neuter(fileServer)))

	return standardMiddleware.Then(mux)
}
