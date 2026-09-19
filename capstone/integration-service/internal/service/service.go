package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/ports"
)

type IntegrationService struct {
	storage ports.JobStorage
	queue   ports.JobQueue
	idGen   func() string
}

func NewIntegrationService(storage ports.JobStorage, queue ports.JobQueue) *IntegrationService {
	return &IntegrationService{
		storage: storage,
		queue:   queue,
		idGen: func() string {
			return fmt.Sprintf("job-%d", time.Now().UnixNano())
		},
	}
}

func (s *IntegrationService) SetIDGenerator(fn func() string) {
	s.idGen = fn
}

func (s *IntegrationService) SubmitJob(ctx context.Context, jobType string, payload map[string]string, maxRetries int) (*domain.Job, error) {
	jobID := s.idGen()
	job, err := domain.NewJob(jobID, jobType, payload, maxRetries)
	if err != nil {
		return nil, fmt.Errorf("creating job: %w", err)
	}

	if err := s.storage.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("saving job to storage: %w", err)
	}

	if err := s.queue.Enqueue(ctx, job); err != nil {
		return nil, fmt.Errorf("enqueuing job: %w", err)
	}

	return job, nil
}

func (s *IntegrationService) GetJob(ctx context.Context, id string) (*domain.Job, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *IntegrationService) ListJobs(ctx context.Context, limit int) ([]*domain.Job, error) {
	return s.storage.List(ctx, limit)
}
