package tokenbucket

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestTokenBucket_BurstCapacity(t *testing.T) {
	// Rate: 1 token/sec, Capacity: 3
	tb, err := NewTokenBucket(1.0, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First 3 should succeed immediately
	for i := 1; i <= 3; i++ {
		if !tb.Allow() {
			t.Fatalf("expected request %d to be allowed", i)
		}
	}

	// 4th request must be denied
	if tb.Allow() {
		t.Fatalf("expected 4th request to be denied after exhausting burst capacity")
	}
}

func TestTokenBucket_Refill(t *testing.T) {
	// Rate: 10 tokens/sec, Capacity: 2
	tb, _ := NewTokenBucket(10.0, 2)

	// Consume all 2
	if !tb.AllowN(2) {
		t.Fatalf("expected AllowN(2) to succeed")
	}
	if tb.Allow() {
		t.Fatalf("expected bucket to be empty")
	}

	// Wait 250ms -> should regenerate at least 2 tokens
	time.Sleep(250 * time.Millisecond)

	if !tb.Allow() {
		t.Fatalf("expected token to be regenerated after sleep")
	}
}

func TestTokenBucket_ConcurrentAccess(t *testing.T) {
	// Capacity: 50, Rate: 100/sec
	tb, _ := NewTokenBucket(100.0, 50)

	var allowedCount int32
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if tb.Allow() {
				atomic.AddInt32(&allowedCount, 1)
			}
		}()
	}

	wg.Wait()

	if atomic.LoadInt32(&allowedCount) != 50 {
		t.Fatalf("expected all 50 concurrent requests to consume initial burst, got %d", allowedCount)
	}
}

func TestTokenBucket_InvalidConfig(t *testing.T) {
	_, err := NewTokenBucket(0, 10)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}

	_, err = NewTokenBucket(10, -1)
	if !errors.Is(err, ErrInvalidConfig) {
		t.Fatalf("expected ErrInvalidConfig, got %v", err)
	}
}
