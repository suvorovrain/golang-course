package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"collector/internal/domain"
)

var (
	ErrRepoNotFound = errors.New("repository not found")
	ErrRateLimit    = errors.New("GitHub API rate limit exceeded")
)

type Config struct {
	BaseURL   string
	UserAgent string
	AuthToken string
	Timeout   time.Duration
}

type Client struct {
	config     Config
	httpClient *http.Client
}

func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.github.com"
	}
	if config.UserAgent == "" {
		config.UserAgent = "golang-course-task2/1.0"
	}
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

func (c *Client) GetRepository(owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("%s/repos/%s/%s", c.config.BaseURL, owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", c.config.UserAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	if c.config.AuthToken != "" {
		req.Header.Set("Authorization", "token "+c.config.AuthToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return c.handleErrorResponse(resp.StatusCode, body, owner, repo)
	}

	var githubRepo GitHubResponse
	if err := json.Unmarshal(body, &githubRepo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return c.toDomain(githubRepo), nil
}

func (c *Client) handleErrorResponse(statusCode int, body []byte, owner, repo string) (*domain.Repository, error) {
	switch statusCode {
	case http.StatusNotFound:
		return nil, fmt.Errorf("%w: %s/%s", ErrRepoNotFound, owner, repo)

	case http.StatusForbidden:
		if contains(string(body), "rate limit") {
			return nil, ErrRateLimit
		}
		return nil, fmt.Errorf("access forbidden to %s/%s", owner, repo)

	case http.StatusUnauthorized:
		return nil, fmt.Errorf("authentication required")

	default:
		return nil, fmt.Errorf("GitHub API error: status %d, response: %s", statusCode, string(body))
	}
}

func (c *Client) toDomain(githubRepo GitHubResponse) *domain.Repository {
	createdAt, err := time.Parse(time.RFC3339, githubRepo.CreatedAt)
	if err != nil {
		createdAt = time.Time{}
	}

	return &domain.Repository{
		FullName:    githubRepo.FullName,
		Description: githubRepo.Description,
		Stars:       githubRepo.StargazersCount,
		Forks:       githubRepo.ForksCount,
		CreatedAt:   createdAt,
		Language:    githubRepo.Language,
	}
}

type GitHubResponse struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description"`
	StargazersCount int    `json:"stargazers_count"`
	ForksCount      int    `json:"forks_count"`
	CreatedAt       string `json:"created_at"`
	Language        string `json:"language"`
	HTMLURL         string `json:"html_url"`
	DefaultBranch   string `json:"default_branch"`
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
