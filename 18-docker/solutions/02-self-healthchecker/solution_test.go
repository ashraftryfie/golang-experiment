package selfhealth

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCLIArgParsing(t *testing.T) {
	tests := []struct {
		args        []string
		wantProbe   bool
		wantURL     string
	}{
		{args: []string{"server", "-healthcheck"}, wantProbe: true, wantURL: DefaultHealthURL},
		{args: []string{"server", "--healthcheck"}, wantProbe: true, wantURL: DefaultHealthURL},
		{args: []string{"server", "-healthcheck=http://localhost:9090/live"}, wantProbe: true, wantURL: "http://localhost:9090/live"},
		{args: []string{"server", "--port=8080"}, wantProbe: false, wantURL: ""},
	}

	for _, tc := range tests {
		probe, url := ParseCLIArgs(tc.args)
		if probe != tc.wantProbe || url != tc.wantURL {
			t.Errorf("args %v: got (%v, %s), want (%v, %s)", tc.args, probe, url, tc.wantProbe, tc.wantURL)
		}
	}
}

func TestHealthProbe_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	if err := RunHealthProbe(ts.URL, 200*time.Millisecond); err != nil {
		t.Fatalf("expected probe success, got %v", err)
	}
}

func TestHealthProbe_UnhealthyStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	err := RunHealthProbe(ts.URL, 200*time.Millisecond)
	if !errors.Is(err, ErrUnhealthyStatus) {
		t.Fatalf("expected ErrUnhealthyStatus, got %v", err)
	}
}

func TestHealthProbe_Timeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	err := RunHealthProbe(ts.URL, 20*time.Millisecond)
	if err == nil {
		t.Fatalf("expected timeout error on slow health endpoint, got nil")
	}
}
