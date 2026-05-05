package usecase

import (
	"context"
	"repo-stat/processor/internal/domain"
)

type Repository interface {
	ListSubscriptions(ctx context.Context) ([]*domain.Subscription, error)
	ReplaceAllSubscriptions(ctx context.Context, subs []*domain.Subscription) error

	GetRepoFromCache(ctx context.Context, owner, repo string) (*domain.Repository, error)
	UpsertRepoCache(ctx context.Context, repo *domain.Repository) error
}

type MessageProducer interface {
	PublishRepoRequest(ctx context.Context, owner, repo string) error
}
