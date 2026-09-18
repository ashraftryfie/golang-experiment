# Exercise 03: Graceful HTTP Server Lifecycle

## Objective
Implement a production HTTP server wrapper in Go that manages listener lifecycle, drains in-flight requests, coordinates background tasks, and performs graceful shutdown bounded by a timeout.

## Requirements
1. **`Start() error`**:
   - Creates a TCP listener and begins serving HTTP requests in a non-blocking background goroutine.
2. **`TrackTask(fn func(ctx context.Context))`**:
   - Launches a background worker goroutine tracked by internal synchronization.
   - Passes a lifecycle context that is canceled when shutdown begins.
3. **`Shutdown(timeout time.Duration) error`**:
   - Derives a `context.WithTimeout`.
   - Signals background tasks to exit by canceling the lifecycle context.
   - Shuts down the HTTP listener and drains in-flight requests via `httpServer.Shutdown(ctx)`.
   - Waits for all background tasks to terminate cleanly.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./17-backend-patterns/exercises/03-graceful-server`
