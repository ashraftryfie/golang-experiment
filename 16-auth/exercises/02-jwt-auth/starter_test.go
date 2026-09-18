package jwtauth

import (
	"errors"
	"testing"
	"time"
)

func TestJWT_Starter(t *testing.T) {
	secret := []byte("16bytessecretkey!!")
	claims := Claims{
		Subject: "user-1",
		Role:    "admin",
		Expires: time.Now().Add(1 * time.Hour).Unix(),
	}

	token, err := CreateToken(claims, secret)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: CreateToken")
	}
	if err != nil {
		t.Fatalf("unexpected error creating token: %v", err)
	}

	parsed, err := VerifyToken(token, secret)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: VerifyToken")
	}
	if err != nil {
		t.Fatalf("unexpected error verifying token: %v", err)
	}
	if parsed.Subject != "user-1" {
		t.Errorf("expected subject user-1, got %s", parsed.Subject)
	}
}
