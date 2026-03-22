package client

import (
	"collector/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	ErrInvalidURL   = errors.New("invalid repository URL")
	ErrRepoNotFound = errors.New("repository not found")
)

type UseCase struct {
	githubToken string
}

func NewUseCase(githubToken string) *UseCase {
	return &UseCase{
		githubToken: githubToken,
	}
}

func (uc *UseCase) Execute(ctx context.Context, repoURL string) (*domain.Repository, error) {
	owner, name, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidURL, err)
	}

	return uc.fetchFromGitHub(owner, name)
}

func (uc *UseCase) fetchFromGitHub(owner, repo string) (*domain.Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "golang-course-task2/1.0")

	if uc.githubToken != "" {
		req.Header.Set("Authorization", "token "+uc.githubToken)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrRepoNotFound
		}
		return nil, fmt.Errorf("GitHub API error: status %d", resp.StatusCode)
	}

	var githubRepo struct {
		FullName        string `json:"full_name"`
		Description     string `json:"description"`
		StargazersCount int    `json:"stargazers_count"`
		ForksCount      int    `json:"forks_count"`
		CreatedAt       string `json:"created_at"`
		Language        string `json:"language"`
	}

	if err := json.Unmarshal(body, &githubRepo); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, githubRepo.CreatedAt)

	return &domain.Repository{
		FullName:    githubRepo.FullName,
		Description: githubRepo.Description,
		Stars:       githubRepo.StargazersCount,
		Forks:       githubRepo.ForksCount,
		CreatedAt:   createdAt,
		Language:    githubRepo.Language,
	}, nil
}

func parseGitHubURL(repoURL string) (owner, repo string, err error) {
	repoURL = strings.TrimSpace(repoURL)
	repoURL = strings.TrimSuffix(repoURL, ".git")

	if !strings.HasPrefix(repoURL, "http://") && !strings.HasPrefix(repoURL, "https://") {
		repoURL = "https://" + repoURL
	}

	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", "", err
	}

	if parsedURL.Host != "github.com" {
		return "", "", fmt.Errorf("not a GitHub URL")
	}

	path := strings.Trim(parsedURL.Path, "/")
	parts := strings.Split(path, "/")

	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid format: expected /owner/repo")
	}

	return parts[0], parts[1], nil
}
