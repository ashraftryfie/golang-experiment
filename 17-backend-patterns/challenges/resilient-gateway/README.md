# Engineering Challenge: Resilient API Gateway

## Challenge Overview
Modern cloud microservices face unpredictable upstream outages and malicious traffic spikes. An API Gateway must act as a protective bulkhead, preventing bad traffic from destabilizing internal services while isolating downstream callers from cascading failures.

You must build a production-grade **Resilient API Gateway** combining:
1. **Per-Client Token-Bucket Rate Limiting**:
   - Identifies clients via `X-Client-ID` header (fallback to remote address IP).
   - Enforces burst capacity and smooth token refill.
   - Rejects throttled requests with `429 Too Many Requests` and a `Retry-After` header.
2. **Upstream Circuit Breaker**:
   - Monitors upstream HTTP responses.
   - Trips to `Open` after consecutive 5xx errors.
   - When `Open`, fails fast immediately with `503 Service Unavailable`, preventing useless network roundtrips.
   - Automatically probes upstream health after a cooldown timeout via a `Half-Open` canary.
3. **Graceful Zero-Downtime Server Lifecycle**:
   - Supports non-blocking `Start(addr string)` and bounded `Shutdown(timeout time.Duration)`.
   - Drains active proxy requests before releasing port bindings.

---

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./17-backend-patterns/challenges/resilient-gateway/...
```
