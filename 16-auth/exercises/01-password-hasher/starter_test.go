package passwordhasher

import (
	"errors"
	"testing"
)

func TestPasswordHasher_Starter(t *testing.T) {
	hash, err := HashPassword("test-password-123", 10)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: HashPassword")
	}
	if err != nil {
		t.Fatalf("unexpected hash error: %v", err)
	}

	if !VerifyPassword(hash, "test-password-123") {
		t.Errorf("expected valid password verification to succeed")
	}
}
