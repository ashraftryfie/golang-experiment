package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

func TestCircuitBreaker_Starter(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		MaxFailures:      2,
		Timeout:          50 * time.Millisecond,
		SuccessThreshold: 1,
	})
	if cb == nil {
		t.Skip("Skipping unimplemented exercise: NewCircuitBreaker")
	}

	err := cb.Execute(func() error { return nil })
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: Execute")
	}
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
