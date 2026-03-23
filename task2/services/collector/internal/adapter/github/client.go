package github

import (
	"context"
	"encoding/json"
	"fmt"
	"golang-course/task2/services/collector/internal/domain"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		token: token,
	}
}

type RepoDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int64  `json:"stargazers_count"`
	Forks       int64  `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func CreatePath(name string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", name, repo)
}

// GetRepo
// Get general info about GitHub repo by url.
func (c *Client) GetRepo(ctx context.Context, owner, name string) (*domain.Repo, error) {
	url := CreatePath(owner, name)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "GetRepoInfo-App")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:

		return nil, fmt.Errorf("repo not found: %s", url)

	case http.StatusForbidden:
		return nil, fmt.Errorf("access is not accepted: %s", url)

	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var dto RepoDTO

	if err := json.Unmarshal(body, &dto); err != nil {
		return nil, err
	}

	return &domain.Repo{
		Name:        dto.Name,
		Description: dto.Description,
		Stars:       dto.Stars,
		Forks:       dto.Forks,
		CreatedAt:   dto.CreatedAt,
	}, nil
}
