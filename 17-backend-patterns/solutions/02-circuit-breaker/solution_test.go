package circuitbreaker

import (
	"errors"
	"testing"
	"time"
)

var errTestDownstream = errors.New("downstream service 500")

func TestCircuitBreaker_StateTransitions(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		MaxFailures:      2,
		Timeout:          50 * time.Millisecond,
		SuccessThreshold: 2,
	})

	// 1. Initially Closed
	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed, got %v", cb.State())
	}

	// 2. First failure
	_ = cb.Execute(func() error { return errTestDownstream })
	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after 1 failure, got %v", cb.State())
	}

	// 3. Second failure -> Trips to Open
	_ = cb.Execute(func() error { return errTestDownstream })
	if cb.State() != StateOpen {
		t.Fatalf("expected StateOpen after reaching MaxFailures, got %v", cb.State())
	}

	// 4. While Open, calls fail fast without running fn
	fnExecuted := false
	err := cb.Execute(func() error {
		fnExecuted = true
		return nil
	})
	if !errors.Is(err, ErrCircuitOpen) {
		t.Fatalf("expected ErrCircuitOpen, got %v", err)
	}
	if fnExecuted {
		t.Fatalf("fn should not have run while circuit is open")
	}

	// 5. Wait for timeout -> Should transition to Half-Open
	time.Sleep(60 * time.Millisecond)
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected StateHalfOpen after timeout, got %v", cb.State())
	}

	// 6. Canary 1 succeeds
	err = cb.Execute(func() error { return nil })
	if err != nil {
		t.Fatalf("expected canary 1 to succeed, got %v", err)
	}

	// 7. Canary 2 succeeds -> Heals back to Closed
	err = cb.Execute(func() error { return nil })
	if err != nil {
		t.Fatalf("expected canary 2 to succeed, got %v", err)
	}
	if cb.State() != StateClosed {
		t.Fatalf("expected StateClosed after reaching SuccessThreshold, got %v", cb.State())
	}
}

func TestCircuitBreaker_HalfOpenCanaryFails(t *testing.T) {
	cb := NewCircuitBreaker(Config{
		MaxFailures:      1,
		Timeout:          20 * time.Millisecond,
		SuccessThreshold: 2,
	})

	// Trip to open
	_ = cb.Execute(func() error { return errTestDownstream })
	if cb.State() != StateOpen {
		t.Fatalf("expected open state")
	}

	// Wait for half-open
	time.Sleep(30 * time.Millisecond)
	if cb.State() != StateHalfOpen {
		t.Fatalf("expected half-open state")
	}

	// Canary fails -> Immediately trips back to Open
	_ = cb.Execute(func() error { return errTestDownstream })
	if cb.State() != StateOpen {
		t.Fatalf("expected open state after canary failure, got %v", cb.State())
	}
}
