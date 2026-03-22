package controller

import (
	"errors"
	"net/url"
	"strings"
)

var errInvalidGitHubRepoURL = errors.New("invalid github repository url")

func ParseGitHubRepoURL(raw string) (string, string, error) {
	normalized := strings.TrimSpace(raw)
	if normalized == "" {
		return "", "", errInvalidGitHubRepoURL
	}

	parsed, err := url.Parse(normalized)
	if err != nil {
		return "", "", errInvalidGitHubRepoURL
	}
	if parsed.Scheme == "" && parsed.Host == "" {
		parsed, err = url.Parse("https://" + normalized)
		if err != nil {
			return "", "", errInvalidGitHubRepoURL
		}
	}
	if parsed.Scheme != "" && parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", "", errInvalidGitHubRepoURL
	}

	host := strings.ToLower(parsed.Host)
	if host == "www.github.com" {
		host = "github.com"
	}
	if host != "github.com" {
		return "", "", errInvalidGitHubRepoURL
	}

	path := strings.Trim(parsed.Path, "/")
	if path == "" {
		return "", "", errInvalidGitHubRepoURL
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "", "", errInvalidGitHubRepoURL
	}

	owner := strings.TrimSpace(parts[0])
	repo := strings.TrimSpace(parts[1])
	repo = strings.TrimSuffix(repo, ".git")
	if owner == "" || repo == "" {
		return "", "", errInvalidGitHubRepoURL
	}

	return owner, repo, nil
}
