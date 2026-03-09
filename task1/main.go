package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	nameUser := flag.String("name", "", "User's name")
	nameRepo := flag.String("repo", "", "Repository's name")

	flag.Parse()

	if *nameUser == "" || *nameRepo == "" {
		fmt.Println("Error: both -name and -repo flags are required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Fetching data for: %s/%s\n", *nameUser, *nameRepo)
	infoRepo, err := getRepoInfo(*nameUser, *nameRepo)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Repository: %s\n Stars: %d\n Description: %s\n Forks: %d\n Date created: %s\n",
		infoRepo.FullName, infoRepo.Stars, infoRepo.Description, infoRepo.Forks, infoRepo.CreatedAt)
}

type Repository struct {
	FullName    string `json:"full_name"`
	Stars       int    `json:"stargazers_count"`
	Description string `json:"description"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func getRepoInfo(nameUser, nameRepo string) (*Repository, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", nameUser, nameRepo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "MyGoApp/1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var repo Repository
	err = json.Unmarshal(body, &repo)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &repo, nil
}
