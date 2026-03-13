package main

import (
	"fmt"
	"githubcli/api"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Error: no arguments provided")
		fmt.Println("Run 'githubcli help' for usage information")
		os.Exit(1)
	}

	command := os.Args[1]
	if command == "help" {
		printHelp()
		return
	}

	owner := os.Args[1]
	repo := os.Args[2]
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	flags := []string{}
	for i := 3; i < len(os.Args); i++ {
		flags = append(flags, strings.TrimPrefix(os.Args[i], "-"))
	}

	api.GetInfo(url, flags)
}

func printHelp() {
	fmt.Println("GitHub CLI - Usage")
	fmt.Println("Commands:")
	fmt.Println("  githubcli <owner> <repo>   Show repository info")
	fmt.Println("  githubcli help     Show this help message")

	fmt.Println("Flags for <repo> command:")
	fmt.Println("  -name      Display repository name")
	fmt.Println("  -desc      Display repository description")
	fmt.Println("  -star      Display number of stars")
	fmt.Println("  -forks     Display number of forks")
	fmt.Println("  -date      Display creation date")

	fmt.Println("Example:")
	fmt.Println("  githubcli golang go -name -star")
}
