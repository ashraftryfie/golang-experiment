package rbac

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRBAC_Starter(t *testing.T) {
	secret := []byte("secret-key-must-be-16bytes!!")
	mw := RequireRoles(secret, "admin")

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	ts := httptest.NewServer(mw(dummyHandler))
	defer ts.Close()

	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotImplemented {
		t.Skip("Skipping unimplemented exercise: RequireRoles")
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for missing token, got %d", resp.StatusCode)
	}
}
