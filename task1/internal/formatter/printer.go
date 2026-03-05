package formatter

import (
	"fmt"
	repository "task1/internal/models"
)

func PrintInfo(repo *repository.Repository) {
	fmt.Println("Доступ успешно получен: ...\n")
	fmt.Printf("Имя Репозитория: %s\n", repo.Name)
	fmt.Printf("Описание: %s\n", repo.Description)
	fmt.Printf("Количество форков: %d\n", repo.Forks)
	fmt.Printf("Количество звезд: %d\n", repo.Stars)
	fmt.Printf("Дата создания: %s\n", repo.CreateAt[:10])
}
