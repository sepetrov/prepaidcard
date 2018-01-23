package createcard_test

import (
	"errors"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/sepetrov/prepaidcard/pkg/internal/event"
	"github.com/sepetrov/prepaidcard/pkg/internal/model"
	"github.com/sepetrov/prepaidcard/pkg/internal/service"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
	h "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestService_CreateCard(t *testing.T) {
	t.Run("saves, dispatches and returns the same card", func(t *testing.T) {
		s := &saver{}
		d := &dispatcher{}
		svc := createcard.New(s, d)
		r, err := svc.CreateCard()
		h.MustNotErr(t, err, "got svc.CreateCard() = %T, %#v, want nil", r)
		res, ok := r.(createcard.Response)
		h.Must(t, ok, "got svc.CreateCard() = %T, want createcard.Response", r)
		h.Must(t, d.e.UUID != uuid.Nil, "got dispatcher event UUID %q == uuid.Nil, want !uuid.Nil", d.e.UUID)
		h.MustE(t, s.c.UUID(), d.e.CardUUID, "got saved card UUID %q != dispatched card UUID %q, want the same")
		h.MustE(t, res.UUID, s.c.UUID().String(), "got response card UUID %q != saver card UUID %q, want them equal")
	})
	t.Run("returns error response and error if saver returns error", func(t *testing.T) {
		s := &saver{err: errors.New("test saver failed")}
		d := &dispatcher{}
		svc := createcard.New(s, d)
		r, err := svc.CreateCard()
		h.MustErr(t, err, "got svc.CreateCard() = %T, nil, want service.ErrorResponse, error", r)
		_, ok := r.(service.ErrorResponse)
		h.Must(t, ok, "got svc.CreateCard() = %T, want service.ErrorResponse", r)
	})
}

type saver struct {
	c   *model.Card
	err error
}

var _ createcard.Saver = &saver{}

func (s *saver) SaveCard(c *model.Card) error {
	s.c = c
	return s.err
}

type dispatcher struct {
	e event.CardCreated
}

var _ createcard.Dispatcher = &dispatcher{}

func (d *dispatcher) DispatchCardCreated(e event.CardCreated) {
	d.e = e
}
