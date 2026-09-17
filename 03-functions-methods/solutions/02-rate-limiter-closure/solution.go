package ratelimiter

import (
	"sync"
	"time"
)

// NewRateLimiter creates a sliding window rate limiter closure.
func NewRateLimiter(maxRequests int, windowDuration time.Duration) func() bool {
	if maxRequests <= 0 || windowDuration <= 0 {
		return func() bool { return false }
	}

	var mu sync.Mutex
	timestamps := make([]time.Time, 0, maxRequests)

	return func() bool {
		mu.Lock()
		defer mu.Unlock()

		now := time.Now()
		cutoff := now.Add(-windowDuration)

		// Prune expired timestamps in-place
		validIdx := 0
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				timestamps[validIdx] = ts
				validIdx++
			}
		}
		timestamps = timestamps[:validIdx]

		if len(timestamps) < maxRequests {
			timestamps = append(timestamps, now)
			return true
		}

		return false
	}
}
