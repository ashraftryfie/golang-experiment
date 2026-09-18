# Exercise 02: API Key Authentication Middleware (Tier 2 - Medium)

## 🎯 Problem Statement
In Go web development, cross-cutting concerns (authentication, rate limiting, logging) are implemented as decorator functions adhering to:
```go
type Middleware func(http.Handler) http.Handler
```

Implement `RequireAPIKey`:
```go
func RequireAPIKey(expectedKey string) Middleware
```

### Requirements
1. Inspect the incoming request header: `r.Header.Get("X-API-Key")`.
2. If the header is missing or does not equal `expectedKey`:
   - Set header `Content-Type: application/json`.
   - Write HTTP status `http.StatusUnauthorized` (401).
   - Write payload: `{"error":"unauthorized"}`.
   - Do NOT invoke the next handler in the chain.
3. If the header matches, invoke `next.ServeHTTP(w, r)`.

## 🛠️ Instructions
1. Open [`starter.go`](./starter.go).
2. Implement `RequireAPIKey`.
3. Run tests:
   ```powershell
   go test -v ./09-http/exercises/02-middleware-pipeline
   ```
