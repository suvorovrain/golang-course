package usecase

import (
	"context"
	"fmt"

	"repo-stat/subscriber/internal/adapter/github"
	"repo-stat/subscriber/internal/domain"
)

type SubscriptionUseCase struct {
	repo   SubscriptionRepository
	github *github.Adapter
}

func NewSubscriptionUseCase(repo SubscriptionRepository, gh *github.Adapter) *SubscriptionUseCase {
	return &SubscriptionUseCase{
		repo:   repo,
		github: gh,
	}
}

func (uc *SubscriptionUseCase) Create(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("owner and repo cannot be empty")
	}

	exists, err := uc.github.IsRepoExist(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to check repository on GitHub: %w", err)
	}
	if !exists {
		return fmt.Errorf("repository %s/%s does not exist on GitHub", owner, repo)
	}

	sub := domain.NewSubscription(owner, repo)

	exists, err = uc.repo.Exists(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to check subscription existence: %w", err)
	}
	if exists {
		return fmt.Errorf("subscription for %s/%s already exists", owner, repo)
	}

	if err := uc.repo.Create(ctx, sub); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (uc *SubscriptionUseCase) Delete(ctx context.Context, owner, repo string) error {
	if owner == "" || repo == "" {
		return fmt.Errorf("owner and repo cannot be empty")
	}

	return uc.repo.Delete(ctx, owner, repo)
}

func (uc *SubscriptionUseCase) List(ctx context.Context) ([]*domain.Subscription, error) {
	return uc.repo.List(ctx)
}
