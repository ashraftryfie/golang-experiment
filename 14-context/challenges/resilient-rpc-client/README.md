# Engineering Challenge: Resilient Context-Aware RPC Client

## Challenge Overview
In distributed cloud architectures, network calls encounter transient spikes, node failovers, and slow queries. A naive retry loop can easily deadlock or amplify cascades (the "thundering herd" problem) if it fails to respect parent deadlines or leaks timeout timers.

You must implement a production-grade **Resilient RPC Invocation Client** in Go that combines:
1. **Parent Context Propagation**: If the incoming request context expires or cancels, all work and waiting must abort immediately.
2. **Per-Attempt Timeout**: Each individual RPC invocation is strictly bounded by a child `context.WithTimeout`, preventing an unresponsive downstream service from stalling the entire retry budget.
3. **Exponential Backoff**: Between retries, backoff duration scales exponentially up to a configured cap.
4. **Context-Aware Sleep**: Waiting during backoff MUST be interruptible via `select` on `timer.C` and `ctx.Done()`, guaranteeing zero delay on parent cancellation and no leaked timers.
5. **Clean Resource Deferrals**: All derived cancellation functions (`cancel()`) and timers must be reclaimed immediately.

---

## Architecture

```
Incoming Request (ctx, Total Deadline = 500ms)
   │
   ├─ Attempt 1 (attemptCtx = 50ms) -> Downstream Fails / Times out
   │    └─ Backoff Sleep (10ms) [select on timer vs ctx.Done()]
   │
   ├─ Attempt 2 (attemptCtx = 50ms) -> Downstream Fails / Times out
   │    └─ Backoff Sleep (20ms) [select on timer vs ctx.Done()]
   │
   ├─ Attempt 3 (attemptCtx = 50ms) -> Downstream Succeeds! -> Return Result
```

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./14-context/challenges/resilient-rpc-client/...
```
