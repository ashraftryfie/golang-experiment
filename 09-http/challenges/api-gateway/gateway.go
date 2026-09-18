package apigateway

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// GatewayConfig holds configuration parameters for the API gateway.
type GatewayConfig struct {
	APIKey string
}

// NewGateway builds and returns the complete API Gateway handler.
func NewGateway(cfg GatewayConfig) http.Handler {
	mux := http.NewServeMux()

	// Public health check route
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	// Protected services route
	mux.HandleFunc("GET /v1/services/{name}", func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.PathValue("name")
		if serviceName == "" {
			http.Error(w, `{"error":"missing service name"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"service": serviceName,
			"status":  "active",
		})
	})

	// Middleware 1: Request ID Injection
	withReqID := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")
			if reqID == "" {
				buf := make([]byte, 8)
				_, _ = rand.Read(buf)
				reqID = hex.EncodeToString(buf)
			}
			w.Header().Set("X-Request-ID", reqID)
			ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	// Middleware 2: Path-scoped Auth Middleware
	withAuth := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/v1/") {
				clientKey := r.Header.Get("X-API-Key")
				if clientKey == "" || clientKey != cfg.APIKey {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					_ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}

	// Wrap mux with middlewares
	return withReqID(withAuth(mux))
}

// RequestIDFromContext extracts the request ID injected by the gateway.
func RequestIDFromContext(ctx context.Context) string {
	if val, ok := ctx.Value(RequestIDKey).(string); ok {
		return val
	}
	return ""
}

func init() {
	_ = fmt.Sprintf("")
}
