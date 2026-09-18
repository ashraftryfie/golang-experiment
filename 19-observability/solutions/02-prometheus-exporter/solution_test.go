package promexporter

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

func TestPrometheusExporter_CounterAndGauge(t *testing.T) {
	reg := NewRegistry()

	counter, err := reg.RegisterCounter("http_requests_total", "Total requests", []string{"method", "code"})
	if err != nil {
		t.Fatalf("failed registering counter: %v", err)
	}

	gauge, err := reg.RegisterGauge("active_connections", "Active client connections", nil)
	if err != nil {
		t.Fatalf("failed registering gauge: %v", err)
	}

	counter.Inc("GET", "200")
	counter.Inc("GET", "200")
	counter.Inc("POST", "201")

	gauge.Set(10)
	gauge.Inc()
	gauge.Dec()

	out := reg.Export()

	// Verify header comments
	if !strings.Contains(out, "# HELP http_requests_total Total requests\n# TYPE http_requests_total counter\n") {
		t.Errorf("missing or incorrect counter header: %s", out)
	}
	if !strings.Contains(out, "# HELP active_connections Active client connections\n# TYPE active_connections gauge\n") {
		t.Errorf("missing or incorrect gauge header: %s", out)
	}

	// Verify values
	if !strings.Contains(out, `http_requests_total{method="GET",code="200"} 2`) {
		t.Errorf("expected GET 200 count 2, got: %s", out)
	}
	if !strings.Contains(out, `http_requests_total{method="POST",code="201"} 1`) {
		t.Errorf("expected POST 201 count 1, got: %s", out)
	}
	if !strings.Contains(out, `active_connections 10`) {
		t.Errorf("expected active_connections 10, got: %s", out)
	}
}

func TestPrometheusExporter_DuplicateRegistration(t *testing.T) {
	reg := NewRegistry()
	_, _ = reg.RegisterCounter("my_metric", "help", nil)
	_, err := reg.RegisterCounter("my_metric", "help", nil)
	if !errors.Is(err, ErrDuplicateMetric) {
		t.Fatalf("expected ErrDuplicateMetric, got %v", err)
	}
}

func TestPrometheusExporter_ConcurrentSafety(t *testing.T) {
	reg := NewRegistry()
	c, _ := reg.RegisterCounter("concurrency_test", "help", []string{"worker"})

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 50; j++ {
				c.Inc("worker-1")
			}
		}()
	}
	wg.Wait()

	out := reg.Export()
	if !strings.Contains(out, `concurrency_test{worker="worker-1"} 1000`) {
		t.Fatalf("expected 1000 count from concurrent increments, got: %s", out)
	}
}
