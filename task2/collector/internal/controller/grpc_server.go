package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/artem-smola/golang-course/task2/collector/internal/domain"
	"github.com/artem-smola/golang-course/task2/proto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CollectorUsecase interface {
	Execute(owner, repoName string) (*domain.RepoInfo, error)
}

type CollectorGRPCServerController struct {
	gen.UnimplementedGRPCServiceServer
	useCase CollectorUsecase
}

func NewCollectorGRPCServerController(useCase CollectorUsecase) *CollectorGRPCServerController {
	return &CollectorGRPCServerController{useCase: useCase}
}

func (cc *CollectorGRPCServerController) GetRepoInfo(_ context.Context, req *gen.GetRepoInfoRequest) (*gen.GetRepoInfoResponse, error) {
	repoInfo, err := cc.useCase.Execute(req.GetOwner(), req.GetRepoName())
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get repo info: %s", err.Error()))
	}

	return &gen.GetRepoInfoResponse{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		StarsCount:  int64(repoInfo.StarsCount),
		ForksCount:  int64(repoInfo.ForksCount),
		CreatedAt:   repoInfo.CreatedAt.Format(time.RFC3339),
	}, nil
}
