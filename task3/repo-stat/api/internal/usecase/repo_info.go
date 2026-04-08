package usecase

import (
	"context"
	"repo-stat/api/internal/domain"
)

type ProcessorClient interface {
	Ping(context.Context) domain.PingStatus
	GetInfoRepo(context.Context, string, string) (*domain.RepoInfo, error)
}

type ApiGatewayUsecase struct {
	pc ProcessorClient
}

func NewUsecaseApiGateway(pc ProcessorClient) *ApiGatewayUsecase {
	return &ApiGatewayUsecase{
		pc: pc,
	}
}

func (agu *ApiGatewayUsecase) GetInfoRep(ctx context.Context, owner string, repo string) (*domain.RepoInfo, error) {
	return agu.pc.GetInfoRepo(ctx, owner, repo)
}

func (agu *ApiGatewayUsecase) Ping(ctx context.Context) domain.PingStatus {
	return agu.pc.Ping(ctx)
}
