package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"task2/pkg/domain"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GitHubAdapter struct {
	client *http.Client
}

func NewGitHubAdapter() *GitHubAdapter {
	return &GitHubAdapter{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

type githubResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func (a *GitHubAdapter) GetRepoInfo(ctx context.Context, owner, repo string) (domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.Repository{}, status.Error(codes.Internal, err.Error())
	}

	req.Header.Set("User-Agent", "task2-collector")

	resp, err := a.client.Do(req)
	if err != nil {
		return domain.Repository{}, status.Error(codes.Unavailable, "github api unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return domain.Repository{}, status.Error(codes.NotFound, "repository not found on github")
	}

	if resp.StatusCode != http.StatusOK {
		return domain.Repository{}, status.Errorf(codes.Internal, "github error: %d", resp.StatusCode)
	}

	var gResp githubResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return domain.Repository{}, status.Error(codes.Internal, "failed to decode github response")
	}

	createdAt, _ := time.Parse(time.RFC3339, gResp.CreatedAt)

	return domain.Repository{
		Name:        gResp.Name,
		Description: gResp.Description,
		Stars:       gResp.Stars,
		Forks:       gResp.Forks,
		CreatedAt:   createdAt,
	}, nil
}
