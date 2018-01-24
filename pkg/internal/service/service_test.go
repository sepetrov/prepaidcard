package service_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/sepetrov/prepaidcard/pkg/internal/service"
	h "github.com/sepetrov/prepaidcard/pkg/internal/testing"
)

func TestNewInternalServerErrorResponse(t *testing.T) {
	r := service.NewInternalServerErrorResponse()
	t.Run("sets status code 500", func(t *testing.T) {
		h.MustE(t, r.StatusCode(), 500, "got status code %#v, want %#v")
	})
	t.Run("sets title Internal Server Error", func(t *testing.T) {
		h.MustE(t, r.Title, http.StatusText(500), "got title %q, want %q")
	})
}

func TestErrorResponse_StatusCode(t *testing.T) {
	t.Run("default status code of StatusCoder interface is 500", func(t *testing.T) {
		r := service.ErrorResponse{}
		h.MustE(t, r.StatusCode(), 500, "got status code %#v, want %#v")
	})
}

func TestErrorResponse_Headers(t *testing.T) {
	t.Run("sets header content-type: application/problem+json", func(t *testing.T) {
		r := service.ErrorResponse{}
		h.MustE(t, r.Headers().Get("content-type"), "application/problem+json", "got Content-Type: %#v, want %#v")
	})
}

func TestErrorResponse_String(t *testing.T) {
	t.Run("uses the title if not empty", func(t *testing.T) {
		r := service.ErrorResponse{Title: "Foo"}
		h.MustE(t, r.String(), "Foo", "got title %q, want %q")
	})
	t.Run("uses the HTTP status code text if title is empty", func(t *testing.T) {
		r := service.ErrorResponse{Status: 201}
		h.MustE(t, r.String(), "Created", "got title %q, want %q")
	})
	t.Run("uses status code 500 if the status code is not recognised", func(t *testing.T) {
		r := service.ErrorResponse{Status: 13}
		h.MustE(t, r.String(), "Internal Server Error", "got title %q, want %q")
	})
}

func TestErrorResponse_Body(t *testing.T) {
	t.Run("must return json encoded error response", func(t *testing.T) {
		r := service.ErrorResponse{}
		d := struct {
			Title  string `json:"title"`
			Status int    `json:"status"`
		}{
			http.StatusText(500),
			500,
		}

		got, err := r.MarshalJSON()
		h.MustNotErr(t, err, "got JSON-encoding error; %v")

		want, err := json.Marshal(d)
		h.MustNotErr(t, err, "got JSON-encoding error; %v")

		h.MustE(t, string(got), string(want), "got %s != %s, want them equal")
	})
}
