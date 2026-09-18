package tokenbucket

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrNotImplemented = errors.New("exercise not implemented yet")
	ErrInvalidConfig  = errors.New("refillRate and capacity must be greater than zero")
)

type TokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewTokenBucket(refillRate float64, capacity int) (*TokenBucket, error) {
	return nil, ErrNotImplemented
}

func (tb *TokenBucket) Allow() bool {
	return tb.AllowN(1)
}

func (tb *TokenBucket) AllowN(n int) bool {
	return false
}

func (tb *TokenBucket) AvailableTokens() float64 {
	return 0
}
