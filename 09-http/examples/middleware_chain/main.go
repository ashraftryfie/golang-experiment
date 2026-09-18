package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"
)

type contextKey string

const RequestIDKey contextKey = "request_id"

// Middleware defines a standard Go HTTP middleware signature.
type Middleware func(http.Handler) http.Handler

// Chain combines multiple middlewares from left to right.
func Chain(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// RequestIDMiddleware injects a unique correlation ID into the request context and response headers.
func RequestIDMiddleware(next http.Handler) http.Handler {
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

// LoggingMiddleware logs the duration and method of every request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID, _ := r.Context().Value(RequestIDKey).(string)

		next.ServeHTTP(w, r)

		log.Printf("[%s] %s %s took %v\n", reqID, r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := r.Context().Value(RequestIDKey).(string)
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprintf(w, "pong (request-id: %s)", reqID)
	})

	// Wrap entire router with chained middlewares
	handler := Chain(mux, RequestIDMiddleware, LoggingMiddleware)

	_ = &http.Server{
		Addr:    ":8080",
		Handler: handler,
	}
	fmt.Println("Server with middleware pipeline initialized.")
}
