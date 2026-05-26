
############################
# 1. Builder stage
############################
FROM golang:1.23-alpine AS builder

WORKDIR /app

# System deps (for migrations and git-based versioning)
RUN apk add --no-cache git ca-certificates

# Dependencies first (better caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code
COPY . .

# Build single binary with CLI support (api + migrate)
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o app ./cmd/api

############################
# 2. Runtime stage
############################
FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary only
COPY --from=builder /app/app .

# Copy migrations
COPY --from=builder /app/internal/infrastructure/database/postgres/migrations \
    ./internal/infrastructure/database/postgres/migrations

# Non-root user (security best practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Default port
EXPOSE 8080

# Default entrypoint
ENTRYPOINT ["./app"]

# Default command (API mode)
CMD ["api"]
