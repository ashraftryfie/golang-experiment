package tokenbucket

import (
	"errors"
	"testing"
)

func TestTokenBucket_Starter(t *testing.T) {
	tb, err := NewTokenBucket(10.0, 5)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: NewTokenBucket")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !tb.Allow() {
		t.Fatalf("expected initial Allow to succeed")
	}
}
