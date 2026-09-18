package promexporter

import (
	"errors"
	"strings"
	"testing"
)

func TestPrometheusExporter_Starter(t *testing.T) {
	reg := NewRegistry()
	if reg == nil {
		t.Skip("Skipping unimplemented exercise: NewRegistry")
	}

	counter, err := reg.RegisterCounter("http_requests_total", "Total requests", []string{"method"})
	if errors.Is(err, ErrNotImplemented) {
		t.Skip("Skipping unimplemented exercise: RegisterCounter")
	}
	if err != nil {
		t.Fatalf("unexpected error registering counter: %v", err)
	}

	counter.Inc("GET")
	output := reg.Export()

	if !strings.Contains(output, "http_requests_total{method=\"GET\"} 1") {
		t.Errorf("export missing expected metric line: %s", output)
	}
}
