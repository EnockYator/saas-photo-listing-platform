.PHONY: help run build clean test

APP_NAME=api

help:
	@echo "Available commands:"
	@echo "  make run    - Run application"
	@echo "  make run-dev    - Run development server"
	@echo "  make build  - Build the application"
	@echo "  make test   - Run tests"
	@echo "  make vet   - Run tests"
	@echo "  make lint   - Run linter"
	@echo "  make fmt   - Format go source codes"
	@echo "  make migrate-up  - Apply migrations"
	@echo "  make migrate-down  - Downgrade migrations"
	@echo "  make migrate-create  - Create migrations"


include .env
export

run:
	APP_ENV=development go run ./cmd/api
	
run-staging:
	APP_ENV=staging go run ./cmd/api

run-production:
	APP_ENV=production go run ./cmd/api

build-api:
	mkdir -p bin
	go build -o bin/api ./cmd/api

build-migrate:
	go build -o bin/migrate ./cmd/migrate
	./migrate version

download-deps:
	go mod download

test:
	go test ./... -race

lint:
	golangci-lint run ./...

fmt:
	test -z "$(gofmt -l .)"

tidy:
	go mod tidy

vet:
	go vet ./...

migrate-up:
# 	migrate -path ./internal/infrastructure/database/postgres/migrations/ -database "$(DB_URL)" up
	go run ./cmd/migrate up

migrate-down:
# 	migrate -path ./internal/infrastructure/database/postgres/migrations/ -database "$(DB_URL)" down
	go run ./cmd/migrate down

# migrate-create:
# 	migrate create -ext sql -dir ./internal/infrastructure/database/postgres/migrations/ -seq $(name)

# migrate-version:
# 	go run ./cmd/migrate version

# migrate-force:
# 	go run ./cmd/migrate force $(version)
# # 	migrate -path ./internal/infrastructure/database/postgres/migrations/ -database "$(DB_URL)" force $(version)

# createdb:
# 	createdb $(name)

# dropdb:
# 	dropdb $(name)
