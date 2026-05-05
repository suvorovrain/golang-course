package controller

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"repo-stat/collector/internal/domain"
	repo "repo-stat/collector/internal/usecase"
	proto "repo-stat/proto/collector"
)

type Handler struct {
	proto.UnimplementedCollectorServer

	repoUsecase         *repo.GetRepoInfo
	subscriptionUsecase *repo.GetSubscriptionsInfo
}

func NewHandler(repoUC *repo.GetRepoInfo, subUC *repo.GetSubscriptionsInfo) *Handler {
	return &Handler{repoUsecase: repoUC, subscriptionUsecase: subUC}
}

func (h *Handler) GetRepoInfo(ctx context.Context, req *proto.GetRepoRequest) (*proto.GetRepoResponse, error) {
	repoData, err := h.repoUsecase.Get(ctx, req.Owner, req.Repo)
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
	subscriptions, err := h.subscriptionUsecase.GetSubscriptionsInfo(ctx)
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

	repositories := proto.GetSubscriptionsInfoResponse{
		Repositories: make([]*proto.GetRepoResponse, 0, len(subscriptions.Repositories)),
	}
	for _, sub := range subscriptions.Repositories {
		repositories.Repositories = append(repositories.Repositories,
			&proto.GetRepoResponse{
				Owner:       sub.Owner,
				Repo:        sub.Repo,
				FullName:    sub.FullName,
				Description: sub.Description,
				Stars:       int32(sub.Stars),
				Forks:       int32(sub.Forks),
				CreatedAt:   sub.CreatedAt,
				Visibility:  sub.Visibility,
			})
	}

	return &repositories, nil
}
