# Exercise 01: Bounded Parallel Fetcher with Timeout

## Objective
Implement a concurrent URL fetcher that derives a child context with a timeout, invokes a fetcher function across multiple URLs in parallel, and returns aggregated results.

## Requirements
1. **Derive Timeout Context**: Create a child context using `context.WithTimeout(ctx, timeout)` and ensure `cancel()` is deferred.
2. **Concurrent Execution**: Run fetches concurrently using goroutines.
3. **Cancellation Awareness**: If the derived context times out or parent context is canceled, any in-flight workers must terminate when observing `ctx.Done()`.
4. **Order Preservation**: The returned slice `[]FetchResult` should correspond to the input `urls` order.
5. **No Goroutine Leaks**: Ensure all spawned goroutines terminate and results can be received without deadlocking.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./14-context/exercises/01-timeout-fetcher`
