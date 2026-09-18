# Engineering Challenge: Full-Stack Telemetry Middleware & Metrics Service

## Challenge Overview
Observability is not an afterthought added via fragmented print statements. In production microservice meshes, every HTTP request must automatically carry distributed trace context, increment Prometheus operational metrics, and emit sanitized structured JSON logs.

You must implement a production-grade **Full-Stack Telemetry Middleware Suite** in Go that provides:
1. **W3C Distributed Trace Injection**:
   - Inspects incoming `traceparent` header.
   - If present, preserves the caller's `TraceID`.
   - If missing, generates a fresh W3C TraceContext and sets `traceparent` in the response header.
2. **Prometheus Operational Metrics**:
   - Serves standard Prometheus text exposition format at `GET /metrics`.
   - Tracks `http_requests_total{method, status}` counter.
   - Tracks `http_active_requests` gauge (incremented upon entry, decremented upon completion).
3. **Structured Contextual Logging (`slog`)**:
   - Emits a structured JSON log entry for every completed request recording:
     `trace_id`, `method`, `path`, `status`, `duration_ms`, and `remote_ip`.

---

## Running Tests
Run tests with race detection:
```bash
go test -v -race ./19-observability/challenges/observability-middleware/...
```
