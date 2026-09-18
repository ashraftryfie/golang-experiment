package greeting_solution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGreetingRouterSolution(t *testing.T) {
	mux, err := NewGreetingRouter()
	if err != nil {
		t.Fatalf("unexpected error creating router: %v", err)
	}

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
}
