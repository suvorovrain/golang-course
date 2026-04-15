package grpcController

import (
	"context"
	"repo-stat/collector/internal/usecase"
	"repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RepoServer struct {
	collector.UnimplementedCollectorServiceServer
	usecase *usecase.RepoProvider
}

func NewRepoServer(usecase *usecase.RepoProvider) *RepoServer {
	return &RepoServer{
		usecase: usecase,
	}
}

func (s *RepoServer) GetRepo(ctx context.Context, req *collector.GetRepoRequest) (*collector.GetRepoResponse, error) {
	if req.GetOwner() == "" || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner, name, required")
	}

	repo, err := s.usecase.GetRepo(ctx, req.GetOwner(), req.GetName())
	if err != nil {
		return nil, err
	}

	return &collector.GetRepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}
