package rbac

import (
	"context"
	"net/http"
)

type UserClaims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	Expires int64  `json:"exp"`
}

// GetUserClaims extracts UserClaims from context.
func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	return nil, false
}

// RequireRoles returns HTTP middleware that validates JWT tokens and requires one of allowedRoles.
func RequireRoles(secret []byte, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Exercise not implemented yet: return 501 Not Implemented
			http.Error(w, "not implemented", http.StatusNotImplemented)
		})
	}
}
