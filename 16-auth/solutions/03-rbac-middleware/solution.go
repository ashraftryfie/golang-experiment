package rbac

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type userContextKey struct{}

var claimsKey = userContextKey{}

type UserClaims struct {
	Subject string `json:"sub"`
	Role    string `json:"role"`
	Expires int64  `json:"exp"`
}

// GetUserClaims extracts UserClaims from context.
func GetUserClaims(ctx context.Context) (*UserClaims, bool) {
	val := ctx.Value(claimsKey)
	if val == nil {
		return nil, false
	}
	claims, ok := val.(*UserClaims)
	return claims, ok
}

func verifyToken(tokenString string, secret []byte) (*UserClaims, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}

	signedData := parts[0] + "." + parts[1]
	expectedSig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid signature encoding")
	}

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(signedData))
	calculatedSig := h.Sum(nil)

	if subtle.ConstantTimeCompare(expectedSig, calculatedSig) != 1 {
		return nil, errors.New("invalid signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid payload encoding")
	}

	var claims UserClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, errors.New("invalid claims json")
	}

	if claims.Expires > 0 && time.Now().Unix() > claims.Expires {
		return nil, errors.New("token expired")
	}

	return &claims, nil
}

// RequireRoles returns HTTP middleware that validates JWT tokens and requires one of allowedRoles.
func RequireRoles(secret []byte, allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, `{"error":"missing or invalid authorization header"}`, http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			claims, err := verifyToken(tokenStr, secret)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// Check role permissions if specified
			if len(allowedRoles) > 0 {
				authorized := false
				for _, role := range allowedRoles {
					if claims.Role == role {
						authorized = true
						break
					}
				}
				if !authorized {
					http.Error(w, `{"error":"forbidden: insufficient permissions"}`, http.StatusForbidden)
					return
				}
			}

			// Attach claims to context
			ctx := context.WithValue(r.Context(), claimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
