package revproxy

import (
	"errors"
	"net/http"
)

var ErrNotImplemented = errors.New("TODO: implement NewProxyHandler")

// NewProxyHandler forwards all incoming requests to the specified targetBaseURL.
func NewProxyHandler(targetBaseURL string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Build outbound request, copy headers, execute request, copy response headers/status/body
		http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
	})
}
