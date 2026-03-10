package main

import (
	"encoding/json"
	"fmt"
	"io"
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

func printRepoInfo(info RepoInfo) {
	fields := []struct {
		key string
		val string
	}{
		{"name", info.Name},
		{"description", func() string {
			if info.Description == "" {
				return "no description"
			}
			return info.Description
		}()},
		{"starred", fmt.Sprintf("%d", info.Stars)},
		{"forks", fmt.Sprintf("%d", info.Forks)},
		{"create date", func() string {
			createdStr := info.CreatedAt
			if t, err := time.Parse(time.RFC3339, info.CreatedAt); err == nil {
				createdStr = t.Format("02 Jan 2006 15:04:05")
			}
			return createdStr
		}()},
	}

	maxLen := 0
	for _, f := range fields {
		if l := len(f.key); l > maxLen {
			maxLen = l
		}
	}

	for _, f := range fields {
		fmt.Printf("%-*s %s\n", maxLen+3, f.key, f.val)
	}
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: owner repo")
		os.Exit(1)
	}

	owner := os.Args[1]
	repo := os.Args[2]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	client := &http.Client{Timeout: 1 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err response %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			fmt.Fprintf(os.Stderr, "err %s/%s not found \n", owner, repo)
		case http.StatusForbidden:
			fmt.Fprintf(os.Stderr, "err forbidden\n")
		default:
			fmt.Fprintf(os.Stderr, "err %s\n", resp.Status)
		}
		os.Exit(1)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "err responde %v\n", err)
		os.Exit(1)
	}

	var info RepoInfo
	if err := json.Unmarshal(body, &info); err != nil {
		fmt.Fprintf(os.Stderr, "err json %v\n", err)
		os.Exit(1)
	}

	printRepoInfo(info)
}
