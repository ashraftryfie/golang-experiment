package greeting

import (
	"errors"
	"net/http"
)

var ErrNotImplemented = errors.New("TODO: implement NewGreetingRouter")

type GreetResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewGreetingRouter initializes an http.ServeMux with a GET /greet/{name} route.
func NewGreetingRouter() (*http.ServeMux, error) {
	// TODO: Create http.NewServeMux(), register "GET /greet/{name}", return error if unimplemented
	return nil, ErrNotImplemented
}
