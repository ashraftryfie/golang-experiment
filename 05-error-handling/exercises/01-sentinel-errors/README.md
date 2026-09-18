# Exercise 01: Sentinel Errors & Repository Lookup (Tier 1 - Easy)

## 🎯 Problem Statement
Implement an in-memory user repository that uses sentinel errors and `%w` wrapping:

### Sentinel Errors
```go
var (
    ErrUserNotFound    = errors.New("user not found")
    ErrDuplicateEmail  = errors.New("email address already registered")
    ErrInvalidUsername = errors.New("username cannot be empty")
)
```

### Methods to Implement
1. `(r *UserRepository) CreateUser(id int, username, email string) error`:
   - Validates `username != ""` (`ErrInvalidUsername`).
   - Validates email is not already used (`ErrDuplicateEmail`).
   - Stores user.
2. `(r *UserRepository) FindByID(id int) (*User, error)`:
   - Returns user if found.
   - If missing, wraps and returns: `fmt.Errorf("user lookup id %d: %w", id, ErrUserNotFound)`.

---

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement functions and wrap errors using `%w`.
3. Run tests:
   ```powershell
   go test -v ./05-error-handling/exercises/01-sentinel-errors
   ```
