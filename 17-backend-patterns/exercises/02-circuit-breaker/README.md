# Exercise 02: Resilient Circuit Breaker State Machine

## Objective
Implement a production circuit breaker state machine that isolates failing downstream services, fails fast when open, and safely self-heals using a half-open canary phase.

## Requirements
1. **States**:
   - `StateClosed`: Normal operation. Invocations execute normally. Consecutive errors increment failure counter. If failures reach `MaxFailures`, transition to `StateOpen`.
   - `StateOpen`: Fails fast immediately returning `ErrCircuitOpen`. If `Timeout` elapses, next call transitions to `StateHalfOpen`.
   - `StateHalfOpen`: Allows canary probe calls. If a call fails, immediately trips back to `StateOpen`. If calls succeed up to `SuccessThreshold`, heals back to `StateClosed` and resets failures.
2. **Thread Safety**: All state transitions and counter increments must be synchronized with `sync.Mutex`.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./17-backend-patterns/exercises/02-circuit-breaker`
