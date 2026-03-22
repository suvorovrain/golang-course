package usecase

import "github.com/artem-smola/golang-course/task2/gateway/internal/domain"

type GatewayAdapter interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type GatewayUsecase struct {
	adapter GatewayAdapter
}

func NewGatewayUsecase(adapter GatewayAdapter) *GatewayUsecase {
	return &GatewayUsecase{adapter: adapter}
}

func (gu *GatewayUsecase) Execute(owner, repoName string) (*domain.RepoInfo, error) {
	return gu.adapter.GetRepoInfo(owner, repoName)
}
