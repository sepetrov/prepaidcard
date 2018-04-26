package middleware_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/internal/handler"
	"github.com/sepetrov/prepaidcard/pkg/internal/handler/middleware"
	assert "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestErrorMiddleware(t *testing.T) {
	t.Run("renders net/http package error if not service.ErrorResponse", func(t *testing.T) {
		m := middleware.ErrorMiddleware()

		var h handler.Handler
		h = handler.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
			return errors.New("foo")
		})
		h = m(h)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		th := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Handle(context.TODO(), w, r)
		})
		th(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.MustE(t, resp.StatusCode, 500, "")
		assert.MustE(t, strings.TrimSpace(string(body)), "Internal Server Error", "")
	})
	t.Run("renders the information from service.ErrorResponse error", func(t *testing.T) {
		t.Skipf("TODO")
	})
	t.Run("does nothing if the handle does not return error", func(t *testing.T) {
		m := middleware.ErrorMiddleware()

		var h handler.Handler
		h = handler.HandlerFunc(func(_ context.Context, w http.ResponseWriter, _ *http.Request) error {
			w.Write([]byte("test"))
			return nil
		})
		h = m(h)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		th := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Handle(context.TODO(), w, r)
		})
		th(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.MustE(t, resp.StatusCode, 200, "")
		assert.MustE(t, strings.TrimSpace(string(body)), "test", "")
	})
}
