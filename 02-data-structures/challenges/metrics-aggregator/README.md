# Challenge: Log Metrics Aggregator

## 🎯 Goal
Build a high-performance in-memory log stream metrics aggregator that processes structured request records and computes aggregated statistics per service using slices and maps.

---

## 📋 Requirements

### Input Data
The aggregator ingests `LogRecord` instances:
```go
type LogRecord struct {
    Service    string
    StatusCode int
    LatencyMs  float64
}
```

### Metrics Output
Aggregates records into `ServiceMetrics`:
```go
type ServiceMetrics struct {
    TotalRequests int
    ErrorCount    int     // Status codes >= 500
    AvgLatencyMs  float64 // Mean latency
    MaxLatencyMs  float64 // Peak latency
}
```

The aggregator must provide:
1. `Ingest(record LogRecord)`: Ingests a new record in `O(1)` amortized time.
2. `GetMetrics(service string) (ServiceMetrics, bool)`: Retrieves metrics for a specific service.
3. `Services() []string`: Returns a sorted list of all known services.
4. `TopFailingServices(limit int) []ServiceSummary`: Returns top services ranked by highest error count.

---

## 🧪 Testing
Run unit tests:
```powershell
go test -v ./02-data-structures/challenges/metrics-aggregator
```
