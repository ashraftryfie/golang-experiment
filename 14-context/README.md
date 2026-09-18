# Stage 14: Context & Cancellation

Welcome to **Stage 14** of your Go mastery journey. In high-performance backend systems, distributed services, and HTTP microservices, handling timeouts, cancellations, and request lifecycles without leaking resources or goroutines is a core responsibility. The standard library `context` package provides the unified abstraction for this in Go.

---

## 1. Mental Model: Context as a Cancellation Tree

A `context.Context` represents a immutable node in a directed acyclic cancellation tree:

```
                  context.Background() (Root)
                           │
             ┌─────────────┴─────────────┐
             ▼                           ▼
   ctxCancel, cancel()          ctxTimeout (500ms)
             │                           │
      ┌──────┴──────┐                    ▼
      ▼             ▼               Child Worker
  Database Query   gRPC Call
```

### Key Properties:
1. **Immutable & Derived**: Contexts are never mutated in-place. You derive child contexts using `context.WithCancel`, `context.WithTimeout`, `context.WithDeadline`, or `context.WithValue`.
2. **Top-Down Cancellation**: When a parent context is cancelled, all descendant contexts derived from it are cancelled immediately. A child cannot cancel its parent.
3. **Goroutine-Safe**: The same `Context` can be passed to arbitrarily many goroutines simultaneously; all methods on `Context` are safe for concurrent use.

---

## 2. The `context.Context` Interface

```go
type Context interface {
    Deadline() (deadline time.Time, ok bool)
    Done() <-chan struct{}
    Err() error
    Value(key any) any
}
```

- **`Done() <-chan struct{}`**: Returns a channel that is closed when the context is cancelled or times out. Receiving from a closed channel unblocks immediately (`case <-ctx.Done():`).
- **`Err() error`**: Returns `nil` if still active; returns `context.Canceled` or `context.DeadlineExceeded` once `Done()` is closed.
- **`Deadline()`**: Returns the clock time when work on behalf of this context should stop.
- **`Value(key)`**: Retrieves request-scoped values associated with this context.

---

## 3. Idiomatic Rules of `context`

1. **First Parameter**: Always pass `ctx context.Context` as the very first argument to any function or method performing I/O or blocking operations:
   ```go
   func QueryUser(ctx context.Context, id string) (*User, error)
   ```
2. **Never Store in Structs**: Do not store a `Context` inside a struct type. Pass it explicitly through the call stack. (The only exception is HTTP Request structs in standard library adapters).
3. **Always Call Cancel**: Every call to `WithCancel`, `WithTimeout`, or `WithDeadline` returns a `cancel` function. Always call it, typically via `defer cancel()`:
   ```go
   ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
   defer cancel() // Releases timers and detaches child from parent
   ```
4. **Context Values**: Use `context.WithValue` strictly for **request-scoped transit data** (e.g. distributed trace IDs, authentication tokens, client IP). Never pass optional function parameters or database connections through `Context.Value`.
5. **Key Types**: Always use an unexported custom type for context keys to eliminate package collisions:
   ```go
   type contextKey struct{}
   var requestIDKey = contextKey{}
   ```

---

## 4. Stage Directory Structure

```
14-context/
├── README.md
├── examples/
│   ├── cancellation_tree/main.go  # Demonstrates parent-to-child cancellation
│   └── http_timeout/main.go       # Demonstrates HTTP request timeout with context
├── exercises/
│   ├── 01-timeout-fetcher/        # Parallel fetcher bounded by timeout
│   ├── 02-context-propagation/    # Pipeline aborting on context cancellation
│   └── 03-request-scoped-tracer/  # Safe request-scoped tracing with custom keys
├── solutions/
│   ├── 01-timeout-fetcher/
│   ├── 02-context-propagation/
│   └── 03-request-scoped-tracer/
└── challenges/
    └── resilient-rpc-client/      # Retrying RPC client with backoff and deadline
```
