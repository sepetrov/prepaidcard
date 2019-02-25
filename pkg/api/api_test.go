// +build !integration

package api_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/api"
	assert "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestVersionHandler(t *testing.T) {
	t.Run("default version is unknown", func(t *testing.T) {
		a, err := api.New(
			api.RepositoryOption(&assert.Repository{}),
		)
		if err != nil {
			t.Fatalf("cannot create new API: %v", err)
		}
		h := a.VersionHandler()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot create request: %v", err)
		}
		h.Handle(w, r)
		assert.MustE(t, strings.TrimSpace(w.Body.String()), fmt.Sprintf(`{"version":%q}`, api.Version), "")
	})
	t.Run("returns the API version", func(t *testing.T) {
		v := "FooBar v123.456-beta"
		a, err := api.New(
			api.VersionOption(v),
			api.RepositoryOption(&assert.Repository{}),
		)
		if err != nil {
			t.Fatalf("cannot create new API: %v", err)
		}
		h := a.VersionHandler()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot create request: %v", err)
		}
		h.Handle(w, r)
		assert.MustE(t, strings.TrimSpace(w.Body.String()), fmt.Sprintf(`{"version":%q}`, v), "")
	})
}

func TestMiddlewareOption(t *testing.T) {
	var isCalled bool
	m := func(h api.Handler) api.Handler {
		return api.HandlerFunc(func(w http.ResponseWriter, r *http.Request) error {
			isCalled = true
			return h.Handle(w, r)
		})
	}
	a, err := api.New(
		api.MiddlewareOption(m),
		api.RepositoryOption(&assert.Repository{}),
	)
	if err != nil {
		t.Fatalf("cannot create API: %v", err)
	}
	a.VersionHandler().Handle(httptest.NewRecorder(), httptest.NewRequest("GET", "http://example.com", nil))
	if !isCalled {
		t.Error("middeware not used")
	}
}
