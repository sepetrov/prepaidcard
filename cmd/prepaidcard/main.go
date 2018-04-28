// Command prepaidcard starts the API server on port 8080.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

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
	logger := log.New(os.Stderr, "", log.LstdFlags)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logger.Printf("%s %s", r.Method, r.URL)
		setCorsHeaders(w)
		w.WriteHeader(http.StatusNotFound)
	})

	api, err := api.New(
		api.LoggerOption(logger),
		api.MiddlewareOption(func(h api.Handler) api.Handler {
			return api.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
				logger.Printf("%s %s", r.Method, r.URL)
				setCorsHeaders(w)
				return h.Handle(ctx, w, r)
			})
		}),
	)
	if err != nil {
		logger.Fatalf("cannot create an API instance: %v", err)
	}
	api.Attach(http.DefaultServeMux)
	logger.Printf("Listenging on port %s", *port)
	logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *port), nil))
}
