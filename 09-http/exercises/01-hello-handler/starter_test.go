package greeting

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreetingRouter(t *testing.T) {
	mux, err := NewGreetingRouter()
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: NewGreetingRouter is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error creating router: %v", err)
	}

	// 1. Valid greeting
	req := httptest.NewRequest(http.MethodGet, "/greet/Gopher", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}
	if rec.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected application/json header, got %q", rec.Header().Get("Content-Type"))
	}

	var resp GreetResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}
	if resp.Message != "Hello, Gopher!" || resp.Status != "success" {
		t.Errorf("unexpected body payload: %+v", resp)
	}

	// 2. Disallowed method (e.g. POST)
	postReq := httptest.NewRequest(http.MethodPost, "/greet/Gopher", nil)
	postRec := httptest.NewRecorder()
	mux.ServeHTTP(postRec, postReq)
	if postRec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 Method Not Allowed for POST, got %d", postRec.Code)
	}
}
