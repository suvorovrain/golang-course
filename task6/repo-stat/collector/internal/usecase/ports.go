package usecase

import (
	"context"
	"repo-stat/collector/internal/domain"
)

type RepoGetter interface {
	Get(ctx context.Context, owner, repo string) (*domain.Repository, error)
}
