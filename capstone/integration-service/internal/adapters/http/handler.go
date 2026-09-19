package httpadapter

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/service"
)

type ServerHandler struct {
	svc          *service.IntegrationService
	mux          *http.ServeMux
	jobsEnqueued uint64
}

func NewServerHandler(svc *service.IntegrationService) *ServerHandler {
	h := &ServerHandler{
		svc: svc,
		mux: http.NewServeMux(),
	}
	h.registerRoutes()
	return h
}

func (h *ServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (h *ServerHandler) registerRoutes() {
	h.mux.HandleFunc("POST /api/v1/jobs", h.handleCreateJob)
	h.mux.HandleFunc("GET /api/v1/jobs/{id}", h.handleGetJob)
	h.mux.HandleFunc("GET /api/v1/jobs", h.handleListJobs)
	h.mux.HandleFunc("GET /healthz", h.handleHealthz)
	h.mux.HandleFunc("GET /metrics", h.handleMetrics)
}

type createJobRequest struct {
	Type       string            `json:"type"`
	Payload    map[string]string `json:"payload"`
	MaxRetries int               `json:"max_retries"`
}

func (h *ServerHandler) handleCreateJob(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	job, err := h.svc.SubmitJob(r.Context(), req.Type, req.Payload, req.MaxRetries)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	atomic.AddUint64(&h.jobsEnqueued, 1)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(job)
}

func (h *ServerHandler) handleGetJob(w http.ResponseWriter, r *http.Request) {
	jobID := r.PathValue("id")
	if jobID == "" {
		http.Error(w, "missing job id", http.StatusBadRequest)
		return
	}

	job, err := h.svc.GetJob(r.Context(), jobID)
	if err != nil {
		if errors.Is(err, domain.ErrJobNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(job)
}

func (h *ServerHandler) handleListJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.svc.ListJobs(r.Context(), 50)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(jobs)
}

func (h *ServerHandler) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"UP"}`))
}

func (h *ServerHandler) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	total := atomic.LoadUint64(&h.jobsEnqueued)
	fmt.Fprintf(w, "# HELP integration_jobs_submitted_total Total integration jobs submitted\n")
	fmt.Fprintf(w, "# TYPE integration_jobs_submitted_total counter\n")
	fmt.Fprintf(w, "integration_jobs_submitted_total %d\n", total)
}
