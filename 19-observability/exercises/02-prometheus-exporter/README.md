# Exercise 02: Prometheus Metric Registry & Exporter

## Objective
Implement an in-memory Prometheus metric collector and exporter that supports labeled counters and gauges formatted to standard Prometheus exposition specifications.

## Requirements
1. **Metric Types**:
   - `Counter`: Monotonically increasing metric (`Inc()`, `Add(val)`).
   - `Gauge`: Fluctuating metric (`Set(val)`, `Inc()`, `Dec()`).
2. **Labels**:
   - Support arbitrary dimensional labels (e.g. `method="GET"`, `status="200"`).
3. **`Export() string`**:
   - Renders output in the official Prometheus 0.0.4 text format:
     ```text
     # HELP metric_name Help text
     # TYPE metric_name counter
     metric_name{label="val"} 42
     ```
   - Separated by blank lines between metric families.

## Starter & Tests
- Starter: `starter.go`
- Run starter tests: `go test -v ./19-observability/exercises/02-prometheus-exporter`
