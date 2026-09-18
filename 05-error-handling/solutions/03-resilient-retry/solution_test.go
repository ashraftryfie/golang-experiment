package retry

import (
	"errors"
	"testing"
	"time"
)

func TestSolutionRetry(t *testing.T) {
	calls := 0
	err := Retry(3, time.Millisecond, func() error {
		calls++
		if calls < 2 {
			return NetworkError{Msg: "slow pipe"}
		}
		return nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}

	// Non-retryable error
	err = Retry(5, time.Millisecond, func() error {
		return AuthError{Msg: "bad creds"}
	})
	var authErr AuthError
	if !errors.As(err, &authErr) {
		t.Errorf("expected AuthError")
	}
}
