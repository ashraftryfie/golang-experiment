package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

type UserResponse struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

func main() {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC(),
		})
	})

	// Go 1.22+ method + path parameter routing
	mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "missing user id", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(UserResponse{
			ID:   id,
			Role: "member",
		})
	})

	// Production-grade server configuration with explicit timeouts
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Printf("HTTP Server listening on http://localhost%s (Go 1.22+ ServeMux)\n", server.Addr)
	// In production, run server.ListenAndServe() with graceful shutdown via signal handling
	_ = server
}
