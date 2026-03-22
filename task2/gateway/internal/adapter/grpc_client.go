package adapter

import (
	"context"
	"time"

	"github.com/artem-smola/golang-course/task2/gateway/internal/domain"
	"github.com/artem-smola/golang-course/task2/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GatewayGRPCClientAdapter struct {
	client gen.GRPCServiceClient
}

func NewGatewayGRPCClientAdapter(addr string) (*GatewayGRPCClientAdapter, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := gen.NewGRPCServiceClient(conn)
	return &GatewayGRPCClientAdapter{client: client}, nil
}

func (ga *GatewayGRPCClientAdapter) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := ga.client.GetRepoInfo(ctx, &gen.GetRepoInfoRequest{
		Owner:    owner,
		RepoName: repoName,
	})
	if err != nil {
		return nil, err
	}
	return &domain.RepoInfo{
		Name:        resp.GetName(),
		Description: resp.GetDescription(),
		StarsCount:  int(resp.GetStarsCount()),
		ForksCount:  int(resp.GetForksCount()),
		CreatedAt:   resp.GetCreatedAt(),
	}, nil
}
