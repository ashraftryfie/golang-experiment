package telemetry

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

type traceKey struct{}

var traceCtxKey = traceKey{}

func GetTraceID(ctx context.Context) string {
	val := ctx.Value(traceCtxKey)
	if val == nil {
		return ""
	}
	return val.(string)
}

// responseRecorder captures the status code written by downstream handlers
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// In-memory Prometheus metrics
type MetricsService struct {
	requestsTotal  map[string]int64
	activeRequests int64
	mu             sync.RWMutex
}

func NewMetricsService() *MetricsService {
	return &MetricsService{
		requestsTotal: make(map[string]int64),
	}
}

func (m *MetricsService) recordRequest(method string, status int) {
	key := fmt.Sprintf(`{method="%s",status="%d"}`, method, status)
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsTotal[key]++
}

func (m *MetricsService) incActive() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequests++
}

func (m *MetricsService) decActive() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.activeRequests--
}

func (m *MetricsService) Export() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var sb strings.Builder

	sb.WriteString("# HELP http_requests_total Total number of HTTP requests processed\n")
	sb.WriteString("# TYPE http_requests_total counter\n")

	keys := make([]string, 0, len(m.requestsTotal))
	for k := range m.requestsTotal {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("http_requests_total%s %d\n", k, m.requestsTotal[k]))
	}

	sb.WriteString("\n# HELP http_active_requests Number of requests currently in flight\n")
	sb.WriteString("# TYPE http_active_requests gauge\n")
	sb.WriteString(fmt.Sprintf("http_active_requests %d\n", m.activeRequests))

	return sb.String()
}

func generateTraceID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func generateParentID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type TelemetryApp struct {
	Metrics *MetricsService
	Logger  *slog.Logger
	Router  http.Handler
}

func NewTelemetryApp(logOut io.Writer) *TelemetryApp {
	metrics := NewMetricsService()
	logger := slog.New(slog.NewJSONHandler(logOut, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	mux := http.NewServeMux()

	// Metrics endpoint
	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(metrics.Export()))
	})

	// Sample application endpoints
	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		traceID := GetTraceID(r.Context())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(fmt.Sprintf(`{"message":"hello","trace_id":"%s"}`, traceID)))
	})

	mux.HandleFunc("GET /api/error", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":"internal failure"}`, http.StatusInternalServerError)
	})

	app := &TelemetryApp{
		Metrics: metrics,
		Logger:  logger,
	}

	// Wrap entire mux in telemetry middleware
	app.Router = app.middleware(mux)
	return app
}

func (app *TelemetryApp) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// 1. Trace Extraction / Generation
		rawTraceparent := r.Header.Get("traceparent")
		var traceID string
		var parentID string

		parts := strings.Split(rawTraceparent, "-")
		if len(parts) == 4 && len(parts[1]) == 32 {
			traceID = parts[1]
			parentID = parts[2]
		} else {
			traceID = generateTraceID()
			parentID = generateParentID()
		}

		responseTraceparent := fmt.Sprintf("00-%s-%s-01", traceID, parentID)
		w.Header().Set("traceparent", responseTraceparent)

		ctx := context.WithValue(r.Context(), traceCtxKey, traceID)
		r = r.WithContext(ctx)

		// 2. Metrics (In-Flight)
		app.Metrics.incActive()
		defer app.Metrics.decActive()

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // default if WriteHeader not explicitly called
		}

		// 3. Process Request
		next.ServeHTTP(rec, r)

		elapsed := time.Since(start)

		// 4. Metrics (Completed)
		app.Metrics.recordRequest(r.Method, rec.statusCode)

		// 5. Structured Logging
		app.Logger.InfoContext(ctx, "http_request",
			slog.String("trace_id", traceID),
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", rec.statusCode),
			slog.Int64("duration_ms", elapsed.Milliseconds()),
			slog.String("remote_ip", r.RemoteAddr),
		)
	})
}

func (app *TelemetryApp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	app.Router.ServeHTTP(w, r)
}
