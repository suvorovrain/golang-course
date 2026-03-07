package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func main() {

	owner, repo := parseFlags()
	result, err := fetchRepoInfo(owner, repo)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	printInfo(result)
}

func parseFlags() (string, string) {
	owner := flag.String("owner", "", "Owner name")
	repo := flag.String("repo", "", "Repository name")
	flag.Parse()

	if *owner == "" || *repo == "" {
		fmt.Println("Usage: go run main.go [-owner owner] [-repo repo]")
		os.Exit(1)
	}

	return *owner, *repo
}

func fetchRepoInfo(owner, repo string) (*RepoInfo, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("network request failed: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("repository not found or unavailable (%d)\n", response.StatusCode)
	}

	result := RepoInfo{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

func printInfo(result *RepoInfo) {
	createdAt := makeDate(result.CreatedAt)

	fmt.Printf(`--- Repository Information ---
Name:          %s
Description:   %s
Stars:         %d
Forks:         %d
Created At:    %s
`, result.Name, result.Description, result.Stars, result.Forks, createdAt)

}

func makeDate(raw string) string {
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return raw
	}
	return t.Format("Jan 02, 2006")
}
