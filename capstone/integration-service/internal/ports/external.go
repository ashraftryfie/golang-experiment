package ports

import (
	"context"

	"github.com/ashraftryfie/golang-experiment/capstone/integration-service/internal/domain"
)

type ExternalAPIClient interface {
	SyncJob(ctx context.Context, job *domain.Job) error
}
