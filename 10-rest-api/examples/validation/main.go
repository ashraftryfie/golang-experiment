package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type CreateProductRequest struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
}

func (r *CreateProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("name is required and cannot be empty")
	}
	if r.Price <= 0 {
		return errors.New("price must be greater than zero")
	}
	if strings.TrimSpace(r.Category) == "" {
		return errors.New("category is required")
	}
	return nil
}

// DecodeAndValidate decodes a JSON request, rejecting unknown fields and empty bodies.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst interface{ Validate() error }) bool {
	// Limit request body size to 1MB to prevent DOS
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalTypeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxErr):
			http.Error(w, fmt.Sprintf("Malformed JSON at position %d", syntaxErr.Offset), http.StatusBadRequest)
		case errors.As(err, &unmarshalTypeErr):
			http.Error(w, fmt.Sprintf("Invalid type for field %q", unmarshalTypeErr.Field), http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			http.Error(w, "Request body cannot be empty", http.StatusBadRequest)
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			http.Error(w, fmt.Sprintf("Request body contains unknown key %s", fieldName), http.StatusBadRequest)
		default:
			http.Error(w, fmt.Sprintf("Decode error: %v", err), http.StatusBadRequest)
		}
		return false
	}

	if err := dst.Validate(); err != nil {
		http.Error(w, fmt.Sprintf("Validation failed: %v", err), http.StatusUnprocessableEntity)
		return false
	}

	return true
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/products", func(w http.ResponseWriter, r *http.Request) {
		var req CreateProductRequest
		if !DecodeAndValidate(w, r, &req) {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":  "created",
			"product": req,
		})
	})

	fmt.Println("REST Request Validation Example initialized.")
}
