package microservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type Server struct {
	ready     int32 // 1 = ready, 0 = degraded
	startTime time.Time
	mux       *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		startTime: time.Now(),
		mux:       http.NewServeMux(),
	}
	atomic.StoreInt32(&s.ready, 1)

	s.mux.HandleFunc("/healthz", s.handleLiveness)
	s.mux.HandleFunc("/readyz", s.handleReadiness)
	s.mux.HandleFunc("/api/status", s.handleStatus)

	return s
}

func (s *Server) SetReady(ready bool) {
	if ready {
		atomic.StoreInt32(&s.ready, 1)
	} else {
		atomic.StoreInt32(&s.ready, 0)
	}
}

func (s *Server) handleLiveness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"alive"}`))
}

func (s *Server) handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if atomic.LoadInt32(&s.ready) == 1 {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ready"}`))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"status":"degraded","reason":"database unavailable"}`))
	}
}

type statusResponse struct {
	UptimeSeconds float64 `json:"uptime_seconds"`
	Ready         bool    `json:"ready"`
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(statusResponse{
		UptimeSeconds: time.Since(s.startTime).Seconds(),
		Ready:         atomic.LoadInt32(&s.ready) == 1,
	})
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// RunHealthProbe provides an in-binary health check for distroless execution.
func RunHealthProbe(url string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("probe failed with status: %d", resp.StatusCode)
	}
	return nil
}
