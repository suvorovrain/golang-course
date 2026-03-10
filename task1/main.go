package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Repository struct {
	Name        string    `json:"name"`
	Visibility  string    `json:"visibility"`
	Description string    `json:"description"`
	Forks       int       `jsom:"forks_count"`
	Stars       int       `json:"stargazers_count"`
	CreatedAt   time.Time `json:"created_at"`
}

func main() {
	var user, repo string

	fmt.Print("Enter user name: ")
	fmt.Scan(&user)

	fmt.Print("Enter name of the repository: ")
	fmt.Scan(&repo)

	if user == "" || repo == "" {
		fmt.Println("Error: name of user and repository is required")
		return
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", user, repo)
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while getting repository: ", err)
		return
	}

	switch response.StatusCode {
	case http.StatusOK:
		var repo Repository
		if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
			fmt.Println("Error while parsing json: ", err)
			return
		}

		fmt.Printf("\033[35m"+"Name: %s\n", repo.Name)
		fmt.Printf("Visibility: %s\n", repo.Visibility)
		fmt.Printf("Description: %s\n", repo.Description)
		fmt.Printf("Amount of forks: %d\n", repo.Forks)
		fmt.Printf("Amount of stars: %d\n", repo.Stars)
		fmt.Printf("Created at: %s\n", repo.CreatedAt.Format("2006-01-02 15:04:05 MST"))

	case http.StatusNotFound:
		fmt.Println("Repository not found")

	case http.StatusMovedPermanently:
		fmt.Println("Repository was moved")
		location := response.Header.Get("Location")
		if location != "" {
			fmt.Printf("New location: %s\n", location)
		}
	}

}
