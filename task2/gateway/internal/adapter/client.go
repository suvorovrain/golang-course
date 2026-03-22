package adapter

import (
	"context"
	"task2/pkg/domain"
	"task2/proto"
	"time"

	"google.golang.org/grpc"
)

type GRPCCollectorAdapter struct {
	client proto.RepositoryServiceClient
}

func NewGRPCCollectorAdapter(conn *grpc.ClientConn) *GRPCCollectorAdapter {
	return &GRPCCollectorAdapter{
		client: proto.NewRepositoryServiceClient(conn),
	}
}

func (a *GRPCCollectorAdapter) GetRepository(ctx context.Context, owner, repo string) (domain.Repository, error) {
	resp, err := a.client.GetRepository(ctx, &proto.GetRepoRequest{
		Owner: owner,
		Repo:  repo,
	})
	if err != nil {
		return domain.Repository{}, err
	}

	createdAt, _ := time.Parse(time.RFC3339, resp.CreatedAt)

	return domain.Repository{
		Name:        resp.Name,
		Description: resp.Description,
		Stars:       int(resp.Stars),
		Forks:       int(resp.Forks),
		CreatedAt:   createdAt,
	}, nil
}
