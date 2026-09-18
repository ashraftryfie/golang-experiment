package rbac

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testSecret = []byte("my-super-secret-key-at-least-16bytes")

func makeToken(sub, role string, exp time.Time) string {
	headerEnc := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	claimsBytes, _ := json.Marshal(UserClaims{Subject: sub, Role: role, Expires: exp.Unix()})
	payloadEnc := base64.RawURLEncoding.EncodeToString(claimsBytes)

	signedData := headerEnc + "." + payloadEnc
	h := hmac.New(sha256.New, testSecret)
	h.Write([]byte(signedData))
	sigEnc := base64.RawURLEncoding.EncodeToString(h.Sum(nil))

	return signedData + "." + sigEnc
}

func TestRBAC_AuthorizedRole(t *testing.T) {
	mw := RequireRoles(testSecret, "admin", "editor")

	var extractedSub string
	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := GetUserClaims(r.Context())
		if ok {
			extractedSub = claims.Subject
		}
		w.WriteHeader(http.StatusOK)
	}))

	ts := httptest.NewServer(handler)
	defer ts.Close()

	validToken := makeToken("admin-user", "admin", time.Now().Add(1*time.Hour))
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
	if extractedSub != "admin-user" {
		t.Fatalf("expected extracted subject 'admin-user', got '%s'", extractedSub)
	}
}

func TestRBAC_MissingHeaderReturns401(t *testing.T) {
	mw := RequireRoles(testSecret, "admin")
	ts := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized, got %d", resp.StatusCode)
	}
}

func TestRBAC_InsufficientRoleReturns403(t *testing.T) {
	mw := RequireRoles(testSecret, "admin")
	ts := httptest.NewServer(mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	defer ts.Close()

	viewerToken := makeToken("regular-viewer", "viewer", time.Now().Add(1*time.Hour))
	req, _ := http.NewRequest(http.MethodGet, ts.URL, nil)
	req.Header.Set("Authorization", "Bearer "+viewerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		fmt.Printf("Status was %d\n", resp.StatusCode)
		t.Errorf("expected 403 Forbidden for viewer accessing admin route, got %d", resp.StatusCode)
	}
}
