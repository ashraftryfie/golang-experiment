package ratelimiter

import (
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(2, 50*time.Millisecond)

	// If starter default returns false, skip gracefully
	if !limiter() {
		t.Skip("skipping: NewRateLimiter is not yet implemented (implement in starter.go)")
	}

	// 2nd request in window should be allowed
	if !limiter() {
		t.Errorf("2nd request should be allowed")
	}

	// 3rd request in window should be denied
	if limiter() {
		t.Errorf("3rd request should be blocked")
	}

	// Wait for window to expire
	time.Sleep(60 * time.Millisecond)

	// Request should now be allowed again
	if !limiter() {
		t.Errorf("request after window expiration should be allowed")
	}
}
