# Project 4: `event-dispatcher` (High-Throughput Event Pipeline)

A robust, concurrent event dispatching and processing engine implementing fan-out and fan-in concurrency patterns, bounded buffers, backpressure protection, and atomic telemetry.

## Features
- **Fan-Out / Fan-In Topology**: Work is distributed among worker routines and aggregated into unified completion streams.
- **Backpressure Handling**: Non-blocking dispatch options with overflow tracking.
- **Race Protection**: Clean concurrency tested with Go's race detector.

## Usage
```powershell
go test -v ./projects/04-concurrent-worker/...
```
