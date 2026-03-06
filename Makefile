.PHONY: help build run stop clean test docker-build docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the Go application
	go build -o bin/api ./cmd/api

run: ## Run the application locally
	go run ./cmd/api

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf data/

docker-build: ## Build Docker images
	docker-compose build

docker-up: ## Start services with Docker Compose (local embeddings)
	docker-compose up -d

docker-up-openai: ## Start services with OpenAI embeddings
	docker-compose -f docker-compose.yml -f docker-compose.openai.yml up -d

docker-down: ## Stop Docker services
	docker-compose down

docker-logs: ## Show Docker logs
	docker-compose logs -f

deps: ## Download Go dependencies
	go mod download
	go mod tidy

setup: ## Setup environment
	cp .env.example .env
	mkdir -p data

.DEFAULT_GOAL := help
