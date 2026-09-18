package useraccount

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestSerializeUser(t *testing.T) {
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
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: SerializeUser is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jsonStr := string(data)

	// Verify omitted/hidden fields
	if strings.Contains(jsonStr, "password_hash") || strings.Contains(jsonStr, "$2a$12") {
		t.Errorf("password_hash should be hidden with json:\"-\", got %s", jsonStr)
	}
	if strings.Contains(jsonStr, "api_secret") || strings.Contains(jsonStr, "super_secret_token_123") {
		t.Errorf("api_secret should be hidden with json:\"-\", got %s", jsonStr)
	}
	if strings.Contains(jsonStr, "phone_number") {
		t.Errorf("empty phone_number should be omitted with omitempty, got %s", jsonStr)
	}

	// Verify unmarshaling back
	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal generated json: %v", err)
	}

	if parsed["username"] != "gopher" {
		t.Errorf("username mismatch: got %v, want gopher", parsed["username"])
	}
	if parsed["role"] != "admin" {
		t.Errorf("role mismatch: got %v, want admin", parsed["role"])
	}
}
