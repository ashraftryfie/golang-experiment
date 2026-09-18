package problem

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteProblem(t *testing.T) {
	rec := httptest.NewRecorder()
	prob := ProblemDetails{
		Title:    "Unprocessable Entity",
		Status:   http.StatusUnprocessableEntity,
		Detail:   "The provided payload has 1 validation error.",
		Instance: "/api/v1/users",
		InvalidParams: []InvalidParam{
			{Name: "email", Reason: "must be a valid email address"},
		},
	}

	err := WriteProblem(rec, prob)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: WriteProblem is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// 1. Verify Content-Type
	if rec.Header().Get("Content-Type") != "application/problem+json" {
		t.Errorf("expected application/problem+json header, got %q", rec.Header().Get("Content-Type"))
	}

	// 2. Verify Status Code
	if rec.Code != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}

	// 3. Verify Body Decoding
	var decoded ProblemDetails
	if err := json.Unmarshal(rec.Body.Bytes(), &decoded); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if decoded.Type != "about:blank" {
		t.Errorf("expected default type 'about:blank', got %q", decoded.Type)
	}
	if len(decoded.InvalidParams) != 1 || decoded.InvalidParams[0].Name != "email" {
		t.Errorf("invalid params mismatch: %+v", decoded.InvalidParams)
	}
}
