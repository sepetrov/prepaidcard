// Package api is the public interface of the prepaid card API.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sepetrov/prepaidcard/pkg/internal/event"
	"github.com/sepetrov/prepaidcard/pkg/internal/handler"
	"github.com/sepetrov/prepaidcard/pkg/internal/handler/middleware"
	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
)

const basePath = "/api"

// API is the prepaid card application.
type API struct {
	saver      createcard.Saver
	dispatcher createcard.Dispatcher
	middleware Middleware
	version    string
}

// Middleware is API handler middleware.
type Middleware func(http.Handler) http.Handler

// Option configures an API instance.
type Option func(*API) (*API, error)

// VersionOption returns new option for setting API version.
func VersionOption(version string) Option {
	return func(api *API) (*API, error) {
		api.version = version
		return api, nil
	}
}

// MiddlewareOption returns new option for setting middleware to api.
func MiddlewareOption(middleware Middleware) Option {
	return func(api *API) (*API, error) {
		api.middleware = middleware
		return api, nil
	}
}

// New returns new API configured with options.
func New(options ...Option) (*API, error) {
	api := &API{
		saver:      &saver{},
		dispatcher: &dispatcher{},
		version:    Version,
	}
	var err error
	for _, option := range options {
		api, err = option(api)
		if err != nil {
			return &API{}, err
		}
	}
	return api, nil
}

// withMiddleware wraps handler h with the configuration middleware.
func (api *API) withMiddleware(h http.Handler) http.Handler {
	if api.middleware == nil {
		return h
	}
	return api.middleware(h)
}

// Attach attaches the API handlers to mux.
func (api *API) Attach(mux *http.ServeMux) {
	mux.Handle(fmt.Sprintf("%s/card", basePath), api.withMiddleware(handlerAdapter(api.CreateCardHandler())))
	mux.Handle(fmt.Sprintf("%s/version", basePath), api.withMiddleware(handlerAdapter(api.VersionHandler())))
}

// VersionHandler returns the handler for API version.
func (api *API) VersionHandler() handler.Handler {
	return handler.HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) error {
		enc := json.NewEncoder(w)
		enc.Encode(struct {
			Version string `json:"version"`
		}{
			Version: api.version,
		})
		return nil
	})
}

// CreateCardHandler returns the handler for registration of new cards.
func (api *API) CreateCardHandler() handler.Handler {
	m := middleware.ErrorMiddleware()

	var h handler.Handler
	h = handler.NewCreateCard(createcard.New(api.saver, api.dispatcher))
	h = m(h)

	return h
}

func handlerAdapter(h handler.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Handle(context.TODO(), w, r)
	})
}

// no-operational implementations

type saver struct{}

var _ createcard.Saver = &saver{}

func (s *saver) SaveCard(_ *model.Card) error { return nil }

type dispatcher struct{}

var _ createcard.Dispatcher = &dispatcher{}

func (d *dispatcher) DispatchCardCreated(_ event.CardCreated) {}
