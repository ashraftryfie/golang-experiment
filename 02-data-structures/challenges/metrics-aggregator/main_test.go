package main

import (
	"math"
	"reflect"
	"testing"
)

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-4
}

func TestAggregator(t *testing.T) {
	agg := NewAggregator()

	records := []LogRecord{
		{Service: "auth", StatusCode: 200, LatencyMs: 50.0},
		{Service: "auth", StatusCode: 503, LatencyMs: 150.0},
		{Service: "auth", StatusCode: 200, LatencyMs: 100.0},
		{Service: "payment", StatusCode: 200, LatencyMs: 80.0},
		{Service: "payment", StatusCode: 500, LatencyMs: 120.0},
		{Service: "payment", StatusCode: 502, LatencyMs: 160.0},
	}

	for _, r := range records {
		agg.Ingest(r)
	}

	// Verify Auth metrics
	authMetrics, ok := agg.GetMetrics("auth")
	if !ok {
		t.Fatalf("expected auth metrics to exist")
	}
	if authMetrics.TotalRequests != 3 {
		t.Errorf("got total requests %d, want 3", authMetrics.TotalRequests)
	}
	if authMetrics.ErrorCount != 1 {
		t.Errorf("got error count %d, want 1", authMetrics.ErrorCount)
	}
	if !almostEqual(authMetrics.AvgLatencyMs, 100.0) {
		t.Errorf("got avg latency %v, want 100.0", authMetrics.AvgLatencyMs)
	}
	if !almostEqual(authMetrics.MaxLatencyMs, 150.0) {
		t.Errorf("got max latency %v, want 150.0", authMetrics.MaxLatencyMs)
	}

	// Verify Services
	services := agg.Services()
	wantServices := []string{"auth", "payment"}
	if !reflect.DeepEqual(services, wantServices) {
		t.Errorf("Services() = %v, want %v", services, wantServices)
	}

	// Verify Top Failing
	top := agg.TopFailingServices(1)
	if len(top) != 1 || top[0].Service != "payment" || top[0].ErrorCount != 2 {
		t.Errorf("TopFailingServices() = %v, want payment with 2 errors", top)
	}
}
