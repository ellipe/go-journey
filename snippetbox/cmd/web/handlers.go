package main

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

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

	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

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
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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

// Add a createSnippet handler function.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// If it's not, then send a 405 status response code indicating that it is
		// not allowed.
		w.Header().Set("Allow", http.MethodPost)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
