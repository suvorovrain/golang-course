package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type RepoInfo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stargazers  uint64 `json:"stargazers_count"`
	Forks       uint64 `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

func GetRepoInfo(repo string, owner string) (*RepoInfo, error) {

	client := http.Client{}
	url := "https://api.github.com/repos/" + owner + "/" + repo

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "my-github-cli-tool")

	responce, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer responce.Body.Close()

	if responce.StatusCode != http.StatusOK {
		if responce.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("couldn't find repo")
		}
		return nil, fmt.Errorf("%s", responce.Status)
	}

	RepInfoSlice, err := io.ReadAll(responce.Body)
	if err != nil {
		return nil, err
	}

	var RepInfo RepoInfo
	err = json.Unmarshal(RepInfoSlice, &RepInfo)
	if err != nil {
		return nil, err
	}

	return &RepInfo, nil
}

func main() {

	var repo, owner string

	fmt.Scan(&repo, &owner)

	RepInfo, err := GetRepoInfo(repo, owner)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("FullName:%s\nDescription:%s\nStargazers:%d\nForks:%d\nCreatedAt:%s\n",
		RepInfo.FullName, RepInfo.Description, RepInfo.Stargazers, RepInfo.Forks, RepInfo.CreatedAt,
	)

}
