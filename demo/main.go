package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"sync"
	"text/tabwriter"
)

var apiEndpoint = "https://api.github.com"

type Owner struct {
	Login string `json:"login"`
}

type RepositoryInfo struct {
	Name         string `json:"name"`
	Owner        Owner  `json:"owner"`
	Description  string `json:"description"`
	Forks        int    `json:"forks"`
	Stargazers   int    `json:"stargazers_count"`
	CreatedAt    string `json:"created_at"`
	CommitsCount int
}

func main() {
	var owner string
	var repoName string

	if len(os.Args) == 3 {
		owner, repoName = os.Args[1], os.Args[2]
	} else {
		fmt.Println("Enter github repository owner and repository name (separate with space)")
		if _, err := fmt.Scan(&owner, &repoName); err != nil {
			fmt.Printf("Error reading input: %v\n", err)
			return
		}
	}

	var repoStruct RepositoryInfo

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := GetRepoInfo(owner, repoName, &repoStruct); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := GetCommitsCount(owner, repoName, &repoStruct); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	wg.Wait()

	PrintRepoInformation(&repoStruct)
}

func PrintRepoInformation(repo *RepositoryInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	_, _ = fmt.Fprintln(w, "--------------------------------------------------")
	_, _ = fmt.Fprintf(w, "🚀 REPOSITORY:\t%s/%s\n", repo.Owner.Login, repo.Name)
	_, _ = fmt.Fprintln(w, "--------------------------------------------------")

	if repo.Description != "" {
		_, _ = fmt.Fprintf(w, "📝 Description:\t%s\n", repo.Description)
	} else {
		_, _ = fmt.Fprintf(w, "📝 Description:\t[No description provided]\n")
	}

	_, _ = fmt.Fprintf(w, "⭐ Stars:\t%d\n", repo.Stargazers)
	_, _ = fmt.Fprintf(w, "🍴 Forks:\t%d\n", repo.Forks)
	_, _ = fmt.Fprintf(w, "📦 Commits:\t%d\n", repo.CommitsCount)
	_, _ = fmt.Fprintf(w, "📅 Created:\t%s\n", repo.CreatedAt)
	_, _ = fmt.Fprintln(w, "--------------------------------------------------")

	if err := w.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, "Can't flush buffer: %v", err)
	}
}

func GetRepoInfo(owner, repoName string, repoStruct *RepositoryInfo) error {
	url := fmt.Sprintf("%s/repos/%s/%s", apiEndpoint, owner, repoName)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(repoStruct)
	if err != nil {
		return fmt.Errorf("can't parse JSON: %w", err)
	}

	return nil
}

func GetCommitsCount(owner, repoName string, repoStruct *RepositoryInfo) error {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?per_page=1", apiEndpoint, owner, repoName)
	resp, err := http.Head(url)
	if err != nil {
		return err
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Can't close response body: %v\n", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned error: %d", resp.StatusCode)
	}

	link := resp.Header.Get("link")
	if link == "" {
		return fmt.Errorf("can't get commits")
	}

	re := regexp.MustCompile(`page=(\d+)>; rel="last"`)
	matches := re.FindStringSubmatch(link)

	if len(matches) > 1 {
		count, _ := strconv.Atoi(matches[1])
		repoStruct.CommitsCount = count
	}

	return nil
}
