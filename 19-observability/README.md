# Stage 19: Observability & Telemetry

Welcome to **Stage 19** of your Go mastery journey. Operating production microservices requires the "Three Pillars of Observability": **Structured Logging**, **Metrics**, and **Distributed Tracing**. In Go 1.21+, the standard library revolutionized logging with `log/slog`.

---

## 1. Structured Logging with `log/slog`

Unstructured strings (`log.Printf("user %s logged in", id)`) are painful to search in cloud log aggregators (Datadog, Loki, CloudWatch). `log/slog` provides structured, high-performance logging:

```go
// Create a JSON handler
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)

// Log with structured key-value attributes
slog.InfoContext(ctx, "order processed",
    slog.String("order_id", "ord-123"),
    slog.Float64("amount", 99.95),
    slog.Duration("latency", 42*time.Millisecond),
)
```

### Masking Sensitive Data (`slog.LogValuer`)
Types can implement `LogValue() slog.Value` to control their log presentation and mask sensitive credentials:
```go
func (u User) LogValue() slog.Value {
    return slog.GroupValue(
        slog.String("id", u.ID),
        slog.String("password", "[REDACTED]"),
    )
}
```

---

## 2. Prometheus Metrics & Exposition Format

Prometheus uses a line-oriented text format served at `/metrics`:
```text
# HELP http_requests_total Total number of HTTP requests
# TYPE http_requests_total counter
http_requests_total{method="GET",status="200"} 1042

# HELP active_connections Current number of active client connections
# TYPE active_connections gauge
active_connections 27
```

---

## 3. Distributed Tracing (W3C Trace Context)

When requests hop across multiple microservices, W3C `traceparent` headers link all logs and spans to a single end-to-end trace:
```text
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
              │  └─────────────┬────────────────┘ └───────┬──────┘ └─┬┘
           version          trace_id                   parent_id   flags
```

---

## 4. Directory Layout

```
19-observability/
├── README.md
├── examples/
│   ├── slog_json/main.go          # JSON structured logging with PII masking
│   └── metrics_collector/main.go  # Prometheus text format metric collector
├── exercises/
│   ├── 01-slog-sanitizer          # Custom slog handler and LogValuer sanitizer
│   ├── 02-prometheus-exporter     # Thread-safe metric registry and /metrics exporter
│   └── 03-w3c-trace-propagator    # W3C traceparent parser and context injector
├── solutions/
│   ├── 01-slog-sanitizer
│   ├── 02-prometheus-exporter
│   └── 03-w3c-trace-propagator
└── challenges/
    └── observability-middleware   # Full telemetry middleware (logs, metrics, traces)
```
