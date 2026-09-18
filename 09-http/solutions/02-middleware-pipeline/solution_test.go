package authmw_solution

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAPIKeySolution(t *testing.T) {
	const validKey = "secret-token-xyz"
	downstreamCalled := false

	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("access granted"))
	})

	mw := RequireAPIKey(validKey)
	wrapped := mw(protectedHandler)

	// Missing key -> 401 Unauthorized
	reqUnauthorized := httptest.NewRequest(http.MethodGet, "/secure", nil)
	recUnauthorized := httptest.NewRecorder()
	wrapped.ServeHTTP(recUnauthorized, reqUnauthorized)

	if recUnauthorized.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", recUnauthorized.Code)
	}
	if downstreamCalled {
		t.Errorf("downstream should not be called")
	}

	var errPayload map[string]string
	_ = json.Unmarshal(recUnauthorized.Body.Bytes(), &errPayload)
	if errPayload["error"] != "unauthorized" {
		t.Errorf("expected error 'unauthorized', got %v", errPayload)
	}

	// Valid key -> 200 OK
	reqAuthorized := httptest.NewRequest(http.MethodGet, "/secure", nil)
	reqAuthorized.Header.Set("X-API-Key", validKey)
	recAuthorized := httptest.NewRecorder()
	wrapped.ServeHTTP(recAuthorized, reqAuthorized)

	if recAuthorized.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recAuthorized.Code)
	}
	if !downstreamCalled {
		t.Errorf("expected downstream called")
	}
}
