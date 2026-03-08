package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type githubRepoUrl struct {
	owner string
	name  string
}

func isValidUrl(str string) bool {
	match, _ := regexp.MatchString("https://github.com/[A-Za-z-]+/[A-Za-z-]+", str)
	return match
}

func parseUrl(url string) (githubRepoUrl, error) {
	if !isValidUrl(url) {
		return githubRepoUrl{}, errors.New("Url with github repo expected")
	}
	split := strings.Split(url, "/")
	// "https:", "", "github.com", "OWNER", "NAME"
	return githubRepoUrl{owner: split[3], name: split[4]}, nil
}

func usage() {
	usage_message :=
		"usage: %s [url_of_github_repo]\n" +
			"The flags are:\n" +
			"\t--help - print this message\n"

	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, usage_message, os.Args[0])
	os.Exit(2)
}

func getJson(url string) (map[string]interface{}, error) {
	repo, err := parseUrl(url)
	if err != nil {
		return nil, err
	}
	res, err := http.Get("https://api.github.com/repos/" + repo.owner + "/" + repo.name)
	if err != nil {
		return nil, errors.New("Error making http request")
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Bad request")
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("Can't read response body")
	}
	var dat map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &dat); err != nil {
		panic(err)
	}

	return dat, nil
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Input file is missing")
		os.Exit(1)
	}
	args := flag.Args()
	url := args[0]

	json, err := getJson(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Println(json["name"])
	fmt.Println(json["description"])
	fmt.Println(json["stargazers_count"])
	fmt.Println(json["fork_count"])
	fmt.Println(json["created_at"])

}
