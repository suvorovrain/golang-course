package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"task1/internal/models"
)

func FetchRepo(owner, repo string) (*models.Repository, error) {
	//формируем ссылку для запроса
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)
	//запрос
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Ошибка запроса: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка доступа к репозиторию, %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения из ответа %w", err)
	}
	var repository models.Repository
	jsonRepo := json.Unmarshal(body, &repository)
	if jsonRepo != nil {
		return nil, fmt.Errorf("Ошибка парсинга %w", jsonRepo)
	}
	return &repository, nil

}
