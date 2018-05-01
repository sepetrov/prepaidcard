package handler_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/internal/event"
	"github.com/sepetrov/prepaidcard/pkg/internal/handler"
	"github.com/sepetrov/prepaidcard/pkg/internal/service/createcard"
	assert "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestNew(t *testing.T) {
	t.Run("renders the card details on success", func(t *testing.T) {
		s := &assert.Repository{}
		d := &dispatcher{}
		h := handler.NewCreateCard(createcard.New(s, d))

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		th := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Handle(context.TODO(), w, r)
		})
		th(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.MustE(t, resp.StatusCode, 201, "")
		assert.MustE(t, resp.Header.Get("Content-Type"), "application/json; charset=utf-8", "")
		assert.Must(t, strings.Contains(string(body), fmt.Sprintf(`"uuid":"%s"`, s.Card.UUID().String())), "")
	})
}

type dispatcher struct {
	e event.CardCreated
}

var _ createcard.Dispatcher = &dispatcher{}

func (d *dispatcher) DispatchCardCreated(e event.CardCreated) {
	d.e = e
}
