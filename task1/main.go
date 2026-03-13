package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type repoInfo struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Stars       int     `json:"stargazers_count"`
	Forks       int     `json:"forks_count"`
	CreatedAt   string  `json:"created_at"`
}

func parseInput(input string) (owner, repo string, err error) {
	input = strings.TrimSpace(input)

	idx := strings.Index(input, "github.com/")
	if idx != -1 {
		input = input[idx+len("github.com/"):]
	}

	parts := strings.Split(input, "/")

	owner = parts[0]
	repo = parts[1]

	if owner == "" || repo == "" {
		return "", "", fmt.Errorf("Broken repo name\n")
	}

	return owner, repo, nil
}

func parseTime(date string) string {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return date
	}
	return t.Format("02.01.2006 15:04:05") + " UTC"
}

func main() {
	if len(os.Args) != 2 {
		_, err := fmt.Fprintf(os.Stderr, "Usage: %s <repo>\n", os.Args[0])
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}

	owner, repo, err := parseInput(os.Args[1])
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error parse input: %v\n", err)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Request error: %v\n", err)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}

	req.Header.Set("User-Agent", "CLI repo info utility written in Go")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Network error: %v\n", err)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Printf("Error: %v\n", err)
		}
	}()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		_, err := fmt.Fprintf(os.Stderr, "Repo '%s/%s' does not found\n", owner, repo)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	default:
		_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", resp.Status)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}

	var info repoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		os.Exit(1)
	}

	description := "-"
	if info.Description != nil && *info.Description != "" {
		description = *info.Description
	}

	fmt.Printf("Repo \033[33m%s/%s\033[0m info:\n", owner, repo)
	fmt.Println("Name:        ", info.Name)
	fmt.Println("Description: ", description)
	fmt.Println("Stars:       ", info.Stars)
	fmt.Println("Forks:       ", info.Forks)
	fmt.Println("Date:        ", parseTime(info.CreatedAt))
}
