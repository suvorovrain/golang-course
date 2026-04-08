package collectorusecase

import (
	"context"

	collectordomain "repo-stat/collector/internal/domain"
)

type collectorService struct {
	ghc collectordomain.GitHubClient
}

func NewCollectorService(ghc collectordomain.GitHubClient) *collectorService {
	return &collectorService{
		ghc: ghc,
	}
}

func (cs *collectorService) GetInfoRepo(ctx context.Context, owner string, repo string) (*collectordomain.RepoInfo, error) {
	RepoInfo, err := cs.ghc.GetRepoInfo(ctx, owner, repo)

	if err != nil {
		return nil, err
	}

	return &collectordomain.RepoInfo{
		FullName:    RepoInfo.FullName,
		Description: RepoInfo.Description,
		Stargazers:  RepoInfo.Stargazers,
		Forks:       RepoInfo.Forks,
		CreatedAt:   RepoInfo.CreatedAt,
	}, nil
}
