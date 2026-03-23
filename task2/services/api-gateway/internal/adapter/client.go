package adapter

import (
	"context"
	"golang-course/task2/pkg/api"
	"golang-course/task2/services/api-gateway/internal/domain"

	"google.golang.org/grpc"
)

type Client struct {
	api api.RepositoryServiceClient
}

func New(conn *grpc.ClientConn) *Client {
	return &Client{
		api: api.NewRepositoryServiceClient(conn),
	}
}

func (c *Client) GetRepo(ctx context.Context, owner, name string) (*domain.Repo, error) {
	resp, err := c.api.GetRepo(ctx, &api.GetRepoRequest{
		Owner: owner,
		Name:  name,
	})

	if err != nil {
		return nil, err
	}

	return &domain.Repo{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       resp.Stars,
		Forks:       resp.Forks,
		CreatedAt:   resp.CreatedAt,
	}, nil
}
