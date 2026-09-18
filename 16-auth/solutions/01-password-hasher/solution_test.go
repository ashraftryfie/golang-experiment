package passwordhasher

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestPasswordHasher_RoundTrip(t *testing.T) {
	password := "correct-horse-battery-staple"
	hash, err := HashPassword(password, bcrypt.MinCost)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !VerifyPassword(hash, password) {
		t.Fatalf("VerifyPassword failed on correct password")
	}

	if VerifyPassword(hash, "wrong-password") {
		t.Fatalf("VerifyPassword succeeded on incorrect password")
	}
}

func TestPasswordHasher_EmptyPassword(t *testing.T) {
	_, err := HashPassword("", bcrypt.MinCost)
	if !errors.Is(err, ErrEmptyPassword) {
		t.Fatalf("expected ErrEmptyPassword, got %v", err)
	}

	if VerifyPassword("somehash", "") {
		t.Fatalf("expected empty password verification to fail")
	}
}

func TestPasswordHasher_NeedsRehash(t *testing.T) {
	hashLowCost, err := HashPassword("pass123", bcrypt.MinCost) // cost 4
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Should need rehash if target is 10
	if !NeedsRehash(hashLowCost, 10) {
		t.Errorf("expected NeedsRehash=true for cost %d vs target 10", bcrypt.MinCost)
	}

	// Should not need rehash if target is 4
	if NeedsRehash(hashLowCost, bcrypt.MinCost) {
		t.Errorf("expected NeedsRehash=false for identical cost")
	}
}
