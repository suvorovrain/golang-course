package ports

import (
	"context"

	"task2/pkg/domain" // Используем общую доменную модель
)

type CollectorClient interface {
	GetRepository(ctx context.Context, owner, repo string) (domain.Repository, error)
}
