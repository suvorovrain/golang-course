package api

import (
	"encoding/json"
	"fmt"
	"log"
)

type Data struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	CreatedAt   string `json:"created_at"`
}

func parseJson(body []byte) Data {
	var data Data
	err := json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}
	return data
}

func OutputInfo(body []byte, args []string) {
	info := parseJson(body)

	flags := map[string]func(){
		"name":  func() { fmt.Println("Name:", info.Name) },
		"desc":  func() { fmt.Println("Description:", info.Description) },
		"star":  func() { fmt.Println("Stars:", info.Stars) },
		"forks": func() { fmt.Println("Forks:", info.Forks) },
		"date":  func() { fmt.Println("Created At:", info.CreatedAt[:10]) },
	}

	if len(args) == 0 {
		desc := info.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("Name: %s\nDescription: %s\nStars: %d\nForks: %d\nCreated At: %s\n",
			info.Name, desc, info.Stars, info.Forks, info.CreatedAt[:10])
		return
	}

	unknown := make([]string, 0)
	for i := range args {
		output, ok := flags[args[i]]
		if ok {
			output()
		} else {
			unknown = append(unknown, args[i])
		}
	}

	if len(unknown) > 0 {
		log.Fatalf("Unknown flags: (use 'help' for available flags)")
	}
}
