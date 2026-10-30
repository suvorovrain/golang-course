APP_NAME := demo-app
APP_PORT := 8081
PROM_PORT := 9090
COMPOSE := docker compose

.PHONY: help tidy build run up down logs ping ping-fail metrics smoke urls

help:
	@echo "Targets:"
	@echo "  make tidy      - go mod tidy"
	@echo "  make build     - build app binary"
	@echo "  make run       - run app locally on :$(APP_PORT)"
	@echo "  make up        - start app + prometheus via docker compose"
	@echo "  make down      - stop docker compose"
	@echo "  make logs      - show compose logs"
	@echo "  make ping      - successful ping request"
	@echo "  make ping-fail - failed ping request (for error metrics)"
	@echo "  make metrics   - show a slice of /metrics"
	@echo "  make smoke     - run basic checks"
	@echo "  make urls      - print local service URLs"

tidy:
	go mod tidy

build:
	go build -o bin/$(APP_NAME) ./cmd/app

run:
	go run ./cmd/app

up:
	$(COMPOSE) up --build -d

down:
	$(COMPOSE) down

logs:
	$(COMPOSE) logs -f --tail=200

ping:
	curl -s "http://localhost:$(APP_PORT)/api/v1/ping?name=jarvis"

ping-fail:
	curl -s -o /dev/null -w "status=%{http_code}\n" "http://localhost:$(APP_PORT)/api/v1/ping?fail=true"

metrics:
	curl -s "http://localhost:$(APP_PORT)/metrics" | rg '^demo_' | head -n 40

smoke:
	curl -s "http://localhost:$(APP_PORT)/healthz"
	curl -s "http://localhost:$(APP_PORT)/api/v1/ping?name=smoke"
	curl -s "http://localhost:$(APP_PORT)/metrics" | rg '^demo_ping_requests_total' -n

urls:
	@echo "API        http://localhost:$(APP_PORT)"
	@echo "Swagger    http://localhost:$(APP_PORT)/swagger/index.html"
	@echo "Prometheus http://localhost:$(PROM_PORT)"
	@echo "Grafana    http://localhost:3000 (admin/admin)"
