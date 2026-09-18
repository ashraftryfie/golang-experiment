package resilientgateway

import (
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestGateway_RateLimiting(t *testing.T) {
	// Dummy upstream returning 200 OK
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("upstream response"))
	}))
	defer upstream.Close()

	// Rate limiter: 1 req/sec, burst: 3
	gw := NewGateway(GatewayConfig{
		RateLimitPerSec: 1.0,
		BurstCapacity:   3,
		UpstreamURL:     upstream.URL,
	})

	ts := httptest.NewServer(gw)
	defer ts.Close()

	client := &http.Client{}

	// First 3 requests with same client ID should succeed
	for i := 1; i <= 3; i++ {
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/test", nil)
		req.Header.Set("X-Client-ID", "test-client")
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("request %d failed: %v", i, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected 200 OK on burst request %d, got %d", i, resp.StatusCode)
		}
	}

	// 4th request must be rate limited
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/test", nil)
	req.Header.Set("X-Client-ID", "test-client")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("4th request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Fatalf("expected 429 Too Many Requests, got %d", resp.StatusCode)
	}
	if resp.Header.Get("Retry-After") == "" {
		t.Errorf("expected Retry-After header on 429 response")
	}

	// Different client ID should succeed
	reqOther, _ := http.NewRequest(http.MethodGet, ts.URL+"/api/test", nil)
	reqOther.Header.Set("X-Client-ID", "different-client")
	respOther, err := client.Do(reqOther)
	if err != nil {
		t.Fatalf("different client request failed: %v", err)
	}
	defer respOther.Body.Close()
	if respOther.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for separate client, got %d", respOther.StatusCode)
	}
}

func TestGateway_CircuitBreakerTripsAndRecovers(t *testing.T) {
	var upstreamFails int32
	atomic.StoreInt32(&upstreamFails, 1)

	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if atomic.LoadInt32(&upstreamFails) == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("healthy"))
	}))
	defer upstream.Close()

	// Breaker trips after 2 failures, timeout 50ms
	gw := NewGateway(GatewayConfig{
		RateLimitPerSec: 100.0,
		BurstCapacity:   50,
		BreakerFailures: 2,
		BreakerTimeout:  50 * time.Millisecond,
		UpstreamURL:     upstream.URL,
	})

	ts := httptest.NewServer(gw)
	defer ts.Close()

	client := &http.Client{}

	// Failure 1
	resp1, _ := client.Get(ts.URL + "/data")
	resp1.Body.Close()
	if resp1.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 on failure 1, got %d", resp1.StatusCode)
	}

	// Failure 2 -> Breaker should trip to Open
	resp2, _ := client.Get(ts.URL + "/data")
	resp2.Body.Close()
	if resp2.StatusCode != http.StatusInternalServerError {
		t.Fatalf("expected 500 on failure 2, got %d", resp2.StatusCode)
	}

	// Immediate next request should fail fast with 503 Service Unavailable
	resp3, _ := client.Get(ts.URL + "/data")
	body3, _ := io.ReadAll(resp3.Body)
	resp3.Body.Close()
	if resp3.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("expected 503 Service Unavailable when open, got %d (body: %s)", resp3.StatusCode, string(body3))
	}

	// Upstream recovers
	atomic.StoreInt32(&upstreamFails, 0)

	// Wait for breaker timeout to transition to Half-Open
	time.Sleep(60 * time.Millisecond)

	// Next request acts as canary and succeeds -> Breaker heals to Closed
	respCanary, _ := client.Get(ts.URL + "/data")
	respCanary.Body.Close()
	if respCanary.StatusCode != http.StatusOK {
		t.Fatalf("expected canary 200 OK, got %d", respCanary.StatusCode)
	}

	// Subsequent request succeeds
	respAfter, _ := client.Get(ts.URL + "/data")
	respAfter.Body.Close()
	if respAfter.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK after healing, got %d", respAfter.StatusCode)
	}
}

func TestGateway_Lifecycle(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer upstream.Close()

	gw := NewGateway(GatewayConfig{
		RateLimitPerSec: 10.0,
		BurstCapacity:   5,
		UpstreamURL:     upstream.URL,
	})

	if err := gw.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Start failed: %v", err)
	}

	addr := gw.Addr()
	resp, err := http.Get("http://" + addr + "/ping")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200 OK, got %d", resp.StatusCode)
	}

	if err := gw.Shutdown(100 * time.Millisecond); err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
}
