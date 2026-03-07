package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	NumberOfStars int    `json:"stargazers_count"`
	NumberOfForks int    `json:"forks_count"`
	CreationDate  string `json:"created_at"`
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the owner's name: ")
	inputName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Reading error:", err)
		return
	}

	name := strings.TrimSpace(inputName)

	fmt.Println("Enter the repository's name: ")
	inputRepo, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Reading error:", err)
		return
	}

	repo := strings.TrimSpace(inputRepo)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", name, repo)
	checkCLI((url))
}

func checkCLI(link string) {
	request, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Fatalf("Request creation error: %s", err)
	}

	request.Header.Set("User-Agent", "golang-course")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Fatalf("Request execution error: %s", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		log.Fatalf("API returns the error %d: %s", response.StatusCode, string(body))
	}

	var repo Repository
	if err := json.NewDecoder(response.Body).Decode(&repo); err != nil {
		log.Fatalf("JSON parsing error: %s", err)
	}

	fmt.Printf("Repository: %s\n", repo.Name)
	fmt.Printf("Description: %s\n", repo.Description)
	fmt.Printf("Number of stars: %d\n", repo.NumberOfStars)
	fmt.Printf("Number of forks: %d\n", repo.NumberOfForks)

	if t, err := time.Parse(time.RFC3339, repo.CreationDate); err == nil {
		fmt.Printf("Created At: %s\n", t.Format("January 2, 2006 at 15:04"))
	} else {
		fmt.Printf("Created At: %s\n", repo.CreationDate)
	}
}
