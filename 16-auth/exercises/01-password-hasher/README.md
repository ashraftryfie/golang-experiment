# Exercise 01: Adaptive Password Hasher with Bcrypt

## Objective
Implement a secure password management service in Go using `golang.org/x/crypto/bcrypt`.

## Requirements
1. **`HashPassword(password string, cost int) (string, error)`**:
   - Rejects empty passwords with `ErrEmptyPassword`.
   - If cost is out of range (`< bcrypt.MinCost` or `> bcrypt.MaxCost`), fallback safely to `bcrypt.DefaultCost`.
   - Returns string hash.
2. **`VerifyPassword(hashedPassword, password string) bool`**:
   - Compares raw password with the hash; returns `true` if valid, `false` otherwise.
3. **`NeedsRehash(hashedPassword string, targetCost int) bool`**:
   - Inspects the existing hash's work factor with `bcrypt.Cost()`.
   - Returns `true` if the existing cost is strictly less than `targetCost`, allowing seamless hash upgrades on user login.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./16-auth/exercises/01-password-hasher`
