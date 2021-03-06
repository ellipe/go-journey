package main

import (
	"html/template"
	"path/filepath"
	"time"

	"ellipe.party/snippetbox/pkg/forms"
	"ellipe.party/snippetbox/pkg/models"
)

// templateData : holds the information required by the template, since only one dynamic element can be sent
// to the HTML Template you can use an struct to hold data comming from different sources as one.
type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        *forms.Form
}

// humanDate : returns a human readable string date and time.
func humanDate(t time.Time) string {
	return (t.Local()).Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// newTemplateCache : Holds the templates in memory to avoid the i/o operations every single time
// a handler is caller.
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	// Initialize a new map to act as the cache.
	cache := map[string]*template.Template{}
	// Use the filepath.Glob function to get a slice of all filepaths with
	// the extension '.page.tmpl'. This essentially gives us a slice of all the
	// 'page' templates for the application.

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.tmpl"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Extract the file name (like 'home.page.tmpl') from the full file path
		// and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the page template file in to a template set and the custom functions.
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}
		// Use the ParseGlob method to add any 'layout' templates to the
		// template set (in our case, it's just the 'base' layout at the
		// moment).

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.tmpl"))
		if err != nil {
			return nil, err
		}

		// Use the ParseGlob method to add any 'partial' templates to the
		// template set (in our case, it's just the 'footer' partial at the
		// moment).
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.tmpl"))
		if err != nil {
			return nil, err
		}
		// Add the template set to the cache, using the name of the page // (like 'home.page.tmpl') as the key.
		cache[name] = ts
	}

	return cache, nil

}
