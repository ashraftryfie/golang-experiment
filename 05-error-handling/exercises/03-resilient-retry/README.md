# Exercise 03: Resilient Retry & Error Classification (Tier 3 - Hard)

## 🎯 Problem Statement
Implement a retry engine that distinguishes between **transient** (recoverable) errors and **permanent** (fatal) errors:

### Error Interface
```go
type TransientError interface {
    IsTransient() bool
}
```

### Methods to Implement
1. `type NetworkError struct { Msg string }`: Implements `TransientError` returning `true`.
2. `type AuthError struct { Msg string }`: Implements `TransientError` returning `false`.
3. `Retry(maxAttempts int, initialDelay time.Duration, fn func() error) error`:
   - Executes `fn()`.
   - If `fn()` returns `nil`, returns immediately.
   - If `fn()` returns an error:
     - Checks if error implements `TransientError` and `IsTransient() == false`. If permanent, returns immediately without retrying!
     - If transient (or standard unclassified error), sleeps with exponential backoff (`delay * 2`) and retries up to `maxAttempts`.
   - Returns final error if all attempts fail.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement types and `Retry`.
3. Run tests:
   ```powershell
   go test -v ./05-error-handling/exercises/03-resilient-retry
   ```
