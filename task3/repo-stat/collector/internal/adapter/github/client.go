package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"repo-stat/collector/internal/domain"
	"repo-stat/collector/internal/dto"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	httpClient *http.Client
	token      string
}

func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		token: token,
	}
}

func CreatePath(name string, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s", name, repo)
}

// GetRepo
// Get general info about GitHub repo by url.
func (c *Client) GetRepo(ctx context.Context, owner, name string) (domain.Repo, error) {
	url := CreatePath(owner, name)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return domain.Repo{}, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "GetRepoInfo-App")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return domain.Repo{}, fmt.Errorf("do request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:

	case http.StatusNotFound:
		return domain.Repo{}, status.Error(codes.NotFound, "repo not found")

	case http.StatusForbidden:
		return domain.Repo{}, status.Error(codes.ResourceExhausted, "access forbidden")

	case http.StatusUnauthorized:
		return domain.Repo{}, status.Error(codes.Unauthenticated, "unauthorized: invalid token")

	default:
		return domain.Repo{}, status.Error(codes.Internal, "unexpected status code")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.Repo{}, err
	}

	var dto dto.RepoDTO
	if err := json.Unmarshal(body, &dto); err != nil {
		return domain.Repo{}, err
	}

	return domain.Repo{
		Name:        dto.Name,
		Description: dto.Description,
		Stars:       dto.Stars,
		Forks:       dto.Forks,
		CreatedAt:   dto.CreatedAt,
	}, nil
}
