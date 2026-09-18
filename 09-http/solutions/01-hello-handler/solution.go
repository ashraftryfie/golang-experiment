package greeting_solution

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type GreetResponse struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

// NewGreetingRouter initializes an http.ServeMux with Go 1.22+ method-based route.
func NewGreetingRouter() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /greet/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			http.Error(w, `{"error":"name required"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(GreetResponse{
			Message: fmt.Sprintf("Hello, %s!", name),
			Status:  "success",
		})
	})

	return mux, nil
}
