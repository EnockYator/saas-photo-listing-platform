
############################
# 1. Builder stage
############################
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy full source code
COPY . .

# Build single binary with CLI support
RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build -o app ./cmd/app

############################
# 2. Runtime stage (distroless)
############################

FROM alpine:3.20

WORKDIR /app

# Copy binary only
COPY --from=builder /build/app .

# Non-root user (security best practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

# Default port
EXPOSE 8080

# Default entrypoint
ENTRYPOINT ["./app"]
CMD [ "api" ]
