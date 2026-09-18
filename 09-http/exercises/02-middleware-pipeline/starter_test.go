package authmw

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequireAPIKey(t *testing.T) {
	const validKey = "secret-token-xyz"
	downstreamCalled := false

	// Dummy protected handler
	protectedHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		downstreamCalled = true
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("access granted"))
	})

	mw := RequireAPIKey(validKey)
	wrapped := mw(protectedHandler)

	// Case 1: Check unimplemented starter
	testReq := httptest.NewRequest(http.MethodGet, "/secure", nil)
	testRec := httptest.NewRecorder()
	wrapped.ServeHTTP(testRec, testReq)
	if testRec.Code == http.StatusNotImplemented && strings.Contains(testRec.Body.String(), "TODO: implement") {
		t.Skip("skipping: RequireAPIKey is not yet implemented (implement in starter.go)")
	}

	// Case 2: Missing API key -> 401 Unauthorized
	downstreamCalled = false
	reqUnauthorized := httptest.NewRequest(http.MethodGet, "/secure", nil)
	recUnauthorized := httptest.NewRecorder()
	wrapped.ServeHTTP(recUnauthorized, reqUnauthorized)

	if recUnauthorized.Code != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", recUnauthorized.Code)
	}
	if downstreamCalled {
		t.Errorf("downstream handler should NOT have been called on unauthorized request")
	}

	var errPayload map[string]string
	_ = json.Unmarshal(recUnauthorized.Body.Bytes(), &errPayload)
	if errPayload["error"] != "unauthorized" {
		t.Errorf("expected error: 'unauthorized', got %v", errPayload)
	}

	// Case 3: Valid API key -> 200 OK & downstream called
	downstreamCalled = false
	reqAuthorized := httptest.NewRequest(http.MethodGet, "/secure", nil)
	reqAuthorized.Header.Set("X-API-Key", validKey)
	recAuthorized := httptest.NewRecorder()
	wrapped.ServeHTTP(recAuthorized, reqAuthorized)

	if recAuthorized.Code != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", recAuthorized.Code)
	}
	if !downstreamCalled {
		t.Errorf("expected downstream handler to be called when API key is valid")
	}
}
