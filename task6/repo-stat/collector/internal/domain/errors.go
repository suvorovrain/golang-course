package domain

import "errors"

var (
	ErrRepoNotFound      = errors.New("repository not found")
	ErrInvalidInput      = errors.New("invalid owner or repository name")
	ErrGitHubAPIError    = errors.New("github api returned an error")
	ErrGitHubRateLimited = errors.New("github rate limit exceeded")
)
