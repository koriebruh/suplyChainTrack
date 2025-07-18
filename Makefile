BINARY_NAME=suplyChainTrack

# Definition Go
GO=go
GOCMD=$(GO)
GOTEST=$(GOCMD) test

# Flag Go
GO_FLAGS=-v

ifneq (,$(wildcard .env))
  include .env
  export
endif


help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Docker Development Commands:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Running binary
build:
	$(GOCMD) build $(GO_FLAGS) -o $(BINARY_NAME) ./cmd

# running main apps
run:
	@echo "Running apps..."
	$(GOCMD) run ./cmd

# running all unit test
test:
	@echo "Running tests..."
	$(GOTEST) -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -v -race -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Format kode Go
fmt:
	$(GOCMD) fmt ./...

swagger:
	@echo "Generating swagger or updating documentation..."
	swag init -g cmd/main.go

# Database migrations
create-db:
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -c "CREATE DATABASE $(DB_NAME);"

migrate-up:
	@echo "Running database migrations..."
	@which migrate > /dev/null || (echo "Please install golang-migrate" && exit 1)
	migrate -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" -path db/migrations up

migrate-down:
	@echo "Reverting database migrations..."
	migrate  -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)" -path db/migrations down

migrate-create:
	@echo "Creating new migration: $(name)"
	migrate create -ext sql -dir db/migrations $(name)

build:
	@echo "Building Docker image And Run another container..."
	docker compose -f docker/docker-compose.yml up -d --build

ps:
	@echo "Listing running Docker containers..."
	docker compose -f docker/docker-compose.yml ps

down:
	@echo "Stopping and removing Docker containers, networks, and volumes..."
	docker compose -f docker/docker-compose.yml down --volumes

logs:
	@echo "Showing logs for container: $(C)"
	docker compose -f docker/docker-compose.yml logs -f $(C)

