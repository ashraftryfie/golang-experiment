package telemetry

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTelemetry_EndToEnd(t *testing.T) {
	logBuf := new(bytes.Buffer)
	app := NewTelemetryApp(logBuf)
	ts := httptest.NewServer(app)
	defer ts.Close()

	client := &http.Client{}

	// 1. Test Hello without incoming traceparent -> should generate one
	resp1, err := client.Get(ts.URL + "/api/hello")
	if err != nil {
		t.Fatalf("hello request failed: %v", err)
	}
	defer resp1.Body.Close()

	tp1 := resp1.Header.Get("traceparent")
	if tp1 == "" || !strings.HasPrefix(tp1, "00-") {
		t.Fatalf("expected valid traceparent header in response, got: %s", tp1)
	}

	// Read body and verify trace_id returned matches response header
	var body1 map[string]string
	_ = json.NewDecoder(resp1.Body).Decode(&body1)
	parts1 := strings.Split(tp1, "-")
	if body1["trace_id"] != parts1[1] {
		t.Errorf("trace ID mismatch: body=%s, header=%s", body1["trace_id"], parts1[1])
	}

	// 2. Test Hello WITH existing incoming traceparent -> should preserve trace_id
	customTraceID := "abcdef0123456789abcdef0123456789"
	incomingTP := "00-" + customTraceID + "-1234567890abcdef-01"
	req2, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/hello", nil)
	req2.Header.Set("traceparent", incomingTP)

	resp2, err := client.Do(req2)
	if err != nil {
		t.Fatalf("req2 failed: %v", err)
	}
	defer resp2.Body.Close()

	var body2 map[string]string
	_ = json.NewDecoder(resp2.Body).Decode(&body2)
	if body2["trace_id"] != customTraceID {
		t.Errorf("expected preserved trace ID %s, got %s", customTraceID, body2["trace_id"])
	}

	// 3. Test Error endpoint -> 500
	respErr, _ := client.Get(ts.URL + "/api/error")
	respErr.Body.Close()
	if respErr.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 from /api/error, got %d", respErr.StatusCode)
	}

	// 4. Test /metrics endpoint -> Prometheus output
	respMetrics, _ := client.Get(ts.URL + "/metrics")
	metricsBytes, _ := io.ReadAll(respMetrics.Body)
	respMetrics.Body.Close()

	metricsOut := string(metricsBytes)
	if !strings.Contains(metricsOut, `http_requests_total{method="GET",status="200"} 2`) {
		t.Errorf("metrics missing 2x GET 200: %s", metricsOut)
	}
	if !strings.Contains(metricsOut, `http_requests_total{method="GET",status="500"} 1`) {
		t.Errorf("metrics missing 1x GET 500: %s", metricsOut)
	}

	// 5. Verify Structured JSON Log Output
	logStr := logBuf.String()
	if !strings.Contains(logStr, customTraceID) {
		t.Errorf("expected structured log to contain trace ID %s: %s", customTraceID, logStr)
	}
	if !strings.Contains(logStr, `"status":500`) {
		t.Errorf("expected structured log to contain status 500: %s", logStr)
	}
}
