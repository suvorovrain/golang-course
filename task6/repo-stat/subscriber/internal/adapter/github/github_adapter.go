package github

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const requestTimeout = 10 * time.Second

type Adapter struct {
	client *http.Client
	log    *slog.Logger
}

func NewAdapter(log *slog.Logger) *Adapter {
	return &Adapter{
		client: &http.Client{
			Timeout: requestTimeout,
		},
		log: log,
	}
}

func (a *Adapter) IsRepoExist(ctx context.Context, owner, repo string) (bool, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	a.log.Info("github request", "url", url)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "repo-stat-subscriber/1.0")

	resp, err := a.client.Do(req)
	if err != nil {
		a.log.Error("github request failed", "error", err)
		return false, fmt.Errorf("github request failed: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			a.log.Error("failed to close response body", "error", err)
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %w", err)
	}

	a.log.Info("github response", "status", resp.StatusCode, "body", string(body))

	if resp.StatusCode == http.StatusNotFound {
		return false, nil
	}

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("github api error: status %d, body: %s", resp.StatusCode, string(body))
	}

	return true, nil
}
