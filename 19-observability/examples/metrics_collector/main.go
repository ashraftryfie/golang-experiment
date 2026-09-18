package main

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
)

type Metrics struct {
	httpRequestsTotal int64
	activeWorkers     int64
	mu                sync.RWMutex
}

func main() {
	m := &Metrics{}

	// Simulate activity
	atomic.AddInt64(&m.httpRequestsTotal, 142)
	atomic.AddInt64(&m.activeWorkers, 5)

	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")

		fmt.Fprintln(w, "# HELP http_requests_total Total number of HTTP requests processed")
		fmt.Fprintln(w, "# TYPE http_requests_total counter")
		fmt.Fprintf(w, "http_requests_total %d\n\n", atomic.LoadInt64(&m.httpRequestsTotal))

		fmt.Fprintln(w, "# HELP active_workers Current number of active workers")
		fmt.Fprintln(w, "# TYPE active_workers gauge")
		fmt.Fprintf(w, "active_workers %d\n", atomic.LoadInt64(&m.activeWorkers))
	})

	fmt.Println("[Metrics] Serving Prometheus metrics on :9090/metrics")
	_ = http.ListenAndServe(":9090", nil)
}
