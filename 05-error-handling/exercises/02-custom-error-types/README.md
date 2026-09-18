# Exercise 02: Custom Validation Errors & `errors.As` (Tier 2 - Medium)

## 🎯 Problem Statement
Implement structured validation error types and inspect them using `errors.As()`:

### Struct Definition
```go
type ValidationError struct {
    Field   string
    Message string
    Value   any
}

func (e *ValidationError) Error() string
```

### Methods to Implement
`ValidateRegistration(username, email string, age int) error`:
- If `username` is shorter than 3 characters, returns `*ValidationError` for field `"username"`.
- If `email` does not contain `@`, returns `*ValidationError` for field `"email"`.
- If `age < 18`, returns `*ValidationError` for field `"age"`.
- If all valid, returns `nil`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `Error()` on `*ValidationError` and `ValidateRegistration`.
3. Run tests:
   ```powershell
   go test -v ./05-error-handling/exercises/02-custom-error-types
   ```
