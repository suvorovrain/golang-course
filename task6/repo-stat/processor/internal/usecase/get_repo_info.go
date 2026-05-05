package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"repo-stat/processor/internal/domain"
)

type GetRepoUseCase struct {
	repo     Repository
	producer MessageProducer
	log      *slog.Logger
}

func NewGetRepoUseCase(repo Repository, producer MessageProducer, log *slog.Logger) *GetRepoUseCase {
	return &GetRepoUseCase{
		repo:     repo,
		producer: producer,
		log:      log,
	}
}

func (uc *GetRepoUseCase) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	if owner == "" || repo == "" {
		return nil, domain.ErrInvalidInput
	}

	uc.log.Info("Get: checking cache", "owner", owner, "repo", repo)

	cached, err := uc.repo.GetRepoFromCache(ctx, owner, repo)
	if err == nil && cached != nil {
		uc.log.Info("Get: cache hit", "owner", owner, "repo", repo, "full_name", cached.FullName)
		return cached, nil
	}

	uc.log.Info("Get: cache miss, publishing request", "owner", owner, "repo", repo)

	if err := uc.producer.PublishRepoRequest(ctx, owner, repo); err != nil {
		return nil, fmt.Errorf("failed to publish request to kafka: %w", err)
	}

	return &domain.Repository{
		Owner: owner,
		Repo:  repo,
	}, nil
}

func (uc *GetRepoUseCase) GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error) {
	uc.log.Info("GetSubscriptionsInfo: listing subscriptions")

	subs, err := uc.repo.ListSubscriptions(ctx)
	if err != nil {
		uc.log.Error("GetSubscriptionsInfo: failed to list subscriptions", "error", err)
		return nil, err
	}

	uc.log.Info("GetSubscriptionsInfo: got subscriptions", "count", len(subs))

	var repos []domain.Repository

	for _, sub := range subs {
		uc.log.Info("GetSubscriptionsInfo: checking cache", "owner", sub.Owner, "repo", sub.Repo)

		cached, err := uc.repo.GetRepoFromCache(ctx, sub.Owner, sub.Repo)
		if err == nil && cached != nil {
			uc.log.Info("GetSubscriptionsInfo: cache hit", "owner", sub.Owner, "repo", sub.Repo)
			repos = append(repos, *cached)
			continue
		}

		uc.log.Info("GetSubscriptionsInfo: cache miss, publishing request", "owner", sub.Owner, "repo", sub.Repo)

		_ = uc.producer.PublishRepoRequest(ctx, sub.Owner, sub.Repo)

		repos = append(repos, domain.Repository{
			Owner: sub.Owner,
			Repo:  sub.Repo,
		})
	}

	return &domain.SubscriptionInfo{Repositories: repos}, nil
}
