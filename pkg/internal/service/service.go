package service

import "net/http"

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
	Type   string `json:"type,omitempty"`
	Title  string `json:"title"`
	Status int    `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
}

var _ StatusCoder = &ErrorResponse{}
var _ Headerer = &ErrorResponse{}

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
