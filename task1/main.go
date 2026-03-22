package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type repoInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StarsNumber int    `json:"stargazers_count"`
	ForkNumber  int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func main() {
	client := &http.Client{}

	owner := os.Args[1]
	repo := os.Args[2]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	resp, err := client.Get(url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Network error: %v\n", err)
		os.Exit(1)
	}

	if resp.StatusCode == 404 {
		fmt.Fprintf(os.Stderr, "Repository not found\n")
		os.Exit(1)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	var repos repoInfo
	err = json.Unmarshal(body, &repos)

	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON parse error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("=== Repo Info ===")
	fmt.Printf("Name: %s\n", repos.Name)
	fmt.Printf("Description: %s\n", repos.Description)
	fmt.Printf("Stars: %d\n", repos.StarsNumber)
	fmt.Printf("Forks: %d\n", repos.ForkNumber)
	fmt.Printf("Created: %s\n", repos.CreatedAt)

}
