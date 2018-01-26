// Package api is the public interface of the prepaid card API.
package api

import (
	"context"
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
}

// New returns new API.
func New() *API {
	return &API{
		saver:      &saver{},
		dispatcher: &dispatcher{},
	}
}

// Attach attaches the API handlers to mux.
func (api *API) Attach(mux *http.ServeMux) {
	mux.Handle(fmt.Sprintf("%s/card", basePath), routeAdapter(api.CreateCardHandler()))
}

// CreateCardHandler returns the handler for registration of new cards.
func (api *API) CreateCardHandler() handler.Handler {
	m := middleware.ErrorMiddleware()

	var h handler.Handler
	h = handler.NewCreateCard(createcard.New(api.saver, api.dispatcher))
	h = m(h)

	return h
}

func routeAdapter(h handler.Handler) http.Handler {
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
