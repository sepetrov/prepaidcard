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

// handler returns 501 Not Implemented.
func handler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// corsMiddleware adds CORS headers to allow requests from Swagger UI.
func corsMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(w, r)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST")
		w.Header().Set("Access-Control-Allow-Origin", "*") // TODO: be more strict
	}
}

func main() {
	flag.Parse()
	http.HandleFunc("/", corsMiddleware(handler))
	api := api.New()
	api.Attach(http.DefaultServeMux)
	log.Println(fmt.Sprintf("Listenging on port %s", *port))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
