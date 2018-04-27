// Command prepaidcard starts the API server on port 8080.
package main

import (
	"context"
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

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		setCorsHeaders(w)
		w.WriteHeader(http.StatusNotFound)
	})

	middlewareOption := api.MiddlewareOption(func(h api.Handler) api.Handler {
		return api.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			setCorsHeaders(w)
			return h.Handle(ctx, w, r)
		})
	})
	api, err := api.New(middlewareOption)
	if err != nil {
		log.Fatalf("cannot create an API instance: %v", err)
	}
	api.Attach(http.DefaultServeMux)
	log.Println(fmt.Sprintf("Listenging on port %s", *port))
	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
