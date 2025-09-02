.PHONY: help build run test clean docker-dev docker-prod docker-clean lint format deps

# Variables
BACKEND_DIR = backend
DOCKER_COMPOSE_DEV = docker-compose.yml
DOCKER_COMPOSE_PROD = docker-compose.prod.yml
GO_VERSION = 1.21

# Default target
help: ## Show this help message
	@echo "Available commands:"
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development Commands

deps: ## Install Go dependencies
	@echo "Installing Go dependencies..."
	cd $(BACKEND_DIR) && go mod download && go mod tidy

build: ## Build the backend application
	@echo "Building backend application..."
	cd $(BACKEND_DIR) && go build -o main .

run: ## Run the backend application locally
	@echo "Running backend application..."
	cd $(BACKEND_DIR) && go run main.go

test: ## Run tests
	@echo "Running tests..."
	cd $(BACKEND_DIR) && go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	cd $(BACKEND_DIR) && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: backend/coverage.html"

lint: ## Run linter
	@echo "Running linter..."
	cd $(BACKEND_DIR) && golangci-lint run

format: ## Format Go code
	@echo "Formatting Go code..."
	cd $(BACKEND_DIR) && go fmt ./...

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	cd $(BACKEND_DIR) && rm -f main coverage.out coverage.html

##@ Docker Commands

docker-dev: ## Start development environment with Docker
	@echo "Starting development environment..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) up --build -d

docker-dev-logs: ## View development logs
	docker-compose -f $(DOCKER_COMPOSE_DEV) logs -f

docker-dev-stop: ## Stop development environment
	@echo "Stopping development environment..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) down

docker-prod: ## Start production environment with Docker
	@echo "Starting production environment..."
	@if [ ! -f .env.prod ]; then \
		echo "Error: .env.prod file not found. Copy from env.prod.example and configure."; \
		exit 1; \
	fi
	docker-compose -f $(DOCKER_COMPOSE_PROD) --env-file .env.prod up -d

docker-prod-logs: ## View production logs
	docker-compose -f $(DOCKER_COMPOSE_PROD) logs -f

docker-prod-stop: ## Stop production environment
	@echo "Stopping production environment..."
	docker-compose -f $(DOCKER_COMPOSE_PROD) down

docker-clean: ## Clean Docker resources
	@echo "Cleaning Docker resources..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) down -v --remove-orphans
	docker-compose -f $(DOCKER_COMPOSE_PROD) down -v --remove-orphans
	docker system prune -f

##@ Database Commands

db-migrate: ## Run database migration
	@echo "Running database migration..."
	cd $(BACKEND_DIR) && go run app/pkg/database/cmd/migrate.go

db-reset: ## Reset database (WARNING: This will delete all data)
	@echo "Resetting database..."
	docker-compose -f $(DOCKER_COMPOSE_DEV) down -v
	docker-compose -f $(DOCKER_COMPOSE_DEV) up -d postgres
	@echo "Waiting for database to be ready..."
	sleep 10
	$(MAKE) db-migrate

##@ Documentation Commands

docs-generate: ## Generate Swagger documentation
	@echo "Generating Swagger documentation..."
	cd $(BACKEND_DIR) && swag init -g main.go --output docs

docs-serve: ## Serve documentation locally
	@echo "Documentation available at: http://localhost:3000/swagger/"
	$(MAKE) run

##@ Git & Release Commands

tag: ## Create and push a new tag (usage: make tag VERSION=v1.0.0)
ifndef VERSION
	@echo "Error: VERSION is required. Usage: make tag VERSION=v1.0.0"
	@exit 1
endif
	@echo "Creating tag $(VERSION)..."
	git tag $(VERSION)
	git push origin $(VERSION)
	@echo "Tag $(VERSION) created and pushed successfully!"

release: ## Create a release (build, test, tag)
ifndef VERSION
	@echo "Error: VERSION is required. Usage: make release VERSION=v1.0.0"
	@exit 1
endif
	@echo "Creating release $(VERSION)..."
	$(MAKE) test
	$(MAKE) lint
	$(MAKE) build
	$(MAKE) tag VERSION=$(VERSION)
	@echo "Release $(VERSION) created successfully!"

##@ Health & Monitoring Commands

health: ## Check application health
	@echo "Checking application health..."
	@curl -s http://localhost:3000/health | jq . || echo "Application not running or jq not installed"

status: ## Show Docker containers status
	@echo "Docker containers status:"
	@docker-compose -f $(DOCKER_COMPOSE_DEV) ps

logs: ## Show application logs
	docker-compose -f $(DOCKER_COMPOSE_DEV) logs -f backend

##@ Setup Commands

setup-dev: ## Complete development setup
	@echo "Setting up development environment..."
	$(MAKE) deps
	@if [ ! -f $(BACKEND_DIR)/.env ]; then \
		cp $(BACKEND_DIR)/.env.example $(BACKEND_DIR)/.env; \
		echo "Created .env file from .env.example"; \
		echo "Please edit $(BACKEND_DIR)/.env with your configuration"; \
	fi
	$(MAKE) docker-dev
	@echo "Development environment setup complete!"
	@echo "API available at: http://localhost:3000"
	@echo "Swagger docs at: http://localhost:3000/swagger/"

setup-prod: ## Setup production environment
	@echo "Setting up production environment..."
	@if [ ! -f .env.prod ]; then \
		cp env.prod.example .env.prod; \
		echo "Created .env.prod file from env.prod.example"; \
		echo "Please edit .env.prod with your production configuration"; \
		exit 1; \
	fi
	$(MAKE) docker-prod
	@echo "Production environment setup complete!"

##@ Utility Commands

install-tools: ## Install development tools
	@echo "Installing development tools..."
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed successfully!"

check-deps: ## Check for outdated dependencies
	@echo "Checking for outdated dependencies..."
	cd $(BACKEND_DIR) && go list -u -m all

update-deps: ## Update dependencies
	@echo "Updating dependencies..."
	cd $(BACKEND_DIR) && go get -u ./... && go mod tidy
