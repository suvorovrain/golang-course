package collectorhandler

import (
	"context"

	collectordomain "repo-stat/collector/internal/domain"

	collectorpb "repo-stat/proto/collector"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsecaseCollectorService interface {
	GetInfoRepo(context.Context, string, string) (*collectordomain.RepoInfo, error)
}

type colletorHandler struct {
	collectorpb.UnimplementedCollectorServer
	ucs UsecaseCollectorService
}

func NewHandler(ucs UsecaseCollectorService) *colletorHandler {
	return &colletorHandler{
		ucs: ucs,
	}
}

func (h *colletorHandler) GetInfoRepo(ctx context.Context, req *collectorpb.GetInfoRepoRequest) (*collectorpb.GetInfoRepoResponce, error) {

	RepoInfo, err := h.ucs.GetInfoRepo(ctx, req.GetOwner(), req.GetRepo())

	switch err {
	case collectordomain.ErrorNotFound:
		return nil, status.Error(codes.NotFound, codes.NotFound.String())
	case collectordomain.InternalError:
		return nil, status.Error(codes.Internal, codes.Internal.String())
	case collectordomain.BadRequest:
		return nil, status.Error(codes.InvalidArgument, codes.InvalidArgument.String())
	}

	return &collectorpb.GetInfoRepoResponce{
		Fullname:    RepoInfo.FullName,
		Description: RepoInfo.Description,
		Stargazers:  RepoInfo.Stargazers,
		Forks:       RepoInfo.Forks,
		Createdat:   RepoInfo.CreatedAt,
	}, nil
}
