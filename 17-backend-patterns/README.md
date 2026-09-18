# Stage 17: Production Backend Patterns

Welcome to **Stage 17** of your Go mastery journey. Building backend systems that survive real-world production incidents requires defensive architecture: handling sudden traffic spikes, preventing cascading failures, and zero-downtime server deployments.

---

## 1. Core Production Patterns

### A. Graceful Shutdown
In containerized environments (Kubernetes, Docker, ECS), when a pod terminates, the runtime sends `SIGTERM` followed by `SIGKILL` after a grace period. A naive server aborts in-flight requests, causing user errors and broken transactions.
- **Idiomatic Pattern**:
  1. Catch OS termination signals via `signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)`.
  2. Call `httpServer.Shutdown(shutdownCtx)`. This closes all listeners and idle connections, while waiting for in-flight requests to complete.
  3. Close database connections and flush background message queues.

### B. Rate Limiting (Token Bucket)
Rate limiters protect downstream resources from overload and brute-force attacks:
- **Capacity ($C$)**: Maximum burst size permitted.
- **Refill Rate ($r$)**: Number of tokens regenerated per second.
- **Token Math**:
  $$\text{tokens} = \min(C, \text{tokens} + \Delta t \times r)$$

### C. Circuit Breakers
When an external dependency (payment gateway, third-party API, database) fails or becomes slow, continuing to send traffic wastes threads, exhausts connection pools, and can bring down your own service.
```
  ┌───────────┐       Failure Threshold Reached       ┌──────────┐
  │  CLOSED   │ ────────────────────────────────────> │   OPEN   │
  │ (Normal)  │                                       │ (Fails)  │
  └───────────┘                                       └──────────┘
        ▲                                                   │
        │ Success Threshold Reached                         │ Timeout Elapsed
        │                                                   ▼
  ┌───────────┐                                       ┌──────────┐
  │           │ <──────────────────────────────────── │HALF-OPEN │
  └───────────┘             Probe Fails               │ (Canary) │
                                                      └──────────┘
```

---

## 2. Directory Layout

```
17-backend-patterns/
├── README.md
├── examples/
│   ├── graceful_shutdown/main.go  # Server shutdown on OS interrupt
│   └── rate_limiter/main.go       # Token bucket rate limiter demonstration
├── exercises/
│   ├── 01-token-bucket            # Thread-safe token bucket rate limiter
│   ├── 02-circuit-breaker         # 3-state circuit breaker state machine
│   └── 03-graceful-server         # Lifecycle-managed HTTP server with timeout
├── solutions/
│   ├── 01-token-bucket
│   ├── 02-circuit-breaker
│   └── 03-graceful-server
└── challenges/
    └── resilient-gateway/         # Rate-limited, circuit-broken API gateway
```
