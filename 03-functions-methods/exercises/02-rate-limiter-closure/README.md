# Exercise 02: Rate Limiter Closure (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a stateful sliding-window rate limiter using a Go closure.

### Requirements
Implement `NewRateLimiter(maxRequests int, windowDuration time.Duration) func() bool`:
- Returns a closure function `Allow() bool`.
- Each invocation of `Allow()` records a request timestamp.
- If the number of requests in the most recent `windowDuration` is strictly less than `maxRequests`, `Allow()` records the request and returns `true`.
- If the limit has been reached within the window, `Allow()` does NOT record a new request and returns `false`.
- If `maxRequests <= 0` or `windowDuration <= 0`, return a closure that always returns `false`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Use a slice of `time.Time` inside the closure scope to track timestamps.
3. On each call, prune timestamps older than `now - windowDuration`.
4. Run tests:
   ```powershell
   go test -v ./03-functions-methods/exercises/02-rate-limiter-closure
   ```
