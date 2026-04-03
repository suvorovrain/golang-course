package grpcController

import (
	"context"
	"repo-stat/processor/internal/usecase"
	"repo-stat/proto/processor"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProcessorServer struct {
	processor.UnimplementedProcessorServiceServer
	ping *usecase.Ping
	repo *usecase.Repo
}

func NewProcessorServer(ping *usecase.Ping, repo *usecase.Repo) *ProcessorServer {
	return &ProcessorServer{
		ping: ping,
		repo: repo,
	}
}

func (s *ProcessorServer) Ping(ctx context.Context, req *processor.PingRequest) (*processor.PingResponse, error) {
	status, err := s.ping.Ping(ctx)
	if err != nil {
		return &processor.PingResponse{Reply: "down"}, err
	}

	return &processor.PingResponse{
		Reply: string(status),
	}, nil
}

func (s *ProcessorServer) GetRepo(ctx context.Context, req *processor.GetRepoRequest) (*processor.GetRepoResponse, error) {
	if req.GetUrl() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "empty repo url")
	}

	repo, err := s.repo.GetRepoInfo(ctx, req.GetUrl())
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return nil, st.Err()
		}

		return nil, status.Errorf(codes.Internal, "%s", err.Error())
	}

	return &processor.GetRepoResponse{
		Name:        repo.Name,
		Description: repo.Description,
		Stars:       repo.Stars,
		Forks:       repo.Forks,
		CreatedAt:   repo.CreatedAt,
	}, nil
}
