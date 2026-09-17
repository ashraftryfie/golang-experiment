package ratelimiter

import (
	"testing"
	"time"
)

func TestSolutionRateLimiter(t *testing.T) {
	limiter := NewRateLimiter(2, 40*time.Millisecond)

	if !limiter() {
		t.Fatalf("first request should pass")
	}
	if !limiter() {
		t.Fatalf("second request should pass")
	}
	if limiter() {
		t.Fatalf("third request in window should be blocked")
	}

	time.Sleep(50 * time.Millisecond)

	if !limiter() {
		t.Fatalf("request after cooldown should pass")
	}
}
