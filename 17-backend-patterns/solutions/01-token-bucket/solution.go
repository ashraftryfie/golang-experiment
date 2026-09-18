package tokenbucket

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrInvalidConfig = errors.New("refillRate and capacity must be greater than zero")
)

type TokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(refillRate float64, capacity int) (*TokenBucket, error) {
	if refillRate <= 0 || capacity <= 0 {
		return nil, ErrInvalidConfig
	}

	return &TokenBucket{
		capacity:   float64(capacity),
		tokens:     float64(capacity),
		refillRate: refillRate,
		lastRefill: time.Now(),
	}, nil
}

func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucket) AllowN(n int) bool {
	if n <= 0 {
		return false
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()

	cost := float64(n)
	if tb.tokens >= cost {
		tb.tokens -= cost
		return true
	}
	return false
}

func (tb *TokenBucket) AvailableTokens() float64 {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.refill()
	return tb.tokens
}

// refill updates the token count based on elapsed time. Caller must hold tb.mu.
func (tb *TokenBucket) refill() {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}
}
