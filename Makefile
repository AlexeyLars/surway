.PHONY: help build run test docker-up docker-down clean fmt lint

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build app
	cd backend && go build -o poll-service ./cmd/api

run: ## Run app locally
	cd backend && go run ./cmd/api/main.go

test: ## Run tests
	cd backend && go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Show tests coverage
	cd backend && go tool cover -html=coverage.out

docker-up: ## Run Docker Compose
	docker compose up -d

docker-down: ## Stop Docker Compose
	docker compose down

docker-logs: ## Show Docker Compose logs
	docker compose logs -f

docker-build: ## Build Docker image
	docker compose build

docker-build-backend: ## Build backend image only
	docker build -t poll-service:latest ./backend

docker-build-frontend: ## Build frontend image only
	docker build -t survey-frontend:latest ./frontend

clean: ## Clean build artefacts
	cd backend && rm -f poll-service coverage.out
	go clean

fmt: ## Format code
	cd backend && go fmt ./...

lint: ## Run linter
	cd backend && golangci-lint run

deps: ## Install dependencies
	cd backend && go mod download && go mod tidy

.DEFAULT_GOAL := help
