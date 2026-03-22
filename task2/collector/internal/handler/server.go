package handler

import (
	"context"
	"task2/collector/internal/usecase"
	"task2/proto"
	"time"
)

type Server struct {
	proto.UnimplementedRepositoryServiceServer
	usecase *usecase.RepoUsecase
}

func NewServer(u *usecase.RepoUsecase) *Server {
	return &Server{usecase: u}
}

func (s *Server) GetRepository(ctx context.Context, req *proto.GetRepoRequest) (*proto.RepoResponse, error) {
	repo, err := s.usecase.Execute(ctx, req.Owner, req.Repo)

	if err != nil {
		return nil, err
	}

	return &proto.RepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       int32(repo.Stars),
		Forks:       int32(repo.Forks),
		CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
	}, nil
}
