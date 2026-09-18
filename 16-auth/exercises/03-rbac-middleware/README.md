# Exercise 03: Role-Based Access Control (RBAC) HTTP Middleware

## Objective
Implement a production HTTP middleware in Go that inspects JWT Bearer tokens, injects authenticated user identities into the request context, and enforces role boundaries.

## Requirements
1. **Authorization Header Extraction**:
   - Extract `Authorization: Bearer <token>`.
   - If missing or malformed, return HTTP `401 Unauthorized`.
2. **Token Verification**:
   - Verify HMAC-SHA256 signature and expiration.
   - If invalid or expired, return HTTP `401 Unauthorized`.
3. **Role Enforcement**:
   - If `allowedRoles` are specified, verify that `claims.Role` matches at least one allowed role.
   - If user lacks required role, return HTTP `403 Forbidden`.
4. **Context Injection**:
   - Store user claims into `r.Context()` using a private unexported key type.
   - Provide helper `GetUserClaims(ctx context.Context) (*UserClaims, bool)`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./16-auth/exercises/03-rbac-middleware`
