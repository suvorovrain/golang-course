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
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания запроса: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("неожиданный статус: %s", resp.Status)
	}

	var res RepoInfo
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("неожиданный статус: %s", resp.Status)
	}

	return &res, nil
}

// Подумать над добавлением токена для запросов
func main() {
	if len(os.Args) != 3 {
		fmt.Println("Некорректный ввод")
		return
	}
	res, err := request(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Println(":((")
		return
	}
	fmt.Printf("Репозиторий: %s\nОписание: %s\nЗвёзд: %d\nФорков: %d\nСоздан: %s\n",
		res.Name, res.Description, res.Stars, res.Forks, res.Date)

	fmt.Println(":)")
}
