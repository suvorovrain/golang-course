package processor_domain

import (
	"context"
	"errors"
)

type RepoInfo struct {
	FullName    string
	Description string
	Stargazers  uint64
	Forks       uint64
	CreatedAt   string
}

type Ping struct {
	Reply string
}

var (
	ErrorNotFound = errors.New("NOT_FOUND")
	InternalError = errors.New("INTERNAL_ERROR")
	BadRequest    = errors.New("BAD_REQUEST")
)

type CollectorClient interface {
	GetRepoInfo(context.Context, string, string) (*RepoInfo, error)
	Ping(context.Context) (*Ping, error)
}
