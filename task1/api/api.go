package api

import (
	"io"
	"log"
	"net/http"
)

func GetInfo(url string, args []string) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode == 404 {
		log.Fatal("repository not found")
	}
	if response.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", response.StatusCode, body)
	}

	OutputInfo(body, args)
}
