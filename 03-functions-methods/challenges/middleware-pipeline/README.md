# Challenge: Functional Middleware Pipeline

## 🎯 Goal
Build a composable, functional middleware chaining pipeline modeled after Go's production HTTP handler patterns.

---

## 📋 Architecture & Types

```go
// Request represents an incoming message context.
type Request struct {
    Path      string
    AuthToken string
    Body      string
}

// Response represents the output of processing.
type Response struct {
    StatusCode int
    Body       string
}

// HandlerFunc processes a Request and produces a Response.
type HandlerFunc func(req Request) Response

// Middleware wraps a HandlerFunc with pre- and post-processing logic.
type Middleware func(next HandlerFunc) HandlerFunc
```

### Middlewares to Build
1. `LoggingMiddleware(logger func(string)) Middleware`:
   - Logs entry path and response status code.
2. `AuthMiddleware(validToken string) Middleware`:
   - If `req.AuthToken != validToken`, intercepts and returns `401 Unauthorized` without calling `next`.
3. `RecoveryMiddleware(errLogger func(any)) Middleware`:
   - Uses `defer` and `recover()` to catch panics in downstream handlers, returning `500 Internal Server Error`.
4. `Chain(middlewares ...Middleware) func(HandlerFunc) HandlerFunc`:
   - Composes all middlewares so they execute in standard FIFO wrapping order.

---

## 🧪 Testing
```powershell
go test -v ./03-functions-methods/challenges/middleware-pipeline
```
