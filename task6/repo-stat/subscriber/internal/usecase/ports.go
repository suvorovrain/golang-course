package usecase

import (
	"context"
	"repo-stat/subscriber/internal/domain"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *domain.Subscription) error

	Delete(ctx context.Context, owner, repo string) error

	List(ctx context.Context) ([]*domain.Subscription, error)

	Exists(ctx context.Context, owner, repo string) (bool, error)
}
