package processor_usecase

import (
	"context"
	processor_domain "repo-stat/processor/internal/domain"
)

type processorService struct {
	cc processor_domain.CollectorClient
}

func NewProcessorService(collc processor_domain.CollectorClient) *processorService {
	return &processorService{
		cc: collc,
	}
}

func (ps *processorService) GetRepoInfo(ctx context.Context, owner string, repo string) (*processor_domain.RepoInfo, error) {
	return ps.cc.GetRepoInfo(ctx, owner, repo)
}

func (ps *processorService) Ping(ctx context.Context) (*processor_domain.Ping, error) {
	return ps.cc.Ping(ctx)
}
