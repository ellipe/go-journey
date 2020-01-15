package main

import (
	"log"
	"net/http"
	"strings"
)

func main() {
	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	// Serves the static files using a Http FileServer
	fileServer := http.FileServer(http.Dir("./ui/assets/"))

	// Removes part of the path so the files can be matched againts the filesystem.
	mux.Handle("/assets/", http.StripPrefix("/assets", neuter(fileServer)))

	// Use the http.ListenAndServe() function to start a new web server. We pass in // two parameters: the TCP network address to listen on (in this case ":4000") // and the servemux we just created. If http.ListenAndServe() returns an error // we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

// TODO: Implement sanitization to filesystem.
// https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings

func neuter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
