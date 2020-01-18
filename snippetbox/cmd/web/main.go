package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// Configuration flags
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Initialize a new instance of application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
		// Include the router as part of the application structure ?
	}

	// Use the http.NewServeMux() function to initialize a new servemux, then
	// register the home function as the handler for the "/" URL pattern.
	mux := http.NewServeMux()
	configRouter(mux, app)

	// Serves the static files using a Http FileServer
	fileServer := http.FileServer(http.Dir("./ui/assets/"))

	// Removes part of the path so the files can be matched againts the filesystem.
	mux.Handle("/assets/", http.StripPrefix("/assets", neuter(fileServer)))

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  mux,
	}

	// Use the http.ListenAndServe() function to start a new web server. We pass in // two parameters: the TCP network address to listen on (in this case ":4000") // and the servemux we just created. If http.ListenAndServe() returns an error // we use the log.Fatal() function to log the error message and exit. Note
	// that any error returned by http.ListenAndServe() is always non-nil.
	infoLog.Printf("Starting server on port %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
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
