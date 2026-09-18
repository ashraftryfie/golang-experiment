package authmw_solution

import (
	"encoding/json"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

// RequireAPIKey validates the X-API-Key header before passing the request downstream.
func RequireAPIKey(expectedKey string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.Header.Get("X-API-Key")
			if key == "" || key != expectedKey {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
