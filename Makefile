.PHONY: help run build clean test

help:
	@echo "Available commands:"
	@echo "  make run    - Run the application"
	@echo "  make build  - Build the application"
	@echo "  make test   - Run tests"
	@echo "  make clean  - Clean build artifacts"

run:
	go run cmd/api/main.go

# build:
# 	go build -o bin/api cmd/api/main.go

fmt:
	gofmt -w .

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test -v ./...

# clean:
# 	rm -rf bin/
# 	go clean

ci: fmt vet lint test