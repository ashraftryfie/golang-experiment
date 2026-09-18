package resilientrpc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestResilientClient_FirstAttemptSuccess(t *testing.T) {
	cfg := DefaultRetryConfig()
	client, err := NewResilientClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	var attempts int32
	mockRPC := func(ctx context.Context) (any, error) {
		atomic.AddInt32(&attempts, 1)
		return "payload", nil
	}

	res, err := client.Invoke(context.Background(), mockRPC)
	if err != nil {
		t.Fatalf("expected success, got %v", err)
	}
	if res != "payload" {
		t.Errorf("expected 'payload', got %v", res)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("expected exactly 1 attempt, got %d", attempts)
	}
}

func TestResilientClient_TransientFailureSuccess(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:       4,
		PerAttemptTimeout: 50 * time.Millisecond,
		InitialBackoff:    5 * time.Millisecond,
		MaxBackoff:        20 * time.Millisecond,
		BackoffFactor:     2.0,
	}
	client, err := NewResilientClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	var attempts int32
	errTransient := errors.New("transient network hiccup")

	mockRPC := func(ctx context.Context) (any, error) {
		att := atomic.AddInt32(&attempts, 1)
		if att < 3 {
			return nil, errTransient
		}
		return 42, nil
	}

	res, err := client.Invoke(context.Background(), mockRPC)
	if err != nil {
		t.Fatalf("expected eventual success, got %v", err)
	}
	if res != 42 {
		t.Errorf("expected 42, got %v", res)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestResilientClient_ExhaustAttempts(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:       3,
		PerAttemptTimeout: 20 * time.Millisecond,
		InitialBackoff:    2 * time.Millisecond,
		MaxBackoff:        10 * time.Millisecond,
		BackoffFactor:     2.0,
	}
	client, err := NewResilientClient(cfg)
	if err != nil {
		t.Fatalf("unexpected error creating client: %v", err)
	}

	var attempts int32
	errFatal := errors.New("upstream persistent outage")

	mockRPC := func(ctx context.Context) (any, error) {
		atomic.AddInt32(&attempts, 1)
		return nil, errFatal
	}

	res, err := client.Invoke(context.Background(), mockRPC)
	if !errors.Is(err, errFatal) {
		t.Fatalf("expected errFatal, got %v", err)
	}
	if res != nil {
		t.Errorf("expected nil result, got %v", res)
	}
	if atomic.LoadInt32(&attempts) != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestResilientClient_PerAttemptTimeout(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:       2,
		PerAttemptTimeout: 20 * time.Millisecond,
		InitialBackoff:    5 * time.Millisecond,
		MaxBackoff:        10 * time.Millisecond,
		BackoffFactor:     1.5,
	}
	client, _ := NewResilientClient(cfg)

	var attempts int32
	mockSlowRPC := func(ctx context.Context) (any, error) {
		atomic.AddInt32(&attempts, 1)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return "too slow", nil
		}
	}

	_, err := client.Invoke(context.Background(), mockSlowRPC)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded from attempt timeout, got %v", err)
	}
	if atomic.LoadInt32(&attempts) != 2 {
		t.Errorf("expected 2 attempts before exhausting, got %d", attempts)
	}
}

func TestResilientClient_ParentContextCancelledDuringBackoff(t *testing.T) {
	cfg := RetryConfig{
		MaxAttempts:       5,
		PerAttemptTimeout: 50 * time.Millisecond,
		InitialBackoff:    200 * time.Millisecond, // long sleep
		MaxBackoff:        500 * time.Millisecond,
		BackoffFactor:     2.0,
	}
	client, _ := NewResilientClient(cfg)

	ctx, cancel := context.WithCancel(context.Background())

	var attempts int32
	mockRPC := func(ctx context.Context) (any, error) {
		atomic.AddInt32(&attempts, 1)
		// Cancel parent context during the first attempt's failure
		cancel()
		return nil, errors.New("fail attempt")
	}

	start := time.Now()
	_, err := client.Invoke(ctx, mockRPC)
	elapsed := time.Since(start)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if elapsed >= 150*time.Millisecond {
		t.Errorf("Invoke took %v; should have aborted immediately during backoff without sleeping 200ms", elapsed)
	}
	if atomic.LoadInt32(&attempts) != 1 {
		t.Errorf("expected 1 attempt before cancellation took effect, got %d", attempts)
	}
}

func TestNewResilientClient_Validation(t *testing.T) {
	_, err := NewResilientClient(RetryConfig{MaxAttempts: 0})
	if !errors.Is(err, ErrZeroAttempts) {
		t.Fatalf("expected ErrZeroAttempts, got %v", err)
	}
}
