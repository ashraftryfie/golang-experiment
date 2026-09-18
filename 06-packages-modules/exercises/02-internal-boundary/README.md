# Exercise 02: Internal Security Boundary (Tier 2 - Medium)

## 🎯 Problem Statement
Implement a token service package that uses an `internal/hasher` sub-package for cryptographic hashing:
- `internal/hasher`:
  - Contains hashing logic `Hash(secret, data string) string`.
  - Is strictly unimportable outside this package's parent tree.
- `token.go`:
  - Exposes `GenerateToken(secret, userID string) string`.
  - Exposes `ValidateToken(secret, userID, token string) bool`.

---

## 🛠️ Instructions
1. Implement `Hash` in [`internal/hasher/hasher.go`](./internal/hasher/hasher.go).
2. Implement token logic in [`token.go`](./token.go).
3. Run tests:
   ```powershell
   go test -v ./06-packages-modules/exercises/02-internal-boundary/...
   ```
