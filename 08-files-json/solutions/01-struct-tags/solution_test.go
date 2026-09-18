package useraccount_solution

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestSerializeUserSolution(t *testing.T) {
	user := UserAccount{
		ID:           42,
		Username:     "gopher",
		Email:        "gopher@golang.org",
		Role:         "admin",
		PhoneNumber:  "", // should be omitted
		IsActive:     true,
		PasswordHash: "$2a$12$e8ZbzK8...",
		APISecret:    "super_secret_token_123",
	}

	data, err := SerializeUser(user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonStr := string(data)

	if strings.Contains(jsonStr, "password_hash") || strings.Contains(jsonStr, "$2a$12") {
		t.Errorf("password_hash should be hidden, got %s", jsonStr)
	}
	if strings.Contains(jsonStr, "api_secret") || strings.Contains(jsonStr, "super_secret_token_123") {
		t.Errorf("api_secret should be hidden, got %s", jsonStr)
	}
	if strings.Contains(jsonStr, "phone_number") {
		t.Errorf("empty phone_number should be omitted, got %s", jsonStr)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed["username"] != "gopher" {
		t.Errorf("username got %v, want gopher", parsed["username"])
	}
}
