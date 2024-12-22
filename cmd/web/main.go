package main

import (
	"flag"
	"log"
	"net/http"
	"strings"
)

// Prevent listing directories by sending requests to "/static/" etc.
func neuter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		h.ServeHTTP(w, r)
	})
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("../../ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", neuter(fileServer)))

	log.Println("Starting server on :4000")
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
