.PHONY: help dev build test clean migrate-up migrate-down migrate-status sqlc-generate install-backend install-frontend install

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

dev: ## Start development environment (Docker Compose)
	docker-compose up

dev-build: ## Build and start development environment
	docker-compose up --build

dev-down: ## Stop development environment
	docker-compose down

dev-logs: ## Tail logs from all services
	docker-compose logs -f

dev-logs-api: ## Tail logs from API service
	docker-compose logs -f api

dev-logs-web: ## Tail logs from web service
	docker-compose logs -f web

##@ Backend

backend-run: ## Run backend locally (requires .env.dev)
	cd backend && go run cmd/server/main.go

backend-test: ## Run backend tests
	cd backend && go test -v -cover ./...

backend-test-coverage: ## Run backend tests with coverage report
	cd backend && go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

backend-lint: ## Lint backend code
	cd backend && golangci-lint run

backend-tidy: ## Tidy Go modules
	cd backend && go mod tidy

sqlc-generate: ## Generate sqlc code
	cd backend && sqlc generate

##@ Frontend

frontend-dev: ## Run frontend development server
	cd frontend && npm run dev

frontend-build: ## Build frontend for production
	cd frontend && npm run build

frontend-test: ## Run frontend tests
	cd frontend && npm test

frontend-test-ui: ## Run frontend tests with UI
	cd frontend && npm run test:ui

frontend-test-coverage: ## Run frontend tests with coverage
	cd frontend && npm run test:coverage

frontend-lint: ## Lint frontend code
	cd frontend && npm run lint

##@ Database

migrate-up: ## Apply all database migrations
	cd backend && goose -dir internal/store/migrations postgres "$$DATABASE_URL" up

migrate-down: ## Rollback last migration
	cd backend && goose -dir internal/store/migrations postgres "$$DATABASE_URL" down

migrate-status: ## Check migration status
	cd backend && goose -dir internal/store/migrations postgres "$$DATABASE_URL" status

migrate-create: ## Create a new migration (usage: make migrate-create NAME=create_table_name)
	cd backend && goose -dir internal/store/migrations create $(NAME) sql

db-reset: ## Reset database (WARNING: drops all data)
	cd backend && goose -dir internal/store/migrations postgres "$$DATABASE_URL" reset

##@ Installation

install-backend: ## Install backend dependencies and tools
	cd backend && go mod download
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-frontend: ## Install frontend dependencies
	cd frontend && npm install

install: install-backend install-frontend ## Install all dependencies

##@ Utilities

clean: ## Clean build artifacts and caches
	rm -rf backend/bin backend/dist backend/coverage.out
	rm -rf frontend/dist frontend/coverage frontend/node_modules/.vite

format-backend: ## Format backend code
	cd backend && go fmt ./...

format-frontend: ## Format frontend code
	cd frontend && npx prettier --write "src/**/*.{ts,tsx,css}"

format: format-backend format-frontend ## Format all code

##@ Docker

docker-build-api: ## Build API Docker image
	docker build -t pinecone-api:latest -f backend/Dockerfile backend/

docker-build-web: ## Build web Docker image (production)
	docker build -t pinecone-web:latest -f frontend/Dockerfile frontend/

docker-build: docker-build-api docker-build-web ## Build all Docker images

docker-push: ## Push Docker images to registry (requires REGISTRY env var)
	docker tag pinecone-api:latest $(REGISTRY)/pinecone-api:latest
	docker tag pinecone-web:latest $(REGISTRY)/pinecone-web:latest
	docker push $(REGISTRY)/pinecone-api:latest
	docker push $(REGISTRY)/pinecone-web:latest
