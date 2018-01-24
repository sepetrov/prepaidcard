package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sepetrov/prepaidcard/pkg/internal/handler"
	"github.com/sepetrov/prepaidcard/pkg/internal/service"
)

// Middleware is a handler.Handler wrapper.
type Middleware func(handler.Handler) handler.Handler

// ErrorMiddleware handles error returned by the wrapped handler prev.
// If the error is type service.ErrorResponse, it will be sent as a response.
// For all other errors a generic 500 service.ErrorResponse will be sent.
func ErrorMiddleware() Middleware {
	return func(prev handler.Handler) handler.Handler {
		return handler.HanderFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			err := prev.Handle(ctx, w, r)

			if err == nil {
				return nil
			}

			errRes, ok := err.(service.ErrorResponse)
			if !ok {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return err
			}

			j, err := errRes.MarshalJSON()
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return fmt.Errorf("got %#v.MarshalJSON() error %v; %T", errRes, err, prev)
			}

			w.WriteHeader(errRes.StatusCode())
			for k := range errRes.Headers() {
				w.Header().Set(k, errRes.Headers().Get(k))
			}
			w.Write(j)
			return nil
		})
	}
}
