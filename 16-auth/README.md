# Stage 16: Authentication & Security

Welcome to **Stage 16** of your Go mastery journey. Production backend services cannot rely on naive security assumptions. This stage equips you with rock-solid security engineering in Go: adaptive password hashing, cryptographic token verification, timing-attack resistance, and role-based access control (RBAC).

---

## 1. Core Security Principles in Go

### A. Never Use Fast Hashes for Passwords
- Fast hashes (MD5, SHA-1, SHA-256) are engineered for high-throughput integrity checks. They can be cracked at billions of hashes per second using GPUs and ASICs.
- **Adaptive Key Derivation**: Use `golang.org/x/crypto/bcrypt` or Argon2. `bcrypt` incorporates automatic per-password salting and an adjustable **work factor / cost** (e.g. `cost = 12`) that makes brute-force attacks computationally infeasible.

### B. Cryptographic Randomness
- Never use `math/rand` for security tokens, session IDs, or crypto keys; its pseudo-random sequence is deterministic and predictable.
- Always use `crypto/rand`:
  ```go
  tokenBytes := make([]byte, 32)
  if _, err := rand.Read(tokenBytes); err != nil { ... }
  ```

### C. Timing Attacks & Constant-Time Comparison
- Standard string equality (`a == b`) short-circuits at the first differing byte. An attacker measuring latency over millions of requests can deduce secrets byte-by-byte.
- Always compare HMAC signatures and security tokens with `crypto/subtle.ConstantTimeCompare`:
  ```go
  if subtle.ConstantTimeCompare([]byte(sigA), []byte(sigB)) != 1 {
      return ErrInvalidSignature
  }
  ```

---

## 2. JWT (JSON Web Token) HMAC Architecture

A JWT consists of three Base64URL-encoded components joined by periods:
```
Header.Payload.Signature
```
1. **Header**: Declares algorithm (`{"alg":"HS256","typ":"JWT"}`).
2. **Payload / Claims**: Request identity, expiration (`exp`), and roles (`{"sub":"user_123","role":"admin","exp":1700000000}`).
3. **Signature**: `HMAC-SHA256(Base64URL(Header) + "." + Base64URL(Payload), secretKey)`.

---

## 3. Directory Layout

```
16-auth/
├── README.md
├── examples/
│   ├── password_hashing/main.go  # Adaptive bcrypt hashing & verification
│   └── jwt_token/main.go         # Standard library HMAC-SHA256 JWT generation
├── exercises/
│   ├── 01-password-hasher        # Robust bcrypt password service
│   ├── 02-jwt-auth               # Type-safe HMAC-SHA256 JWT token parser
│   └── 03-rbac-middleware        # Context-injected RBAC HTTP middleware
├── solutions/
│   ├── 01-password-hasher
│   ├── 02-jwt-auth
│   └── 03-rbac-middleware
└── challenges/
    └── secure-auth-service/      # Full register/login/RBAC microservice
```
