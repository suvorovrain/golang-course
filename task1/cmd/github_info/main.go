package main

import (
	"flag"
	"fmt"
	"os"
	"task1/internal/formatter"
	"task1/internal/github"
)

func main() {
	//обработка флагов
	owner := flag.String("owner", "", "owner of repository")
	repo := flag.String("repo", "", "name of repository")
	flag.Parse()

	if *owner == "" || *repo == "" {
		fmt.Fprintf(os.Stderr, "Ошибка: Введите go run main.go --owner \"your_name\", --repo \"your_repository_name\"")
		os.Exit(1)
	}
	//вызываем функцию для запроса из github package
	repository, err := github.FetchRepo(*owner, *repo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	//если все успешно прочитано, то выводим с помощью кастомного вывода из formatter
	formatter.PrintInfo(repository)

}
