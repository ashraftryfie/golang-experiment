# Exercise 01: Thread-Safe Token Bucket Rate Limiter

## Objective
Implement a thread-safe token bucket algorithm that supports burst capacity and smooth token replenishment.

## Requirements
1. **`NewTokenBucket(refillRate float64, capacity int) (*TokenBucket, error)`**:
   - `refillRate`: Tokens added per second (must be > 0).
   - `capacity`: Maximum bucket size (must be > 0).
   - Reject non-positive values with `ErrInvalidConfig`.
   - Starts fully saturated with tokens equal to `capacity`.
2. **`Allow() bool`**:
   - Equivalent to `AllowN(1)`.
3. **`AllowN(n int) bool`**:
   - Updates token count according to time elapsed since last refill.
   - Caps tokens at `capacity`.
   - If at least `n` tokens are available, deducts `n` and returns `true`. Otherwise returns `false`.
   - Thread-safe across concurrent goroutines via `sync.Mutex`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./17-backend-patterns/exercises/01-token-bucket`
