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
	version    string
}

// New returns new API.
func New() *API {
	return &API{
		saver:      &saver{},
		dispatcher: &dispatcher{},
		version:    Version,
	}
}

// Attach attaches the API handlers to mux.
func (api *API) Attach(mux *http.ServeMux) {
	mux.Handle(fmt.Sprintf("%s/card", basePath), handlerAdapter(api.CreateCardHandler()))
	mux.Handle(fmt.Sprintf("%s/version", basePath), handlerAdapter(api.VersionHandler()))
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
