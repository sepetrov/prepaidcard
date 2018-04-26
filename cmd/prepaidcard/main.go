// Command prepaidcard starts the API server on port 8080.
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/sepetrov/prepaidcard/pkg/api"
)

var port = flag.String("port", "8080", "Port number")

// setCorsHeaders adds CORS headers to response writer w.
func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
	w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: be more strict
}

// corsMiddleware adds CORS headers to allow requests from Swagger UI.
func corsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		h.ServeHTTP(w, r)
	})
}

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		w.WriteHeader(http.StatusNotFound)
	})
	api, err := api.New(api.MiddlewareOption(corsMiddleware))
	if err != nil {
		log.Fatalf("cannot create an API instance: %v", err)
	}
	api.Attach(http.DefaultServeMux)
	log.Println(fmt.Sprintf("Listenging on port %s", *port))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
