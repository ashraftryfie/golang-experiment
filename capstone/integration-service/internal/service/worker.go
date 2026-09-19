package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/adapters/memory"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/ports"
)

type WorkerPool struct {
	concurrency    int
	storage        ports.JobStorage
	queue          ports.JobQueue
	externalClient ports.ExternalAPIClient
	jobTimeout     time.Duration
	wg             sync.WaitGroup
	cancel         context.CancelFunc
}

func NewWorkerPool(
	concurrency int,
	storage ports.JobStorage,
	queue ports.JobQueue,
	client ports.ExternalAPIClient,
	jobTimeout time.Duration,
) *WorkerPool {
	if concurrency <= 0 {
		concurrency = 4
	}
	if jobTimeout <= 0 {
		jobTimeout = 5 * time.Second
	}
	return &WorkerPool{
		concurrency:    concurrency,
		storage:        storage,
		queue:          queue,
		externalClient: client,
		jobTimeout:     jobTimeout,
	}
}

func (w *WorkerPool) Start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	w.cancel = cancel
	for i := 0; i < w.concurrency; i++ {
		w.wg.Add(1)
		go w.runWorker(ctx, i)
	}
}

func (w *WorkerPool) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
}

func (w *WorkerPool) runWorker(ctx context.Context, id int) {
	defer w.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		job, err := w.queue.Dequeue(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, memory.ErrQueueClosed) {
				return
			}
			continue
		}

		w.processJob(ctx, job, id)
	}
}

func (w *WorkerPool) processJob(ctx context.Context, job *domain.Job, workerID int) {
	slog.InfoContext(ctx, "processing job",
		slog.String("job_id", job.ID),
		slog.Int("worker_id", workerID),
		slog.Int("attempt", job.Attempts+1),
	)

	_ = job.MarkProcessing()
	_ = w.storage.Save(ctx, job)

	jobCtx, cancel := context.WithTimeout(ctx, w.jobTimeout)
	defer cancel()

	err := w.externalClient.SyncJob(jobCtx, job)
	if err != nil {
		slog.WarnContext(ctx, "job execution failed",
			slog.String("job_id", job.ID),
			slog.String("error", err.Error()),
			slog.Int("attempt", job.Attempts),
		)
		job.MarkFailed(err.Error())
		_ = w.storage.Save(ctx, job)

		// If job failed but has remaining retries, re-enqueue
		if !job.IsTerminal() {
			_ = w.queue.Enqueue(ctx, job)
		}
		return
	}

	_ = job.MarkCompleted()
	_ = w.storage.Save(ctx, job)

	slog.InfoContext(ctx, "job completed successfully",
		slog.String("job_id", job.ID),
	)
}
