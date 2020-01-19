package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ellipe.party/snippetbox/pkg/models"
)

// Define a home handler function which writes a byte slice containing
// "Hello from Snippetbox" as the response body.
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

	// As a JSON
	// for _, snippet := range s {
	// 	data, err := json.MarshalIndent(snippet, "", "")
	// 	if err != nil {
	// 		app.serverError(w, err)
	// 	}
	// 	fmt.Fprintf(w, "%s\n", data)
	// }

	// As a String.
	for _, snippet := range s {
		fmt.Fprintf(w, "%v\n", snippet)
	}

	// // Parses the templates files.
	// // Initialize a slice containing the paths to the two files. Note that the // home.page.tmpl file must be the *first* file in the slice.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }

	// // ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	// ts, err := template.ParseFiles(files...)

	// // If an error occurs during parsing an error is send to the client with a 500 error message.
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }

	// // Compiles the template and passes the variables as a second parameter, in this case nil.
	// err = ts.Execute(w, nil)

	// if err != nil {
	// 	app.serverError(w, err)
	// }

}

// Add a showSnippet handler function.
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
	fmt.Fprintf(w, "%v", s)
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
