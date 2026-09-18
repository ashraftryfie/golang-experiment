package main

import (
	"fmt"
	"sync"
	"time"
)

type SimpleTokenBucket struct {
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	lastRefill time.Time
	mu         sync.Mutex
}

func NewSimpleTokenBucket(rate float64, capacity int) *SimpleTokenBucket {
	return &SimpleTokenBucket{
		capacity:   float64(capacity),
		tokens:     float64(capacity),
		refillRate: rate,
		lastRefill: time.Now(),
	}
}

func (tb *SimpleTokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	// Refill tokens
	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

func main() {
	// Rate = 5 tokens/sec, Capacity = 3
	limiter := NewSimpleTokenBucket(5.0, 3)

	fmt.Println("Attempting 5 immediate requests (Capacity = 3):")
	for i := 1; i <= 5; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED\n", i)
		} else {
			fmt.Printf("Request %d: RATE LIMITED (DENIED)\n", i)
		}
	}

	fmt.Println("\nSleeping 400ms (should regenerate ~2 tokens)...")
	time.Sleep(400 * time.Millisecond)

	for i := 6; i <= 8; i++ {
		if limiter.Allow() {
			fmt.Printf("Request %d: ALLOWED\n", i)
		} else {
			fmt.Printf("Request %d: RATE LIMITED (DENIED)\n", i)
		}
	}
}
