package resilientgateway

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	ErrCircuitOpen      = errors.New("upstream circuit breaker open")
	ErrServerNotRunning = errors.New("gateway server not running")
)

// TokenBucket per-client limiter
type tokenBucket struct {
	tokens     float64
	capacity   float64
	refillRate float64
	lastRefill time.Time
}

func (tb *tokenBucket) allow() bool {
	now := time.Now()
	elapsed := now.Sub(tb.lastRefill).Seconds()
	tb.lastRefill = now

	tb.tokens += elapsed * tb.refillRate
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity
	}

	if tb.tokens >= 1.0 {
		tb.tokens -= 1.0
		return true
	}
	return false
}

// Circuit Breaker
type BreakerState int

const (
	StateClosed BreakerState = iota
	StateOpen
	StateHalfOpen
)

type circuitBreaker struct {
	maxFailures      int
	timeout          time.Duration
	successThreshold int

	state        BreakerState
	failures     int
	successes    int
	lastOpenedAt time.Time
	mu           sync.Mutex
}

func newCircuitBreaker(maxFailures int, timeout time.Duration, successThreshold int) *circuitBreaker {
	return &circuitBreaker{
		maxFailures:      maxFailures,
		timeout:          timeout,
		successThreshold: successThreshold,
		state:            StateClosed,
	}
}

func (cb *circuitBreaker) checkTransition() {
	if cb.state == StateOpen && time.Since(cb.lastOpenedAt) >= cb.timeout {
		cb.state = StateHalfOpen
		cb.successes = 0
	}
}

func (cb *circuitBreaker) canExecute() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.checkTransition()
	return cb.state != StateOpen
}

func (cb *circuitBreaker) recordResult(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.checkTransition()

	if err != nil {
		switch cb.state {
		case StateClosed:
			cb.failures++
			if cb.failures >= cb.maxFailures {
				cb.state = StateOpen
				cb.lastOpenedAt = time.Now()
			}
		case StateHalfOpen:
			cb.state = StateOpen
			cb.lastOpenedAt = time.Now()
		}
	} else {
		switch cb.state {
		case StateClosed:
			cb.failures = 0
		case StateHalfOpen:
			cb.successes++
			if cb.successes >= cb.successThreshold {
				cb.state = StateClosed
				cb.failures = 0
				cb.successes = 0
			}
		}
	}
}

type GatewayConfig struct {
	RateLimitPerSec float64
	BurstCapacity   int
	BreakerFailures int
	BreakerTimeout  time.Duration
	UpstreamURL     string
}

type Gateway struct {
	config  GatewayConfig
	breaker *circuitBreaker
	client  *http.Client

	buckets map[string]*tokenBucket
	buckMu  sync.Mutex

	server   *http.Server
	listener net.Listener
	servMu   sync.Mutex
	running  bool
}

func NewGateway(cfg GatewayConfig) *Gateway {
	if cfg.RateLimitPerSec <= 0 {
		cfg.RateLimitPerSec = 10.0
	}
	if cfg.BurstCapacity <= 0 {
		cfg.BurstCapacity = 5
	}
	if cfg.BreakerFailures <= 0 {
		cfg.BreakerFailures = 3
	}
	if cfg.BreakerTimeout <= 0 {
		cfg.BreakerTimeout = 1 * time.Second
	}

	return &Gateway{
		config:  cfg,
		breaker: newCircuitBreaker(cfg.BreakerFailures, cfg.BreakerTimeout, 1),
		client:  &http.Client{Timeout: 5 * time.Second},
		buckets: make(map[string]*tokenBucket),
	}
}

func (g *Gateway) getClientID(r *http.Request) string {
	cid := r.Header.Get("X-Client-ID")
	if cid != "" {
		return cid
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && host != "" {
		return host
	}
	return r.RemoteAddr
}

func (g *Gateway) allowRequest(clientID string) bool {
	g.buckMu.Lock()
	defer g.buckMu.Unlock()

	bucket, exists := g.buckets[clientID]
	if !exists {
		bucket = &tokenBucket{
			tokens:     float64(g.config.BurstCapacity),
			capacity:   float64(g.config.BurstCapacity),
			refillRate: g.config.RateLimitPerSec,
			lastRefill: time.Now(),
		}
		g.buckets[clientID] = bucket
	}

	return bucket.allow()
}

func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 1. Rate Limiting
	clientID := g.getClientID(r)
	if !g.allowRequest(clientID) {
		w.Header().Set("Retry-After", "1")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"error":"rate limit exceeded"}`))
		return
	}

	// 2. Circuit Breaker Fast-Fail
	if !g.breaker.canExecute() {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte(`{"error":"upstream circuit breaker open"}`))
		return
	}

	// 3. Proxy to Upstream
	targetURL := g.config.UpstreamURL + r.URL.Path
	if r.URL.RawQuery != "" {
		targetURL += "?" + r.URL.RawQuery
	}

	upstreamReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL, r.Body)
	if err != nil {
		http.Error(w, `{"error":"failed to create proxy request"}`, http.StatusInternalServerError)
		return
	}
	for k, vv := range r.Header {
		for _, v := range vv {
			upstreamReq.Header.Add(k, v)
		}
	}

	resp, err := g.client.Do(upstreamReq)
	if err != nil {
		g.breaker.recordResult(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte(`{"error":"upstream unavailable"}`))
		return
	}
	defer resp.Body.Close()

	// If upstream returned 5xx server error, record as failure for circuit breaker
	if resp.StatusCode >= 500 {
		g.breaker.recordResult(fmt.Errorf("upstream error %d", resp.StatusCode))
	} else {
		g.breaker.recordResult(nil)
	}

	// Copy response headers and body to client
	for k, vv := range resp.Header {
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (g *Gateway) Start(addr string) error {
	g.servMu.Lock()
	defer g.servMu.Unlock()

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	g.listener = l
	g.server = &http.Server{Handler: g}
	g.running = true

	go func() {
		_ = g.server.Serve(g.listener)
	}()

	return nil
}

func (g *Gateway) Addr() string {
	g.servMu.Lock()
	defer g.servMu.Unlock()
	if g.listener == nil {
		return ""
	}
	return g.listener.Addr().String()
}

func (g *Gateway) Shutdown(timeout time.Duration) error {
	g.servMu.Lock()
	if !g.running {
		g.servMu.Unlock()
		return ErrServerNotRunning
	}
	g.running = false
	g.servMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return g.server.Shutdown(ctx)
}
