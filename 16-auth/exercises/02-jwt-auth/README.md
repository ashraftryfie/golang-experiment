# Exercise 02: Type-Safe HMAC-SHA256 JWT Parser

## Objective
Implement a lightweight, dependency-free JSON Web Token (JWT) generator and verifier using Go's standard library `crypto/hmac`, `crypto/sha256`, and `crypto/subtle`.

## Requirements
1. **Claims**:
   - `Subject` (`sub`): User identifier (must not be empty).
   - `Role` (`role`): User authorization role (e.g. `admin`, `user`).
   - `Expires` (`exp`): Unix epoch expiration timestamp.
2. **`CreateToken(claims Claims, secret []byte) (string, error)`**:
   - Enforce secret key length of at least 16 bytes (`ErrSecretTooShort`).
   - Base64URL-encode header (`{"alg":"HS256","typ":"JWT"}`) and payload.
   - Compute HMAC-SHA256 signature and return `header.payload.signature`.
3. **`VerifyToken(tokenString string, secret []byte) (*Claims, error)`**:
   - Split into 3 parts.
   - Compare signature using `subtle.ConstantTimeCompare` against timing attacks.
   - Return `ErrInvalidSignature`, `ErrTokenExpired`, or `ErrInvalidToken`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./16-auth/exercises/02-jwt-auth`
