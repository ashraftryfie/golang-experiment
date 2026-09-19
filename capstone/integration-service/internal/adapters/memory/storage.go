package memory

import (
	"context"
	"sync"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

type InMemoryJobStorage struct {
	mu   sync.RWMutex
	jobs map[string]*domain.Job
}

func NewInMemoryJobStorage() *InMemoryJobStorage {
	return &InMemoryJobStorage{
		jobs: make(map[string]*domain.Job),
	}
}

func (s *InMemoryJobStorage) Save(_ context.Context, job *domain.Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clone := cloneJob(job)
	s.jobs[job.ID] = clone
	return nil
}

func (s *InMemoryJobStorage) GetByID(_ context.Context, id string) (*domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	job, exists := s.jobs[id]
	if !exists {
		return nil, domain.ErrJobNotFound
	}
	return cloneJob(job), nil
}

func (s *InMemoryJobStorage) List(_ context.Context, limit int) ([]*domain.Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*domain.Job
	count := 0
	for _, job := range s.jobs {
		if limit > 0 && count >= limit {
			break
		}
		result = append(result, cloneJob(job))
		count++
	}
	return result, nil
}

func cloneJob(j *domain.Job) *domain.Job {
	payloadCopy := make(map[string]string, len(j.Payload))
	for k, v := range j.Payload {
		payloadCopy[k] = v
	}

	clone := *j
	clone.Payload = payloadCopy
	if j.CompletedAt != nil {
		t := *j.CompletedAt
		clone.CompletedAt = &t
	}
	return &clone
}
