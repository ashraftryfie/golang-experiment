# Exercise 02: Sequential Pipeline with Context Propagation

## Objective
Implement a robust execution pipeline that propagates `context.Context` through a sequence of functions, ensuring that cancellation signals instantly halt execution and prevent subsequent operations from running.

## Requirements
1. **Pipeline Execution**:
   - Signature: `ExecutePipeline(ctx context.Context, steps []StepFunc) error`
   - `type StepFunc func(ctx context.Context) error`
   - Check `ctx.Err()` before invoking each step. If canceled, return `ctx.Err()`.
   - Pass `ctx` directly to each `StepFunc`.
   - If any step fails, return the error immediately without invoking subsequent steps.
2. **Timeout Wrapper**:
   - Signature: `RunWithTimeout(ctx context.Context, timeout time.Duration, steps []StepFunc) error`
   - Derives a timeout context and delegates to `ExecutePipeline`, guaranteeing cleanup via `defer cancel()`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./14-context/exercises/02-context-propagation`
