package grpc_server

import (
	"context"
	"golang-course/task2/pkg/api"
	"golang-course/task2/services/collector/internal/usecase"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	api.UnimplementedRepositoryServiceServer
	usecase *usecase.Usecase
}

func New(usecase *usecase.Usecase) *Server {
	return &Server{usecase: usecase}
}

func (s *Server) GetRepo(ctx context.Context, req *api.GetRepoRequest) (*api.GetRepoResponse, error) {
	if req.GetOwner() == "" || req.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "owner, name, required")
	}

	repo, err := s.usecase.GetRepoInfo(ctx, req.GetOwner(), req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &api.GetRepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}
