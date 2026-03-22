package ports

import (
	"context"
	"task2/pkg/domain"
)

type GithubClient interface {
	GetRepoInfo(ctx context.Context, owner, repo string) (domain.Repository, error)
}
