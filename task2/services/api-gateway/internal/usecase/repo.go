package usecase

import (
	"context"
	"golang-course/task2/services/api-gateway/internal/domain"
)

type RepoProvider interface {
	GetRepo(ctx context.Context, owner, name string) (*domain.Repo, error)
}

type Usecase struct {
	provider RepoProvider
}

func New(p RepoProvider) *Usecase {
	return &Usecase{provider: p}
}

func (u *Usecase) GetRepoInfo(ctx context.Context, owner, name string) (*domain.Repo, error) {
	return u.provider.GetRepo(ctx, owner, name)
}
