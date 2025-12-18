# ============================================================================
# SaaS Photo Listing Platform - Makefile
# ============================================================================
# This Makefile provides a consistent interface for building, testing, and
# deploying the photo cloud application.
#
# Usage: make <command>
# Example: make run
# ============================================================================

# ============================================================================
# Configuration & Variables
# ============================================================================

# Project metadata
PROJECT_NAME := saas-photo-listing-platform
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Go settings
GO := go
GO_MODULE := $(shell cat go.mod | head -1 | awk '{print $$2}')
BINARY_NAME := photo-cloud-api
BUILD_DIR := ./bin
MAIN_PACKAGE := ./cmd/api

# Database
DB_NAME := saas-photo-listing-platform
DB_USER := Enock
DB_PASSWORD := enockpostgresql02
DB_HOST := localhost
DB_PORT := 5432
DATABASE_URL ?= postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Docker
DOCKER_COMPOSE_FILE := deployments/docker-compose.yml
DOCKER_REGISTRY ?= ghcr.io  # Change to your registry
DOCKER_IMAGE := $(DOCKER_REGISTRY)/your-username/$(PROJECT_NAME):$(VERSION)

# Development tools (versions can be pinned)
MIGRATE_VERSION := v4.17.0
SQLC_VERSION := v1.24.0
AIR_VERSION := v1.51.0
GOLANGCI_LINT_VERSION := v1.56.2

# Colors for pretty output
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
BOLD := \033[1m
RESET := \033[0m

# ============================================================================
# PHONY Targets
# ============================================================================
# Phony targets are always executed, even if a file with the same name exists
.PHONY: help \
        install-tools install \
        run dev \
        build build-linux build-darwin build-windows \
        test test-verbose test-coverage \
        lint lint-fix \
        fmt \
        clean \
        docker-up docker-down docker-logs docker-clean \
        migrate-create migrate-up migrate-down migrate-status migrate-reset \
        sqlc sqlc-validate \
        db-setup db-reset db-seed \
        setup check-env \
        pre-commit \
        swagger \
        release

# ============================================================================
# Help
# ============================================================================
## Display this help message
help:
	@echo "$(BOLD)$(CYAN)$(PROJECT_NAME) - Available Commands$(RESET)"
	@echo ""
	@echo "$(BOLD)Development:$(RESET)"
	@echo "  $(GREEN)make run$(RESET)          - Start the application with hot reload"
	@echo "  $(GREEN)make dev$(RESET)          - Start the application without hot reload"
	@echo "  $(GREEN)make setup$(RESET)        - Full development environment setup"
	@echo ""
	@echo "$(BOLD)Database:$(RESET)"
	@echo "  $(GREEN)make db-setup$(RESET)     - Setup database (migrations + sqlc)"
	@echo "  $(GREEN)make migrate-create$(RESET) - Create a new migration"
	@echo "  $(GREEN)make migrate-up$(RESET)   - Apply pending migrations"
	@echo "  $(GREEN)make sqlc$(RESET)         - Generate SQLC code"
	@echo "  $(GREEN)make db-reset$(RESET)     - Reset database (DANGER: destroys data!)"
	@echo ""
	@echo "$(BOLD)Testing & Quality:$(RESET)"
	@echo "  $(GREEN)make test$(RESET)         - Run all tests"
	@echo "  $(GREEN)make test-coverage$(RESET) - Run tests with coverage report"
	@echo "  $(GREEN)make lint$(RESET)         - Run linter"
	@echo "  $(GREEN)make lint-fix$(RESET)     - Run linter and auto-fix issues"
	@echo "  $(GREEN)make fmt$(RESET)          - Format Go code"
	@echo "  $(GREEN)make pre-commit$(RESET)   - Run all checks before commit"
	@echo ""
	@echo "$(BOLD)Build & Deploy:$(RESET)"
	@echo "  $(GREEN)make build$(RESET)        - Build application for current OS"
	@echo "  $(GREEN)make build-all$(RESET)    - Build for Linux, macOS, and Windows"
	@echo "  $(GREEN)make docker-up$(RESET)    - Start all Docker services"
	@echo "  $(GREEN)make docker-down$(RESET)  - Stop all Docker services"
	@echo ""
	@echo "$(BOLD)Maintenance:$(RESET)"
	@echo "  $(GREEN)make clean$(RESET)        - Clean build artifacts"
	@echo "  $(GREEN)make install-tools$(RESET) - Install development tools"
	@echo "  $(GREEN)make check-env$(RESET)    - Check environment setup"
	@echo ""
	@echo "$(BOLD)Information:$(RESET)"
	@echo "  Version: $(VERSION)"
	@echo "  Go Module: $(GO_MODULE)"
	@echo ""

# ============================================================================
# Development Tools Installation
# ============================================================================
## Install all required development tools
install-tools:
	@echo "$(CYAN)Installing development tools...$(RESET)"
	
	@echo "$(YELLOW)Installing migrate $(MIGRATE_VERSION)...$(RESET)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$(MIGRATE_VERSION); \
	else \
		echo "✓ migrate already installed"; \
	fi
	
	@echo "$(YELLOW)Installing sqlc $(SQLC_VERSION)...$(RESET)"
	@if ! command -v sqlc >/dev/null 2>&1; then \
		$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION); \
	else \
		echo "✓ sqlc already installed"; \
	fi
	
	@echo "$(YELLOW)Installing air $(AIR_VERSION)...$(RESET)"
	@if ! command -v air >/dev/null 2>&1; then \
		$(GO) install github.com/air-verse/air@$(AIR_VERSION); \
	else \
		echo "✓ air already installed"; \
	fi
	
	@echo "$(YELLOW)Installing golangci-lint $(GOLANGCI_LINT_VERSION)...$(RESET)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION); \
	else \
		echo "✓ golangci-lint already installed"; \
	fi
	
	@echo "$(YELLOW)Installing swag (for OpenAPI docs)...$(RESET)"
	@if ! command -v swag >/dev/null 2>&1; then \
		$(GO) install github.com/swaggo/swag/cmd/swag@latest; \
	else \
		echo "✓ swag already installed"; \
	fi
	
	@echo "$(GREEN)✓ All tools installed successfully$(RESET)"

## Install Go dependencies
install:
	@echo "$(CYAN)Installing Go dependencies...$(RESET)"
	$(GO) mod download
	$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies installed$(RESET)"

# ============================================================================
# Development
# ============================================================================
## Run the application with hot reload (air)
run:
	@echo "$(CYAN)Starting application with hot reload...$(RESET)"
	@if [ ! -f .env ]; then \
		echo "$(YELLOW)⚠ Warning: .env file not found. Copying from .env.example...$(RESET)"; \
		cp -n .env.example .env || true; \
	fi
	@if ! command -v air >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'air' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	air

## Run the application without hot reload
dev:
	@echo "$(CYAN)Starting application...$(RESET)"
	$(GO) run $(MAIN_PACKAGE)

# ============================================================================
# Building
# ============================================================================
## Build application for current OS/architecture
build:
	@echo "$(CYAN)Building $(BINARY_NAME) for $(shell go env GOOS)/$(shell go env GOARCH)...$(RESET)"
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 $(GO) build \
		-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) \
		$(MAIN_PACKAGE)
	@echo "$(GREEN)✓ Built: $(BUILD_DIR)/$(BINARY_NAME)$(RESET)"

## Build for Linux
build-linux:
	@echo "$(CYAN)Building for Linux...$(RESET)"
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 \
		$(MAIN_PACKAGE)

## Build for macOS
build-darwin:
	@echo "$(CYAN)Building for macOS...$(RESET)"
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GO) build \
		-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 \
		$(MAIN_PACKAGE)

## Build for Windows
build-windows:
	@echo "$(CYAN)Building for Windows...$(RESET)"
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 $(GO) build \
		-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe \
		$(MAIN_PACKAGE)

## Build for all platforms
build-all: build-linux build-darwin build-windows
	@echo "$(GREEN)✓ Built for all platforms$(RESET)"

# ============================================================================
# Testing
# ============================================================================
## Run all tests
test:
	@echo "$(CYAN)Running tests...$(RESET)"
	$(GO) test ./... -short

## Run all tests with verbose output
test-verbose:
	@echo "$(CYAN)Running tests with verbose output...$(RESET)"
	$(GO) test ./... -v

## Run tests with coverage report
test-coverage:
	@echo "$(CYAN)Running tests with coverage...$(RESET)"
	$(GO) test ./... -coverprofile=coverage.out -covermode=atomic
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(RESET)"

# ============================================================================
# Code Quality
# ============================================================================
## Run linter
lint:
	@echo "$(CYAN)Running golangci-lint...$(RESET)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'golangci-lint' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	golangci-lint run ./...

## Run linter and auto-fix issues
lint-fix:
	@echo "$(CYAN)Running golangci-lint with auto-fix...$(RESET)"
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'golangci-lint' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	golangci-lint run ./... --fix

## Format Go code
fmt:
	@echo "$(CYAN)Formatting Go code...$(RESET)"
	$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(RESET)"

# ============================================================================
# Database Operations
# ============================================================================
## Create a new migration file
migration-create:
	@echo "$(CYAN)Creating new migration...$(RESET)"
	@if [ -z "$(name)" ]; then \
		read -p "Enter migration name: " name; \
	fi
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'migrate' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	migrate create -ext sql -dir internal/database/migrations -seq $(name)

## Apply pending migrations
migrate-up:
	@echo "$(CYAN)Applying migrations...$(RESET)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'migrate' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" up

## Rollback last migration
migrate-down:
	@echo "$(CYAN)Rolling back last migration...$(RESET)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'migrate' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" down 1

## Show migration status
migrate-status:
	@echo "$(CYAN)Migration status...$(RESET)"
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'migrate' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" version

## Reset database (DANGER: destroys all data!)
migrate-reset:
	@echo "$(RED)⚠ WARNING: This will DESTROY ALL DATA in the database!$(RESET)"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" != "yes" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	@if ! command -v migrate >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'migrate' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" down -f
	migrate -path internal/database/migrations -database "$(DATABASE_URL)" up
	@echo "$(GREEN)✓ Database reset complete$(RESET)"

## Generate SQLC code
sqlc:
	@echo "$(CYAN)Generating SQLC code...$(RESET)"
	@if ! command -v sqlc >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'sqlc' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	sqlc generate
	@echo "$(GREEN)✓ SQLC code generated$(RESET)"

## Validate SQLC configuration and queries
sqlc-validate:
	@echo "$(CYAN)Validating SQLC configuration...$(RESET)"
	@if ! command -v sqlc >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'sqlc' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	sqlc compile
	@echo "$(GREEN)✓ SQLC configuration is valid$(RESET)"

## Setup database (migrations + sqlc)
db-setup: migrate-up sqlc
	@echo "$(GREEN)✓ Database setup complete$(RESET)"

## Reset database and re-run setup
db-reset: migrate-reset sqlc
	@echo "$(GREEN)✓ Database reset and setup complete$(RESET)"

## Seed database with sample data
db-seed:
	@echo "$(CYAN)Seeding database...$(RESET)"
	@if [ -f scripts/seed.go ]; then \
		$(GO) run scripts/seed.go; \
	else \
		echo "$(YELLOW)⚠ No seed script found at scripts/seed.go$(RESET)"; \
	fi

# ============================================================================
# Docker Operations
# ============================================================================
## Start all Docker services
docker-up:
	@echo "$(CYAN)Starting Docker services...$(RESET)"
	@if [ ! -f $(DOCKER_COMPOSE_FILE) ]; then \
		echo "$(RED)✗ Error: Docker compose file not found at $(DOCKER_COMPOSE_FILE)$(RESET)"; \
		exit 1; \
	fi
	docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@echo "$(GREEN)✓ Docker services started$(RESET)"
	@echo "  PostgreSQL: localhost:5432"
	@echo "  Redis: localhost:6379"
	@echo "  MinIO: http://localhost:9000 (Console: http://localhost:9001)"
	@echo "  MinIO credentials: minioadmin / minioadmin"

## Stop all Docker services
docker-down:
	@echo "$(CYAN)Stopping Docker services...$(RESET)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@echo "$(GREEN)✓ Docker services stopped$(RESET)"

## Show Docker service logs
docker-logs:
	@echo "$(CYAN)Showing Docker logs...$(RESET)"
	docker-compose -f $(DOCKER_COMPOSE_FILE) logs -f

## Clean Docker resources (containers, volumes, networks)
docker-clean:
	@echo "$(RED)⚠ WARNING: This will remove ALL Docker resources!$(RESET)"
	@read -p "Are you sure? (yes/no): " confirm; \
	if [ "$$confirm" != "yes" ]; then \
		echo "Aborted."; \
		exit 1; \
	fi
	docker-compose -f $(DOCKER_COMPOSE_FILE) down -v --rmi all --remove-orphans
	@echo "$(GREEN)✓ Docker resources cleaned$(RESET)"

# ============================================================================
# Setup & Environment
# ============================================================================
## Complete development environment setup
setup: install-tools install docker-up db-setup
	@echo ""
	@echo "$(GREEN)========================================$(RESET)"
	@echo "$(GREEN)✅ Development environment ready!$(RESET)"
	@echo "$(GREEN)========================================$(RESET)"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Edit .env file with your configuration"
	@echo "  2. Run $(CYAN)make run$(RESET) to start the app"
	@echo "  3. Visit $(CYAN)http://localhost:8080$(RESET)"
	@echo ""

## Check environment setup
check-env:
	@echo "$(CYAN)Checking environment...$(RESET)"
	
	@echo -n "Go version: "
	@$(GO) version || echo "$(RED)✗ Go not installed$(RESET)"
	
	@echo -n "Docker: "
	@docker --version 2>/dev/null && echo "✓" || echo "$(RED)✗ Docker not installed$(RESET)"
	
	@echo -n "Docker Compose: "
	@docker-compose --version 2>/dev/null && echo "✓" || echo "$(RED)✗ Docker Compose not installed$(RESET)"
	
	@echo -n ".env file: "
	@if [ -f .env ]; then echo "✓"; else echo "$(YELLOW)⚠ Not found (copy .env.example)$(RESET)"; fi
	
	@echo -n "Database connection: "
	@if command -v psql >/dev/null 2>&1; then \
		if PGPASSWORD=$(DB_PASSWORD) psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "SELECT 1" >/dev/null 2>&1; then \
			echo "✓"; \
		else \
			echo "$(YELLOW)⚠ Cannot connect$(RESET)"; \
		fi \
	else \
		echo "$(YELLOW)⚠ psql not installed$(RESET)"; \
	fi

# ============================================================================
# Pre-commit Checks
# ============================================================================
## Run all checks before committing
pre-commit: fmt lint test build
	@echo "$(GREEN)✅ All pre-commit checks passed!$(RESET)"

# ============================================================================
# Documentation
# ============================================================================
## Generate OpenAPI/Swagger documentation
swagger:
	@echo "$(CYAN)Generating Swagger documentation...$(RESET)"
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "$(RED)✗ Error: 'swag' is not installed. Run 'make install-tools' first$(RESET)"; \
		exit 1; \
	fi
	swag init -g cmd/api/main.go -o docs/swagger
	@echo "$(GREEN)✓ Swagger docs generated at docs/swagger$(RESET)"

# ============================================================================
# Cleanup
# ============================================================================
## Clean build artifacts
clean:
	@echo "$(CYAN)Cleaning build artifacts...$(RESET)"
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf tmp/
	$(GO) clean
	@echo "$(GREEN)✓ Cleaned$(RESET)"

# ============================================================================
# Release
# ============================================================================
## Create a release build
release: test lint build-all
	@echo "$(GREEN)✅ Release build complete$(RESET)"
	@echo "Version: $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@ls -lh $(BUILD_DIR)/*

# ============================================================================
# Utility Functions
# ============================================================================
# These are internal functions, not meant to be called directly

# Check if a command exists
_check_command = @if ! command -v $(1) >/dev/null 2>&1; then \
	echo "$(RED)✗ Error: '$(1)' is not installed$(RESET)"; \
	exit 1; \
	fi

# Wait for a service to be ready
_wait_for_service = @echo "Waiting for $(1) at $(2)..." && \
	for i in $$(seq 1 30); do \
		if nc -z $(subst :, ,$(2)) >/dev/null 2>&1; then \
			echo "✓ $(1) ready"; \
			break; \
		fi; \
		echo "  Attempt $$i/30"; \
		sleep 2; \
		if [ $$i -eq 30 ]; then \
			echo "$(RED)✗ $(1) failed to start$(RESET)"; \
			exit 1; \
		fi; \
	done