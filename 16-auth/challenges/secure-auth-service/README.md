# Engineering Challenge: Production Auth & RBAC Microservice

## Challenge Overview
Authentication and access control form the critical security foundation of every enterprise SaaS and cloud application. Subtle flaws—like using weak password hashing, ignoring timing attacks, leaking error details during login, or mishandling token expiration—lead directly to critical vulnerabilities.

You must implement an end-to-end **Production-Grade Auth & RBAC Microservice** in Go featuring:
1. **Secure Registration (`POST /auth/register`)**:
   - Validates email format and enforces minimum password complexity.
   - Prevents duplicate account registrations.
   - Hashes passwords with adaptive `bcrypt` salt and work factor.
2. **Timing-Safe Login (`POST /auth/login`)**:
   - Compares credentials using `bcrypt.CompareHashAndPassword`.
   - Issues HMAC-SHA256 signed JSON Web Tokens (JWT) with an unforgeable signature.
   - Masks authentication failure reasons to prevent user enumeration attacks (`{"error":"invalid email or password"}`).
3. **Context-Driven RBAC Protection**:
   - `GET /api/profile`: Open to any authenticated user. Extracts caller identity from request context.
   - `GET /api/admin/metrics`: Restricted strictly to users with the `admin` role (returns `403 Forbidden` for non-admin accounts).

---

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./16-auth/challenges/secure-auth-service/...
```
