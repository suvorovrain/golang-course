package processor

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	"repo-stat/api/internal/domain"
	processorProto "repo-stat/proto/processor"
)

type Client struct {
	log  *slog.Logger
	conn *grpc.ClientConn
	pb   processorProto.ProcessorClient
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
		pb:   processorProto.NewProcessorClient(conn),
	}, nil
}

func (c *Client) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	c.log.Info("calling GetRepository on processor", "owner", owner, "repo", repo)

	resp, err := c.pb.GetRepository(ctx, &processorProto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	})
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
	c.log.Info("calling GetSubscriptionsInfo on processor", "address", c.conn.Target())

	resp, err := c.pb.GetSubscriptionsInfo(ctx, &processorProto.GetSubscriptionsInfoRequest{})
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

func (c *Client) Ping(ctx context.Context) domain.PingStatus {
	_, err := c.pb.Ping(ctx, &processorProto.PingRequest{})
	if err != nil {
		return domain.PingStatusDown
	}
	return domain.PingStatusUp
}
