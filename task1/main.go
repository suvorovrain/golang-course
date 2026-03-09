package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type RepoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	Date        string `json:"created_at"`
}

func request(owner string, repo string) (*RepoInfo, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "go-client")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		switch resp.StatusCode {
		case http.StatusNotFound:
			return nil, fmt.Errorf("repository not found")
		case http.StatusForbidden:
			return nil, fmt.Errorf("access forbidden or rate limit exceeded")
		default:
			return nil, fmt.Errorf("unexpected http status: %s", resp.Status)
		}
	}

	var res RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode json: %w", err)
	}

	return &res, nil
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Incorrect input")
		return
	}
	res, err := request(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Repository: %s\nDescription: %s\nStars: %d\nForks: %d\nCreated: %s\n",
		res.Name, res.Description, res.Stars, res.Forks, res.Date)
}
