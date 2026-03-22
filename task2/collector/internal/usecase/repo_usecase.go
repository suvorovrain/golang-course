package usecase

import (
	"context"
	"task2/collector/internal/ports"
	"task2/pkg/domain"
)

type RepoUsecase struct {
	githubClient ports.GithubClient
}

func NewRepoUsecase(client ports.GithubClient) *RepoUsecase {
	return &RepoUsecase{
		githubClient: client,
	}
}

func (u *RepoUsecase) Execute(ctx context.Context, owner, repo string) (domain.Repository, error) {
	repoData, err := u.githubClient.GetRepoInfo(ctx, owner, repo)
	if err != nil {
		return domain.Repository{}, err
	}

	return repoData, nil
}
