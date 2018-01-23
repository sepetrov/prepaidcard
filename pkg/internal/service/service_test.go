package service_test

import (
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
