package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
)

// HandlerFunc is an adapter to allow regular functions with the signature of
// Handle method of Handler interface to be wrapped and used as the Handler
// interface. This is useful when writing middleware.
type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request) error

// Handle implements Handler.
func (h HandlerFunc) Handle(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return h(ctx, w, r)
}

// Handler is an interface for handling HTTP request.
type Handler interface {
	Handle(context.Context, http.ResponseWriter, *http.Request) error
}

// CreateCard is handler for new cards.
type CreateCard struct {
	svc *createcard.Service
}

var _ Handler = &CreateCard{}

// NewCreateCard returns CreateCard handler.
func NewCreateCard(svc *createcard.Service) *CreateCard {
	return &CreateCard{svc}
}

// Handle handles requests for new card.
func (h *CreateCard) Handle(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
	res, err := h.svc.CreateCard()
	if err != nil {
		return err
	}

	j, err := json.Marshal(res)
	if err != nil {
		return fmt.Errorf("got json.Marshal(%T) error; %v", res, err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
	return nil
}
