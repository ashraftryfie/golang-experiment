package authservice

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthService_Lifecycle(t *testing.T) {
	secret := []byte("secret-key-must-have-16-bytes!!")
	svc := NewAuthService(secret)
	ts := httptest.NewServer(svc)
	defer ts.Close()

	// 1. Register Regular User
	regPayload := `{"email":"alice@example.com","password":"securepassword123","role":"user"}`
	resp, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBufferString(regPayload))
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	// 2. Duplicate Registration Rejection
	respDup, err := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBufferString(regPayload))
	if err != nil {
		t.Fatalf("duplicate register failed: %v", err)
	}
	defer respDup.Body.Close()
	if respDup.StatusCode != http.StatusConflict {
		t.Errorf("expected 409 Conflict, got %d", respDup.StatusCode)
	}

	// 3. Register Admin User
	adminPayload := `{"email":"admin@example.com","password":"adminpassword123","role":"admin"}`
	respAdminReg, _ := http.Post(ts.URL+"/auth/register", "application/json", bytes.NewBufferString(adminPayload))
	respAdminReg.Body.Close()

	// 4. Login Alice (Correct Password)
	loginPayload := `{"email":"alice@example.com","password":"securepassword123"}`
	respLogin, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBufferString(loginPayload))
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}
	defer respLogin.Body.Close()
	if respLogin.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for login, got %d", respLogin.StatusCode)
	}

	var aliceLogin loginResponse
	_ = json.NewDecoder(respLogin.Body).Decode(&aliceLogin)
	if aliceLogin.Token == "" {
		t.Fatalf("expected non-empty token")
	}

	// 5. Login with Bad Password
	badLoginPayload := `{"email":"alice@example.com","password":"wrongpassword"}`
	respBadLogin, _ := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBufferString(badLoginPayload))
	respBadLogin.Body.Close()
	if respBadLogin.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 Unauthorized for bad login, got %d", respBadLogin.StatusCode)
	}

	// 6. Access /api/profile with Alice's Token
	reqProfile, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/profile", nil)
	reqProfile.Header.Set("Authorization", "Bearer "+aliceLogin.Token)
	respProfile, err := http.DefaultClient.Do(reqProfile)
	if err != nil {
		t.Fatalf("profile request failed: %v", err)
	}
	defer respProfile.Body.Close()
	if respProfile.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for Alice profile, got %d", respProfile.StatusCode)
	}
	var aliceUser User
	_ = json.NewDecoder(respProfile.Body).Decode(&aliceUser)
	if aliceUser.Email != "alice@example.com" {
		t.Errorf("expected alice@example.com, got %s", aliceUser.Email)
	}

	// 7. Alice attempts to access /api/admin/metrics -> Expect 403 Forbidden
	reqAdminAsAlice, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/admin/metrics", nil)
	reqAdminAsAlice.Header.Set("Authorization", "Bearer "+aliceLogin.Token)
	respAdminAsAlice, err := http.DefaultClient.Do(reqAdminAsAlice)
	if err != nil {
		t.Fatalf("admin request as Alice failed: %v", err)
	}
	defer respAdminAsAlice.Body.Close()
	if respAdminAsAlice.StatusCode != http.StatusForbidden {
		t.Errorf("expected 403 Forbidden for Alice accessing admin, got %d", respAdminAsAlice.StatusCode)
	}

	// 8. Admin Logins and Accesses /api/admin/metrics -> Expect 200 OK
	adminLoginPayload := `{"email":"admin@example.com","password":"adminpassword123"}`
	respAdminLogin, _ := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewBufferString(adminLoginPayload))
	var adminLogin loginResponse
	_ = json.NewDecoder(respAdminLogin.Body).Decode(&adminLogin)
	respAdminLogin.Body.Close()

	reqAdmin, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/admin/metrics", nil)
	reqAdmin.Header.Set("Authorization", "Bearer "+adminLogin.Token)
	respAdmin, err := http.DefaultClient.Do(reqAdmin)
	if err != nil {
		t.Fatalf("admin request failed: %v", err)
	}
	defer respAdmin.Body.Close()
	if respAdmin.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for admin metrics, got %d", respAdmin.StatusCode)
	}
}
