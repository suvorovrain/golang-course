package adapter

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/artem-smola/golang-course/task2/collector/internal/domain"
)

type RepoInfo struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StarsCount  int       `json:"stargazers_count"`
	ForksCount  int       `json:"forks_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type CollectorAdapter struct{}

func (ca *CollectorAdapter) GetRepoInfo(owner, repoName string) (*domain.RepoInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repoName)

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-gh-tool")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status %d", resp.StatusCode)
	}

	var repoInfo RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&repoInfo); err != nil {
		return nil, err
	}

	return &domain.RepoInfo{
		Name:        repoInfo.Name,
		Description: repoInfo.Description,
		StarsCount:  repoInfo.StarsCount,
		ForksCount:  repoInfo.ForksCount,
		CreatedAt:   repoInfo.CreatedAt,
	}, nil
}
