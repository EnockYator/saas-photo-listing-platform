# =========================
# Variables
# =========================
DOCKER_COMPOSE=docker-compose -f docker-compose.yml
BACKEND_DIR=./backend
FRONTEND_DIR=./react-frontend
DOCS_DIR=$(BACKEND_DIR)/internal/docs

# =========================
# Docker Commands
# =========================

## Build all Docker services
.PHONY: build
build:
	$(DOCKER_COMPOSE) build

## Start all services in detached mode
.PHONY: up
up:
	$(DOCKER_COMPOSE) up -d

## Stop all services
.PHONY: down
down:
	$(DOCKER_COMPOSE) down

## Restart all services
.PHONY: restart
restart: down up

## Follow logs for all services
.PHONY: logs
logs:
	$(DOCKER_COMPOSE) logs -f

# =========================
# Backend Commands
# =========================

## Run backend API locally (without Docker)
.PHONY: run-backend
run-backend:
	cd $(BACKEND_DIR) && go run ./cmd/api/main.go

## Run backend worker locally (without Docker)
.PHONY: run-worker
run-worker:
	cd $(BACKEND_DIR) && go run ./cmd/worker/main.go

## Run database migrations
.PHONY: migrate
migrate:
	cd $(BACKEND_DIR) && go run ./cmd/tools/migrate.go

## Generate SQLC type-safe queries
.PHONY: sqlc
sqlc:
	cd $(BACKEND_DIR) && sqlc generate

## Generate Swagger / OpenAPI docs
.PHONY: swagger
swagger:
	cd $(BACKEND_DIR) && swag init -g cmd/api/main.go -o $(DOCS_DIR)

## Run backend tests
.PHONY: test-backend
test-backend:
	cd $(BACKEND_DIR) && go test ./... -v

# =========================
# Frontend Commands
# =========================

## Install frontend dependencies
.PHONY: install-frontend
install-frontend:
	cd $(FRONTEND_DIR) && npm ci

## Run frontend in dev mode
.PHONY: run-frontend
run-frontend:
	cd $(FRONTEND_DIR) && npm run dev

## Build frontend for production
.PHONY: build-frontend
build-frontend:
	cd $(FRONTEND_DIR) && npm run build

## Test frontend
.PHONY: test-frontend
test-frontend:
	cd $(FRONTEND_DIR) && npm test

# =========================
# Tools / Utilities
# =========================

## Lint backend code
.PHONY: lint-backend
lint-backend:
	cd $(BACKEND_DIR) && golangci-lint run

## Lint frontend code
.PHONY: lint-frontend
lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

# =========================
# Clean / Reset
# =========================

## Clean backend and frontend build artifacts
.PHONY: clean
clean:
	cd $(BACKEND_DIR) && go clean
	cd $(FRONTEND_DIR) && rm -rf node_modules dist

## Reset database (use with caution)
.PHONY: reset-db
reset-db:
	docker exec -it saas_postgres psql -U $$POSTGRES_USER -d $$POSTGRES_DB -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@echo "Database reset complete."

# =========================
# Docker-specific backend restart
# =========================

## Restart backend services in Docker
.PHONY: restart-backend-docker
restart-backend-docker:
	$(DOCKER_COMPOSE) restart api worker




# # Start fullstack services
# make up

# # Run backend API locally (dev without Docker)
# make run-backend

# # Generate SQLC queries after changing DB schema
# make sqlc

# # Generate Swagger docs
# make swagger

# # Build frontend for production
# make build-frontend

# # Reset database schema
# make reset-db
