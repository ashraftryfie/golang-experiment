package ports

import (
	"context"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

type JobStorage interface {
	Save(ctx context.Context, job *domain.Job) error
	GetByID(ctx context.Context, id string) (*domain.Job, error)
	List(ctx context.Context, limit int) ([]*domain.Job, error)
}
