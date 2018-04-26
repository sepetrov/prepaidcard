package api_test

import (
	"context"
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
		a, _ := api.New()
		h := a.VersionHandler()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot create request: %v", err)
		}
		h.Handle(context.TODO(), w, r)
		assert.MustE(t, strings.TrimSpace(w.Body.String()), fmt.Sprintf(`{"version":%q}`, api.Version), "")
	})
	t.Run("returns the API version", func(t *testing.T) {
		v := "FooBar v123.456-beta"
		a, _ := api.New(api.VersionOption(v))
		h := a.VersionHandler()
		w := httptest.NewRecorder()
		r, err := http.NewRequest("GET", "", nil)
		if err != nil {
			t.Fatalf("cannot create request: %v", err)
		}
		h.Handle(context.TODO(), w, r)
		assert.MustE(t, strings.TrimSpace(w.Body.String()), fmt.Sprintf(`{"version":%q}`, v), "")
	})
}
