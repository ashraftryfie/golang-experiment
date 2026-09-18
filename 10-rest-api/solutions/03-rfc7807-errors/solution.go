package problem_solution

import (
	"encoding/json"
	"net/http"
)

type InvalidParam struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type ProblemDetails struct {
	Type          string         `json:"type"`
	Title         string         `json:"title"`
	Status        int            `json:"status"`
	Detail        string         `json:"detail,omitempty"`
	Instance      string         `json:"instance,omitempty"`
	InvalidParams []InvalidParam `json:"invalid_params,omitempty"`
}

// WriteProblem sets headers, writes the HTTP status code, and serializes the problem details response.
func WriteProblem(w http.ResponseWriter, p ProblemDetails) error {
	if p.Type == "" {
		p.Type = "about:blank"
	}

	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(p.Status)

	return json.NewEncoder(w).Encode(p)
}
