package ratelimiter

import (
	"time"
)

// NewRateLimiter creates a sliding window rate limiter closure.
func NewRateLimiter(maxRequests int, windowDuration time.Duration) func() bool {
	// TODO: Return a closure holding timestamps of recent requests
	return func() bool {
		return false
	}
}
