# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=shortener
BINARY_UNIX=$(BINARY_NAME)_unix

# Docker parameters
DOCKER_COMPOSE=docker-compose
DOCKER_BUILD=docker build

.PHONY: all build clean test coverage deps lint run dev up down logs help

all: test build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) -v cmd/server/main.go

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./... -ignore=./docs
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy
	$(GOMOD) verify

# Lint code
lint:
	golangci-lint run

# Run the application locally
run: build
	./$(BINARY_NAME)

# Run the application in development mode
dev:
	$(GOCMD) run cmd/server/main.go

# Docker compose commands
up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_COMPOSE) down

logs:
	$(DOCKER_COMPOSE) logs -f

restart:
	$(DOCKER_COMPOSE) restart

rebuild:
	$(DOCKER_COMPOSE) down
	$(DOCKER_COMPOSE) up --build -d

# Generate swagger docs
docs:
	swag init -g cmd/server/main.go -o ./docs

# Run migrations (placeholder for future implementation)
migrate:
	@echo "Migrations will be handled by GORM AutoMigrate in this implementation"

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v cmd/server/main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build       - Build the application"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage"
	@echo "  deps        - Download dependencies"
	@echo "  lint        - Lint code"
	@echo "  run         - Build and run the application"
	@echo "  dev         - Run the application in development mode"
	@echo "  up          - Start all services with Docker Compose"
	@echo "  down        - Stop all services"
	@echo "  logs        - Show logs from all services"
	@echo "  restart     - Restart all services"
	@echo "  rebuild     - Rebuild and restart all services"
	@echo "  docs        - Generate swagger documentation"
	@echo "  migrate     - Run database migrations"
	@echo "  build-linux - Cross compile for Linux"
	@echo "  help        - Show this help message" 