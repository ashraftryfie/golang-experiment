package test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/external"
	httpadapter "github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/http"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/memory"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/service"
)

func TestIntegrationService_EndToEnd(t *testing.T) {
	storage := memory.NewInMemoryJobStorage()
	queue := memory.NewChannelJobQueue(50)
	client := external.NewResilientExternalClient()
	client.Latency = 1 * time.Millisecond

	svc := service.NewIntegrationService(storage, queue)

	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	pool := service.NewWorkerPool(3, storage, queue, client, 2*time.Second)
	pool.Start(workerCtx)
	defer pool.Stop()

	handler := httpadapter.NewServerHandler(svc)
	server := httptest.NewServer(handler)
	defer server.Close()

	httpClient := server.Client()

	// 1. Submit Job via HTTP POST
	reqBody := `{"type":"sync_inventory","payload":{"sku":"ITEM-123","qty":"10"},"max_retries":2}`
	resp, err := httpClient.Post(server.URL+"/api/v1/jobs", "application/json", bytes.NewBufferString(reqBody))
	if err != nil {
		t.Fatalf("POST /api/v1/jobs failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 Created, got %d", resp.StatusCode)
	}

	var submitted domain.Job
	if err := json.NewDecoder(resp.Body).Decode(&submitted); err != nil {
		t.Fatalf("failed decoding submitted job: %v", err)
	}
	resp.Body.Close()

	if submitted.ID == "" || submitted.Status != domain.StatusPending {
		t.Fatalf("unexpected submitted job state: %+v", submitted)
	}

	// 2. Poll until worker completes the job (up to 3 seconds)
	var finalJob *domain.Job
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		getResp, err := httpClient.Get(server.URL + "/api/v1/jobs/" + submitted.ID)
		if err != nil {
			t.Fatalf("GET /api/v1/jobs failed: %v", err)
		}

		var j domain.Job
		_ = json.NewDecoder(getResp.Body).Decode(&j)
		getResp.Body.Close()

		if j.Status == domain.StatusCompleted {
			finalJob = &j
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if finalJob == nil {
		t.Fatalf("job did not complete within deadline")
	}
	if finalJob.Attempts < 1 {
		t.Fatalf("expected at least 1 attempt, got %d", finalJob.Attempts)
	}

	// 3. Verify downstream external client received the job
	synced := client.GetSyncedJobs()
	if len(synced) != 1 || synced[0] != submitted.ID {
		t.Fatalf("external client did not record synced job: %v", synced)
	}

	// 4. Verify Health and Metrics
	healthResp, err := httpClient.Get(server.URL + "/healthz")
	if err != nil || healthResp.StatusCode != http.StatusOK {
		t.Fatalf("healthz check failed: %v, status: %d", err, healthResp.StatusCode)
	}
	healthResp.Body.Close()

	metricsResp, err := httpClient.Get(server.URL + "/metrics")
	if err != nil || metricsResp.StatusCode != http.StatusOK {
		t.Fatalf("metrics check failed: %v, status: %d", err, metricsResp.StatusCode)
	}
	metricsBody, _ := io.ReadAll(metricsResp.Body)
	metricsResp.Body.Close()

	if !strings.Contains(string(metricsBody), "integration_jobs_submitted_total 1") {
		t.Fatalf("expected metric integration_jobs_submitted_total 1, got: %s", string(metricsBody))
	}
}

func TestWorkerPool_RetryAndExhaustion(t *testing.T) {
	storage := memory.NewInMemoryJobStorage()
	queue := memory.NewChannelJobQueue(50)
	client := external.NewResilientExternalClient()
	client.Latency = 1 * time.Millisecond
	client.FailForJobID = "failing-job-001" // Always fail this job

	svc := service.NewIntegrationService(storage, queue)
	svc.SetIDGenerator(func() string { return "failing-job-001" })

	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	defer cancelWorkers()

	pool := service.NewWorkerPool(2, storage, queue, client, 2*time.Second)
	pool.Start(workerCtx)
	defer pool.Stop()

	// Submit job with maxRetries = 2 (initial attempt + 2 retries = 3 total attempts before permanent failure)
	ctx := context.Background()
	_, err := svc.SubmitJob(ctx, "critical_sync", map[string]string{"foo": "bar"}, 2)
	if err != nil {
		t.Fatalf("SubmitJob failed: %v", err)
	}

	// Poll until terminal state
	var failedJob *domain.Job
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		j, err := storage.GetByID(ctx, "failing-job-001")
		if err == nil && j.Status == domain.StatusFailed {
			failedJob = j
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	if failedJob == nil {
		t.Fatalf("job did not transition to StatusFailed after retry exhaustion")
	}
	if failedJob.Attempts != 3 {
		t.Fatalf("expected 3 total attempts (1 initial + 2 retries), got %d", failedJob.Attempts)
	}
	if failedJob.LastError == "" {
		t.Fatal("expected LastError to be populated")
	}
}
