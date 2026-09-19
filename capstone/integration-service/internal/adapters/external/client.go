package external

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

var ErrExternalAPI = errors.New("external downstream integration failure")

type ResilientExternalClient struct {
	mu           sync.Mutex
	Latency      time.Duration
	FailForJobID string
	SyncedJobs   []string
}

func NewResilientExternalClient() *ResilientExternalClient {
	return &ResilientExternalClient{
		Latency:    5 * time.Millisecond,
		SyncedJobs: make([]string, 0),
	}
}

func (c *ResilientExternalClient) SyncJob(ctx context.Context, job *domain.Job) error {
	if c.Latency > 0 {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(c.Latency):
		}
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if c.FailForJobID != "" && job.ID == c.FailForJobID {
		return ErrExternalAPI
	}

	c.SyncedJobs = append(c.SyncedJobs, job.ID)
	return nil
}

func (c *ResilientExternalClient) GetSyncedJobs() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	res := make([]string, len(c.SyncedJobs))
	copy(res, c.SyncedJobs)
	return res
}
