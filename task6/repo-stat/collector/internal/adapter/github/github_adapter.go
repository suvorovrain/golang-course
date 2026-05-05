package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"repo-stat/collector/internal/domain"
)

const requestTimeout = 10 * time.Second

type Adapter struct {
	client *http.Client
	log    *slog.Logger
}

type githubRepo struct {
	Owner       githubOwner `json:"owner"`
	Repo        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Description string      `json:"description"`
	ForksCount  int32       `json:"forks_count"`
	Stargazers  int32       `json:"stargazers_count"`
	CreatedAt   string      `json:"created_at"`
	Visibility  string      `json:"visibility"`
}

type githubOwner struct {
	Login string `json:"login"`
}

func NewAdapter(log *slog.Logger) *Adapter {
	return &Adapter{
		client: &http.Client{
			Timeout: requestTimeout,
		},
		log: log,
	}
}

func (a *Adapter) Get(ctx context.Context, owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	a.log.Info("github request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "repo-stat-collector/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		a.log.Error("github request failed", "error", err)
		return nil, fmt.Errorf("github request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.log.Error("failed to close response body", "error", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, domain.ErrRepoNotFound
	case http.StatusForbidden:
		return nil, domain.ErrGitHubRateLimited
	default:
		return nil, fmt.Errorf("%w: status %d, body: %s",
			domain.ErrGitHubAPIError, resp.StatusCode, string(body))
	}

	var gh githubRepo
	if err := json.Unmarshal(body, &gh); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return &domain.Repository{
		Owner:       gh.Owner.Login,
		Repo:        gh.Repo,
		FullName:    gh.FullName,
		Description: gh.Description,
		Stars:       gh.Stargazers,
		Forks:       gh.ForksCount,
		CreatedAt:   gh.CreatedAt,
		Visibility:  gh.Visibility,
	}, nil
}
