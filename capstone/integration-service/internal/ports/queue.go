package ports

import (
	"context"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

type JobQueue interface {
	Enqueue(ctx context.Context, job *domain.Job) error
	Dequeue(ctx context.Context) (*domain.Job, error)
	Close() error
}
