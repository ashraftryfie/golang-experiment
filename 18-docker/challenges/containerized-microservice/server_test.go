package microservice

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestServer_LivenessAndReadiness(t *testing.T) {
	srv := NewServer()
	ts := httptest.NewServer(srv)
	defer ts.Close()

	// 1. Liveness check -> 200 OK
	respLive, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("liveness failed: %v", err)
	}
	defer respLive.Body.Close()
	if respLive.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for liveness, got %d", respLive.StatusCode)
	}

	// 2. Initial readiness check -> 200 OK
	respReady, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("readiness failed: %v", err)
	}
	defer respReady.Body.Close()
	if respReady.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK for initial readiness, got %d", respReady.StatusCode)
	}

	// 3. Degrade readiness -> Expect 503 Service Unavailable
	srv.SetReady(false)
	respDegraded, err := http.Get(ts.URL + "/readyz")
	if err != nil {
		t.Fatalf("degraded readiness failed: %v", err)
	}
	defer respDegraded.Body.Close()
	if respDegraded.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected 503 for degraded readiness, got %d", respDegraded.StatusCode)
	}

	// 4. In-binary self-probe works
	srv.SetReady(true)
	if err := RunHealthProbe(ts.URL+"/healthz", 100*time.Millisecond); err != nil {
		t.Fatalf("self-health probe failed: %v", err)
	}
}

func TestServer_StatusAPI(t *testing.T) {
	srv := NewServer()
	ts := httptest.NewServer(srv)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/status")
	if err != nil {
		t.Fatalf("status API failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}
}

func TestDockerfile_SecurityStandards(t *testing.T) {
	content, err := os.ReadFile("Dockerfile")
	if err != nil {
		t.Fatalf("failed reading challenge Dockerfile: %v", err)
	}

	strContent := string(content)

	// Check multi-stage
	fromCount := strings.Count(strContent, "FROM ")
	if fromCount < 2 {
		t.Errorf("expected at least 2 FROM instructions in Dockerfile, got %d", fromCount)
	}

	// Check static compilation
	if !strings.Contains(strContent, "CGO_ENABLED=0") {
		t.Errorf("Dockerfile missing CGO_ENABLED=0")
	}

	// Check nonroot user
	if !strings.Contains(strContent, "USER nonroot") {
		t.Errorf("Dockerfile missing non-root USER instruction")
	}
}
