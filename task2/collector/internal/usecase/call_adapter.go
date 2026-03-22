package usecase

import "github.com/artem-smola/golang-course/task2/collector/internal/domain"

type CollectorAdapter interface {
	GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error)
}

type CollectorUsecase struct {
	adapter CollectorAdapter
}

func NewCollectorUsecase(adapter CollectorAdapter) *CollectorUsecase {
	return &CollectorUsecase{adapter: adapter}
}

func (cu *CollectorUsecase) Execute(owner, repoName string) (*domain.RepoInfo, error) {
	return cu.adapter.GetRepoInfo(owner, repoName)
}
