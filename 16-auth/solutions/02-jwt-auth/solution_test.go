package jwtauth

import (
	"errors"
	"strings"
	"testing"
	"time"
)

var testSecret = []byte("16-bytes-secret-key-32!!")

func TestJWT_RoundTrip(t *testing.T) {
	claims := Claims{
		Subject: "usr-42",
		Role:    "editor",
		Expires: time.Now().Add(10 * time.Minute).Unix(),
	}

	token, err := CreateToken(claims, testSecret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	parsed, err := VerifyToken(token, testSecret)
	if err != nil {
		t.Fatalf("VerifyToken failed: %v", err)
	}

	if parsed.Subject != "usr-42" || parsed.Role != "editor" {
		t.Errorf("claims mismatch: %+v", parsed)
	}
}

func TestJWT_InvalidSignature(t *testing.T) {
	claims := Claims{
		Subject: "usr-42",
		Role:    "user",
		Expires: time.Now().Add(10 * time.Minute).Unix(),
	}

	token, err := CreateToken(claims, testSecret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	wrongSecret := []byte("different-secret-key-32!")
	_, err = VerifyToken(token, wrongSecret)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	claims := Claims{
		Subject: "usr-expired",
		Role:    "user",
		Expires: time.Now().Add(-1 * time.Minute).Unix(), // expired 1 minute ago
	}

	token, err := CreateToken(claims, testSecret)
	if err != nil {
		t.Fatalf("CreateToken failed: %v", err)
	}

	_, err = VerifyToken(token, testSecret)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("expected ErrTokenExpired, got %v", err)
	}
}

func TestJWT_MalformedToken(t *testing.T) {
	_, err := VerifyToken("not-a-jwt", testSecret)
	if !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("expected ErrInvalidToken, got %v", err)
	}

	// Tampered payload
	parts := strings.Split("abc.def.ghi", ".")
	_, err = VerifyToken(strings.Join(parts, "."), testSecret)
	if err == nil {
		t.Fatalf("expected error on malformed parts")
	}
}

func TestJWT_ValidationErrors(t *testing.T) {
	_, err := CreateToken(Claims{Subject: "user"}, []byte("short"))
	if !errors.Is(err, ErrSecretTooShort) {
		t.Fatalf("expected ErrSecretTooShort, got %v", err)
	}

	_, err = CreateToken(Claims{Subject: ""}, testSecret)
	if !errors.Is(err, ErrEmptySubject) {
		t.Fatalf("expected ErrEmptySubject, got %v", err)
	}
}
