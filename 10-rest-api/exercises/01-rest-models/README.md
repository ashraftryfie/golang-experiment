# Exercise 01: REST Request Validation (Tier 1 - Easy)

## 🎯 Problem Statement
In production REST APIs, client input must be rigorously validated before reaching business logic or database layers.

Implement request validation for `CreateUserRequest`:
```go
type ValidationError struct {
    Field   string `json:"field"`
    Message string `json:"message"`
}

type CreateUserRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Age      int    `json:"age"`
}

func (r *CreateUserRequest) Validate() []ValidationError
```

### Validation Rules
1. `Username`: Must be between 3 and 30 characters (inclusive) and must not contain spaces.
   - If invalid: field `"username"`, message `"username must be 3-30 characters without whitespace"`.
2. `Email`: Must contain both `@` and `.`.
   - If invalid: field `"email"`, message `"must be a valid email address"`.
3. `Age`: Must be between 18 and 120 (inclusive).
   - If invalid: field `"age"`, message `"age must be between 18 and 120"`.
4. If all fields are valid, return an empty or nil slice.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `Validate`.
3. Run tests:
   ```powershell
   go test -v ./10-rest-api/exercises/01-rest-models
   ```
