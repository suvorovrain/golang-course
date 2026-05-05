package collector

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"repo-stat/processor/internal/domain"
	proto "repo-stat/proto/collector"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	grpc proto.CollectorClient
}

func NewClient(address string, log *slog.Logger) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		log:  log,
		conn: conn,
		grpc: proto.NewCollectorClient(conn),
	}, nil
}

func (c *Client) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	c.log.Info("calling GetRepoInfo on collector", "owner", owner, "repo", repo)

	req := &proto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	}
	resp, err := c.grpc.GetRepoInfo(ctx, req)
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
				return nil, fmt.Errorf("%w: %s", domain.ErrGitHubAPIError, st.Message())
			}
		}

		return nil, fmt.Errorf("collector: %w", err)
	}

	return &domain.Repository{
		Owner:       resp.Owner,
		Repo:        resp.Repo,
		FullName:    resp.FullName,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
		Visibility:  resp.Visibility,
	}, nil
}

func (c *Client) GetSubscriptionsInfo(ctx context.Context) (*domain.SubscriptionInfo, error) {
	resp, err := c.grpc.GetSubscriptionsInfo(ctx, &proto.GetSubscriptionsInfoRequest{})
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
				return nil, fmt.Errorf("%w: %s", domain.ErrGitHubAPIError, st.Message())
			}
		}
		return nil, fmt.Errorf("processor client: %w", err)
	}

	repositories := &domain.SubscriptionInfo{Repositories: make([]domain.Repository, 0, len(resp.Repositories))}
	for _, repo := range resp.Repositories {
		repositories.Repositories = append(repositories.Repositories, domain.Repository{
			Owner:       repo.Owner,
			Repo:        repo.Repo,
			FullName:    repo.FullName,
			CreatedAt:   repo.CreatedAt,
			Description: repo.Description,
			Visibility:  repo.Visibility,
			Stars:       repo.Stars,
			Forks:       repo.Forks,
		})
	}

	return repositories, nil
}
