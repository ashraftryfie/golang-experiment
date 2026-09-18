# Challenge: Aggregate Multi-Error Collector

## 🎯 Goal
Implement a composite `MultiError` type that collects multiple independent errors during batch processing or validation passes, implementing both `error` and Go 1.20+ `Unwrap() []error`.

---

## 📋 Architecture & Types

```go
type MultiError struct {
    errors []error
}
```

### Requirements
1. `(m *MultiError) Add(err error)`:
   - Appends non-nil errors to the internal list. Ignores `nil` errors.
2. `(m *MultiError) ErrorOrNil() error`:
   - If empty, returns `nil`! (Critical for avoiding the typed-nil interface trap).
   - If non-empty, returns `m`.
3. `(m *MultiError) Error() string`:
   - Returns a formatted summary of all accumulated errors (e.g. `"encountered 2 errors:\n - error 1\n - error 2"`).
4. `(m *MultiError) Unwrap() []error`:
   - Enables standard `errors.Is` and `errors.As` to search through *all* contained errors simultaneously.

---

## 🧪 Testing
```powershell
go test -v ./05-error-handling/challenges/multierror-collector
```
