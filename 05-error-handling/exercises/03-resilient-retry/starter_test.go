package retry

import (
	"errors"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	attempts := 0
	flakyFn := func() error {
		attempts++
		if attempts < 3 {
			return NetworkError{Msg: "connection reset"}
		}
		return nil
	}

	err := Retry(5, 5*time.Millisecond, flakyFn)
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("skipping: Retry is not yet implemented (implement in starter.go)")
	}
	if err != nil {
		t.Fatalf("expected retry to succeed on 3rd attempt, got %v", err)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}

	// Permanent error aborts immediately
	attempts = 0
	fatalFn := func() error {
		attempts++
		return AuthError{Msg: "invalid api key"}
	}

	err = Retry(5, 5*time.Millisecond, fatalFn)
	if err == nil {
		t.Errorf("expected fatal error")
	}
	if attempts != 1 {
		t.Errorf("expected exactly 1 attempt for permanent error, got %d", attempts)
	}
}
