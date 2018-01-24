package service

import (
	"encoding/json"
	"net/http"
)

const errContentType = "application/problem+json"
const errStatusCode = http.StatusInternalServerError

// NewInternalServerErrorResponse returns 500 Internal Server Error.
func NewInternalServerErrorResponse() ErrorResponse {
	return ErrorResponse{
		Title:  http.StatusText(errStatusCode),
		Status: errStatusCode,
	}
}

// StatusCoder is used to set response status code.
type StatusCoder interface {
	StatusCode() int
}

// Headerer is used to set response headers.
type Headerer interface {
	Headers() http.Header
}

// ErrorResponse represents a business level error. It is returned to the client.
// This response must not include sensitive information. For more information
// about error response see https://tools.ietf.org/html/rfc7807#section-3.1
type ErrorResponse struct {
	Type   string `json:"-"`
	Title  string `json:"-"`
	Status int    `json:"-"`
	Detail string `json:"-"`
}

var _ StatusCoder = &ErrorResponse{}
var _ Headerer = &ErrorResponse{}
var _ json.Marshaler = &ErrorResponse{}

// StatusCoder implements StatusCoder.
func (r *ErrorResponse) StatusCode() int {
	if r.Status > 0 {
		return r.Status
	}
	return errStatusCode
}

// Headers implements Headerer.
func (r *ErrorResponse) Headers() http.Header {
	h := http.Header{}
	h.Set("Content-Type", errContentType)
	return h
}

// String implements Stringer. It returns the title.
func (r *ErrorResponse) String() string {
	if len(r.Title) > 0 {
		return r.Title
	}
	t := http.StatusText(r.StatusCode())
	if len(t) > 0 {
		return t
	}
	return http.StatusText(errStatusCode)
}

// MarshalJSON implements json.Marshaller.
func (r ErrorResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type   string `json:"type,omitempty"`
		Title  string `json:"title"`
		Status int    `json:"status,omitempty"`
		Detail string `json:"detail,omitempty"`
	}{
		r.Type,
		r.String(),
		r.StatusCode(),
		r.Detail,
	})
}
