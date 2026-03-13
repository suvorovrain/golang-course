package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"
)

type Repository struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Stars        int       `json:"stargazers_count"`
	Forks        int       `json:"forks_count"`
	CreationDate time.Time `json:"created_at"`
}

func main() {
	owner := flag.String("owner", "", "имя владельца")
	repoName := flag.String("repo", "", "имя репозитория")

	flag.Parse()
	if *owner == "" || *repoName == "" {
		fmt.Println("Ошибка: переданы не все параметры")
		flag.Usage()
		return
	}

	repoUrl, err := url.Parse("http://api.github.com")
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
	}
	repoUrl.Path = path.Join("repos", *owner, *repoName)

	request, err := http.NewRequest("GET", repoUrl.String(), nil)
	if err != nil {
		fmt.Println("Ошибка формирования запроса:", err)
		return
	}

	request.Header.Set("User-Agent", "inforepo")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Ошибка запроса:", err)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Ошибка: статус", response.StatusCode)
		return
	}

	var repo Repository
	err = json.NewDecoder(response.Body).Decode(&repo)
	if err != nil {
		fmt.Println("Ошибка декодирования:", err)
		return
	}

	fmt.Printf("Имя репозитория: %s\n", repo.Name)
	fmt.Printf("Описание: %s\n", repo.Description)
	fmt.Printf("Количество звёзд: %d\n", repo.Stars)
	fmt.Printf("Количество форков: %d\n", repo.Forks)
	fmt.Printf("Дата создания: %s\n", repo.CreationDate.Format("02.01.2006"))
}
