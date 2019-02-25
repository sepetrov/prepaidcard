// +build !integration

package middleware_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/internal/handler"
	"github.com/sepetrov/prepaidcard/pkg/internal/handler/middleware"
	assert "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestError(t *testing.T) {
	t.Run("renders net/http package error if not service.ErrorResponse", func(t *testing.T) {
		m := middleware.Error()

		var h handler.Handler
		h = handler.Func(func(w http.ResponseWriter, _ *http.Request) error {
			return errors.New("foo")
		})
		h = m(h)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		th := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Handle(w, r)
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
		m := middleware.Error()

		var h handler.Handler
		h = handler.Func(func(w http.ResponseWriter, _ *http.Request) error {
			w.Write([]byte("test"))
			return nil
		})
		h = m(h)

		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		th := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.Handle(w, r)
		})
		th(w, req)

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)

		assert.MustE(t, resp.StatusCode, 200, "")
		assert.MustE(t, strings.TrimSpace(string(body)), "test", "")
	})
}

func TestErrorLog(t *testing.T) {
	b := &bytes.Buffer{}
	l := log.New(b, "", 0)
	m := middleware.ErrorLog(l)
	t.Run("logs errors", func(t *testing.T) {
		defer b.Reset()
		e := errors.New("foo")
		h := m(handler.Func(func(http.ResponseWriter, *http.Request) error { return e }))
		err := h.Handle(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))
		if want := fmt.Sprintf("%s\n", e); b.String() != want {
			t.Errorf("want logged error %q, got %q", want, b.String())
		}
		if err != e {
			t.Errorf("want error %#v, got error %#v", e, err)
		}
	})
	t.Run("ignores sucessfully handled requests", func(t *testing.T) {
		defer b.Reset()
		h := m(handler.Func(func(http.ResponseWriter, *http.Request) error { return nil }))
		err := h.Handle(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))
		if b.String() != "" {
			t.Errorf("want no logs, got %q", b.String())
		}
		if err != nil {
			t.Errorf("want nil, got error %#v", err)
		}
	})
}
