.PHONY: all run build fmt lint test down open docs

run:
	docker compose up --build -d

build:
	go build ./...

fmt:
	go fmt ./...

down:
	docker compose down --rmi local -v

docs:
	go run github.com/swaggo/swag/cmd/swag@v1.8.12 init -g cmd/main.go -d gateway --parseInternal --output gateway/docs

all: run
