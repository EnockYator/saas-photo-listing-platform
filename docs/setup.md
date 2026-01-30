
## /docs/setup.md

```markdown
# Setup Guide

This guide will help you set up a local development environment for Photo Listing SaaS.

## 🚀 Quick Start

### Prerequisites

- **Go 1.21+**: [Installation guide](https://golang.org/doc/install)
- **Docker & Docker Compose**: [Installation guide](https://docs.docker.com/get-docker/)
- **Git**: [Installation guide](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- **Make** (optional): Usually pre-installed on Unix systems

### Clone the Repository

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/photo-listing-saas.git
cd photo-listing-saas

# Add upstream remote
git remote add upstream https://github.com/original-owner/photo-listing-saas.git

One-Command Setup (Linux/macOS)
bash

# Run the setup script
./scripts/setup-dev.sh

# Or use make
make setup

⚙️ Manual Setup
1. Configure Environment Variables
bash

# Copy the example environment file
cp .env.example .env

# Edit the file with your preferences
nano .env  # or use your preferred editor

Key Environment Variables:
bash

# Application
APP_ENV=development
APP_PORT=8080

# Database
DATABASE_URL=postgresql://postgres:postgres@localhost:5432/photo_listing?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379/0

# NATS
NATS_URL=nats://localhost:4222

# Storage (MinIO for development)
STORAGE_ENDPOINT=http://localhost:9000
STORAGE_ACCESS_KEY=admin
STORAGE_SECRET_KEY=admin123
STORAGE_BUCKET=photo-listing

# JWT
JWT_SECRET=your-development-jwt-secret-change-in-production
JWT_EXPIRY=24h

# Email (development)
EMAIL_PROVIDER=console  # Prints to console instead of sending

2. Start Dependencies with Docker
bash

# Start all services
docker-compose up -d

# Check services are running
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down

Services started:

    PostgreSQL: Database on port 5432

    Redis: Cache and sessions on port 6379

    NATS: Message broker on port 4222 (web UI on 8222)

    MinIO: S3-compatible storage on port 9000 (console on 9001)

    Adminer: Database GUI on port 8081 (optional)

3. Database Setup
bash

# Run migrations
go run cmd/migrate/main.go up

# Seed development data (optional)
go run cmd/seed/main.go

# Check database status
go run cmd/migrate/main.go status

4. Build and Run the Application
bash

# Build the application
go build -o photo-listing ./cmd/server

# Run the server
./photo-listing

# Or run directly
go run cmd/server/main.go

# Run with hot reload (using air)
air

Expected output:
text

INFO Starting Photo Listing SaaS server
INFO Environment: development
INFO Server starting on :8080
INFO Connected to PostgreSQL
INFO Connected to Redis
INFO Connected to NATS
INFO Storage configured: MinIO
INFO JWT middleware initialized

5. Verify the Setup
bash

# Check health endpoint
curl http://localhost:8080/health

# Expected response:
{
  "status": "healthy",
  "version": "0.1.0",
  "services": {
    "database": "connected",
    "redis": "connected",
    "nats": "connected",
    "storage": "connected"
  }
}

🔧 Development Tools
IDE Configuration
VS Code

Create .vscode/settings.json:
json

{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintFlags": ["--fast"],
  "go.formatTool": "goimports",
  "go.testFlags": ["-v"],
  "files.autoSave": "onFocusChange",
  "editor.formatOnSave": true,
  "editor.codeActionsOnSave": {
    "source.organizeImports": true
  }
}

GoLand/IntelliJ

    Enable Go modules

    Configure Go 1.21 SDK

    Set up run configurations for server and tests

Useful Make Commands
bash

# Development
make dev              # Start development server with hot reload
make watch            # Watch for changes and restart
make generate         # Generate code (mocks, SQL, etc.)

# Database
make migrate-up       # Run migrations
make migrate-down     # Rollback last migration
make migrate-status   # Check migration status
make seed             # Seed development data

# Testing
make test             # Run all tests
make test-unit        # Run unit tests
make test-integration # Run integration tests
make test-e2e         # Run E2E tests
make test-coverage    # Run tests with coverage
make lint             # Run linters

# Building
make build            # Build binary
make docker-build     # Build Docker image
make clean            # Clean build artifacts

Database Management
Using Adminer

    Visit http://localhost:8081

    System: PostgreSQL

    Server: postgres

    Username: postgres

    Password: postgres

    Database: photo_listing

Using psql
bash

# Connect to database
docker exec -it photo-listing-postgres psql -U postgres -d photo_listing

# Common commands
\dt                    # List tables
\d+ albums            # Describe albums table
SELECT * FROM tenants; # Query tenants
\q                     # Quit

Using Migrations
bash

# Create new migration
migrate create -ext sql -dir migrations -seq add_client_proofing

# Run migrations
go run cmd/migrate/main.go up

# Rollback last migration
go run cmd/migrate/main.go down 1

# Check migration status
go run cmd/migrate/main.go status

Storage Management (MinIO)
Access MinIO Console

    Visit http://localhost:9001

    Login: admin / admin123

    Create bucket: photo-listing

    Set bucket policy to public for development

Using mc (MinIO Client)
bash

# Configure alias
mc alias set local http://localhost:9000 admin admin123

# List buckets
mc ls local/

# List objects
mc ls local/photo-listing

# Copy file to bucket
mc cp image.jpg local/photo-listing/

🧪 Testing Environment
Test-Specific Setup
bash

# Start test containers
docker-compose -f docker-compose.test.yml up -d

# Run integration tests
go test ./tests/integration -v

# Run E2E tests
go test ./tests/e2e -v

# Clean up test containers
docker-compose -f docker-compose.test.yml down -v

Test Data Seeding
go

// tests/testdata/seed.go
func SeedTestTenant(db *sql.DB) (string, error) {
    tenantID := uuid.New().String()
    
    _, err := db.Exec(`
        INSERT INTO tenants (id, name, slug, plan)
        VALUES ($1, 'Test Tenant', 'test-tenant', 'free')
    `, tenantID)
    
    return tenantID, err
}

🔍 Debugging
Logging Levels
bash

# Set log level
LOG_LEVEL=debug go run cmd/server/main.go

# View structured logs
tail -f logs/development.log | jq .

Debug with Delve
bash

# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug server
dlv debug ./cmd/server

# Debug tests
dlv test ./internal/domain

Profiling
bash

# Enable pprof
export DEBUG=true

# Access pprof endpoints
go tool pprof http://localhost:8080/debug/pprof/profile
go tool pprof http://localhost:8080/debug/pprof/heap

# Generate flame graph
go tool pprof -http=:8081 http://localhost:8080/debug/pprof/profile

📦 Dependencies Management
Go Modules
bash

# Add new dependency
go get github.com/new/package@v1.2.3

# Update dependencies
go get -u ./...

# Tidy module
go mod tidy

# Verify dependencies
go mod verify

# Vendor dependencies (optional)
go mod vendor

Docker Dependencies
bash

# Update Docker images
docker-compose pull

# Rebuild containers
docker-compose build --no-cache

# Check for updates
docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.CreatedAt}}"

🐛 Common Issues
Port Conflicts

Problem: Port 8080 already in use
Solution: Change port in .env or kill existing process
bash

# Find process using port
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
APP_PORT=8082 go run cmd/server/main.go

Database Connection Issues

Problem: Cannot connect to PostgreSQL
Solution: Check if container is running and ports are correct
bash

# Check container status
docker-compose ps postgres

# Check logs
docker-compose logs postgres

# Test connection
pg_isready -h localhost -p 5432 -U postgres

Storage Access Issues

Problem: Cannot write to MinIO
Solution: Check bucket permissions and credentials
bash

# Create bucket if missing
docker exec photo-listing-minio mc mb local/photo-listing

# Set public policy
docker exec photo-listing-minio mc anonymous set public local/photo-listing

🚀 Next Steps

After setup is complete:

    Explore the API: Visit http://localhost:8080/docs for Swagger UI

    Run tests: make test to verify everything works

    Check out examples: Look at examples/ directory

    Create your first tenant: Use the seed script or API

    Start developing: Pick an issue from GitHub

📚 Additional Resources

    API Documentation - Complete API reference

    Architecture Overview - System architecture

    Deployment Guide - Production deployment

    Contributing Guide - How to contribute

🆘 Getting Help

If you encounter issues:

    Check this documentation

    Search existing GitHub issues

    Ask in GitHub Discussions

    Review logs in logs/ directory

    Enable debug mode: DEBUG=true go run cmd/server/main.go

Happy developing! If you have suggestions to improve this guide, please submit a PR.