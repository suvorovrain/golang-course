package usecase

import (
	"context"
	"repo-stat/collector/internal/adapter/github"
	"repo-stat/collector/internal/adapter/subscriber"
	"repo-stat/collector/internal/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GetSubscriptionsInfo struct {
	subscriberClient *subscriber.Client
	githubAdapter    *github.Adapter
}

func NewGetSubscriptionsInfo(subClient *subscriber.Client, gh *github.Adapter) *GetSubscriptionsInfo {
	return &GetSubscriptionsInfo{
		subscriberClient: subClient,
		githubAdapter:    gh,
	}
}

func (gri *GetSubscriptionsInfo) GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error) {
	subs, err := gri.subscriberClient.GetSubscriptionsInfo(ctx)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return nil, domain.ErrRepoNotFound
			case codes.InvalidArgument:
				return nil, domain.ErrInvalidInput
			case codes.ResourceExhausted:
				return nil, domain.ErrGitHubRateLimited
			default:
				return nil, domain.ErrGitHubAPIError
			}
		}
		return nil, domain.ErrGitHubAPIError
	}

	repositories := &domain.SubscriptionInfo{Repositories: make([]*domain.Repository, 0, len(subs))}
	for _, sub := range subs {
		info, err := gri.githubAdapter.Get(ctx, sub.Owner, sub.Repo)
		if err != nil {
			continue
		}
		repositories.Repositories = append(repositories.Repositories, info)
	}

	return repositories, nil
}
