package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Stargazers  int       `json:"stargazers_count"`
	Forks       int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
	HTMLURL     string    `json:"html_url"`
}

func (r *Repository) String() string {
	return fmt.Sprintf("=== Repository info ===\n"+
		"Name: %s\n"+
		"Description: %s\n"+
		"Stars: %d\n"+
		"Forks: %d\n"+
		"Created at: %s\n"+
		"URL: %s\n"+
		"========================",
		r.Name, r.Description, r.Stargazers, r.Forks,
		r.CreatedAt.Format("02-01-2006 15:04:05"), r.HTMLURL)
}

type GitHubClient struct {
	httpClient *http.Client
	baseURL    string
	logger     *log.Logger
}

func NewClient() *GitHubClient {
	return &GitHubClient{
		httpClient: &http.Client{Timeout: 10 * time.Second},
		baseURL:    "https://api.github.com/repos",
		logger:     log.Default(),
	}
}

func (ghc GitHubClient) GetRepositoryInfo(url string) (*Repository, error) {
	repoPath, err := extractRepoPath(url)
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("%s/%s", ghc.baseURL, repoPath)

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation error: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Go-Client")

	resp, err := ghc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request execution error: %w", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			ghc.logger.Printf("Warning: failed to close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var repo Repository
	if err := json.NewDecoder(resp.Body).Decode(&repo); err != nil {
		return nil, fmt.Errorf("json decoding error: %w", err)
	}

	return &repo, nil
}

func extractRepoPath(input string) (string, error) {
	if !strings.Contains(input, "https://") {
		parts := strings.Split(input, "/")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid format")
		} else if parts[0] == "" || parts[1] == "" {
			return "", fmt.Errorf("owner or repo is empty")
		}

		return input, nil
	}

	return extractFromURL(input)
}

func extractFromURL(url string) (string, error) {
	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "github.com/")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid GitHub URL format")
	}

	path := parts[1]
	path = strings.TrimPrefix(path, ":")
	path = strings.Trim(path, "/")

	repoParts := strings.Split(path, "/")
	if len(repoParts) < 2 {
		return "", fmt.Errorf("could not extract repository owner and name")
	}

	return strings.Join(repoParts[:2], "/"), nil
}
