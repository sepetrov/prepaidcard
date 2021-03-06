// Package api is the public interface of the prepaid card API.
package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	dispatcher createcard.Dispatcher
	logger     *log.Logger
	middleware Middleware
	repository Repository
	version    string
}

// HandlerFunc is an adapter to allow regular functions with the signature of
// Handle method of Handler interface to be wrapped and used as the Handler
// interface. This is useful when writing middleware.
type HandlerFunc func(http.ResponseWriter, *http.Request) error

// Handle implements Handler.
func (h HandlerFunc) Handle(w http.ResponseWriter, r *http.Request) error {
	return h(w, r)
}

// Handler is an interface for handling HTTP request.
type Handler interface {
	Handle(http.ResponseWriter, *http.Request) error
}

// Middleware is API handler middleware.
type Middleware func(Handler) Handler

var noopMiddleware Middleware = func(h Handler) Handler {
	return h
}

// Repository is an interface that satisfies the individual services' (handlers') repositories.
type Repository interface {
	createcard.Saver
}

// Option configures an API instance.
type Option func(*API) (*API, error)

// VersionOption returns new option for setting API version.
func VersionOption(version string) Option {
	return func(api *API) (*API, error) {
		api.version = version
		return api, nil
	}
}

// LoggerOption returns new option for setting the logger.
func LoggerOption(logger *log.Logger) Option {
	return func(api *API) (*API, error) {
		api.logger = logger
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

// RepositoryOption returns new option for setting a repository.
func RepositoryOption(repository Repository) Option {
	return func(api *API) (*API, error) {
		api.repository = repository
		return api, nil
	}
}

// New returns new API configured with options.
func New(options ...Option) (*API, error) {
	api := &API{
		dispatcher: &dispatcher{},
		middleware: noopMiddleware,
		version:    Version,
	}
	var err error
	for _, option := range options {
		api, err = option(api)
		if err != nil {
			return &API{}, err
		}
	}
	if api.logger == nil {
		api.logger = log.New(ioutil.Discard, "", 0)
	}
	if api.repository == nil {
		return &API{}, errors.New("missing repository option")
	}

	return api, nil
}

// withMiddleware wraps handler h with middleware.
func (api *API) withMiddleware(h Handler) Handler {
	return api.middleware(
		middleware.ErrorLog(api.logger)(
			middleware.Error()(
				h,
			),
		),
	)
}

// Attach attaches the API handlers to mux.
func (api *API) Attach(mux *http.ServeMux) {
	mux.Handle(fmt.Sprintf("%s/card", basePath), handlerAdapter(api.CreateCardHandler()))
	mux.Handle(fmt.Sprintf("%s/version", basePath), handlerAdapter(api.VersionHandler()))
}

// VersionHandler returns the handler for API version.
func (api *API) VersionHandler() Handler {
	h := handler.Func(func(w http.ResponseWriter, r *http.Request) error {
		enc := json.NewEncoder(w)
		enc.Encode(struct {
			Version string `json:"version"`
		}{
			Version: api.version,
		})
		return nil
	})
	return api.withMiddleware(h)
}

// CreateCardHandler returns the handler for registration of new cards.
func (api *API) CreateCardHandler() Handler {
	h := handler.NewCreateCard(createcard.New(api.repository.(createcard.Saver), api.dispatcher))
	return api.withMiddleware(h)
}

func handlerAdapter(h Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Handle(w, r)
	})
}

// no-operational implementations

type saver struct{}

var _ createcard.Saver = &saver{}

func (s *saver) SaveCard(_ *model.Card) error { return nil }

type dispatcher struct{}

var _ createcard.Dispatcher = &dispatcher{}

func (d *dispatcher) DispatchCardCreated(_ event.CardCreated) {}
