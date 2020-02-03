package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"ellipe.party/snippetbox/pkg/forms"
	"ellipe.party/snippetbox/pkg/models"
)

// Define a count handler that returns the number of go routines running using the endpoint /count
func (app *application) count(w http.ResponseWriter, r *http.Request) {
	count := runtime.NumGoroutine()
	w.Write([]byte(strconv.Itoa(count)))
}

// Define a home handler function which writes a byte slice containing
// the latest snippets created.
func (app *application) home(w http.ResponseWriter, r *http.Request) {

	s, err := app.snippets.Latest()

	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

// showSnippet : show the snippet defined by the query string "id"
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	// Attemtps to convert the URL Query String into an Integer
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	// Validates if id is bigger than 0
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

// Add a new createSnippetForm handler, which for now returns a placeholder response.
func (app *application) createSnippetForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

// Add a createSnippet handler function.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	// Parses the request body into a map r.PostForm
	err := r.ParseForm()

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")

	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{
			Form: form,
		})
		return
	}

	id, err := app.snippets.Insert(form.Get("title"), form.Get("content"), form.Get("expires"))

	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/%d", id), http.StatusSeeOther)
}
