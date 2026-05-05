package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/processor/internal/domain"
	repo "repo-stat/processor/internal/usecase"
	proto "repo-stat/proto/processor"
)

type Handler struct {
	proto.UnimplementedProcessorServer

	usecase *repo.GetRepoUseCase
}

func NewHandler(usecase *repo.GetRepoUseCase) *Handler {
	return &Handler{usecase: usecase}
}

func (h *Handler) GetRepository(ctx context.Context, req *proto.GetRepoRequest) (*proto.GetRepoResponse, error) {
	repoData, err := h.usecase.Get(ctx, req.Owner, req.Repo)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrRepoNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case domain.ErrGitHubRateLimited:
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}

	return &proto.GetRepoResponse{
		Owner:       repoData.Owner,
		Repo:        repoData.Repo,
		FullName:    repoData.FullName,
		Description: repoData.Description,
		Stars:       int32(repoData.Stars),
		Forks:       int32(repoData.Forks),
		CreatedAt:   repoData.CreatedAt,
		Visibility:  repoData.Visibility,
	}, nil
}

func (h *Handler) GetSubscriptionsInfo(ctx context.Context, req *proto.GetSubscriptionsInfoRequest) (*proto.GetSubscriptionsInfoResponse, error) {
	resp, err := h.usecase.GetSubscriptionsInfo(ctx)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case domain.ErrRepoNotFound:
			return nil, status.Error(codes.NotFound, err.Error())
		case domain.ErrGitHubRateLimited:
			return nil, status.Error(codes.ResourceExhausted, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error: "+err.Error())
		}
	}

	repositories := proto.GetSubscriptionsInfoResponse{Repositories: make([]*proto.GetRepoResponse, 0, len(resp.Repositories))}
	for _, repo := range resp.Repositories {
		repositories.Repositories = append(repositories.Repositories, &proto.GetRepoResponse{
			Owner:       repo.Owner,
			Repo:        repo.Repo,
			FullName:    repo.FullName,
			Description: repo.Description,
			Stars:       int32(repo.Stars),
			Forks:       int32(repo.Forks),
			CreatedAt:   repo.CreatedAt,
			Visibility:  repo.Visibility,
		})
	}

	return &repositories, nil
}

func (h *Handler) Ping(ctx context.Context, req *proto.PingRequest) (*proto.PingResponse, error) {
	return &proto.PingResponse{Reply: "pong from processor"}, nil
}
