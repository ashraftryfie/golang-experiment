package authmw

import (
	"errors"
	"net/http"
)

var ErrNotImplemented = errors.New("TODO: implement RequireAPIKey")

// Middleware wraps an http.Handler with cross-cutting functionality.
type Middleware func(http.Handler) http.Handler

// RequireAPIKey validates the X-API-Key header before passing the request downstream.
func RequireAPIKey(expectedKey string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO: Check r.Header.Get("X-API-Key"), return 401 if mismatch, else call next.ServeHTTP(w, r)
			http.Error(w, ErrNotImplemented.Error(), http.StatusNotImplemented)
		})
	}
}
