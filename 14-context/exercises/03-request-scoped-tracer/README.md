# Exercise 03: Request-Scoped Tracer with Context Values

## Objective
Implement type-safe, collision-free request metadata and trace propagation using `context.WithValue`.

## Requirements
1. **Unexported Key Type**: Prevent namespace collisions by using a private type struct for the context key (e.g., `type traceContextKey struct{}`).
2. **Type-Safe Storage & Retrieval**:
   - `WithTraceInfo(ctx context.Context, info TraceInfo) context.Context`
   - `GetTraceInfo(ctx context.Context) (TraceInfo, bool)`
3. **Immutable Tag Enrichment**:
   - `AddTraceTag(ctx context.Context, key, value string) context.Context`
   - Contexts must remain immutable. Creating a derived context with a new tag MUST clone the existing tag map so concurrent goroutines sharing the parent context never suffer data races.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./14-context/exercises/03-request-scoped-tracer`
