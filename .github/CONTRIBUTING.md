# Contributing to SaaS-Photo-Listing-Platform

Thank you for contributing to saas-photo-listing-platform.

This project follows production-grade Go SaaS engineering best practices.

---

## Development Requirements

Install the following tools:

- Go 1.22+
- Docker
- Docker Compose
- PostgreSQL
- Git

---

## Project High-Level Structure

```text
saas-photo-listing-platform
├── .air.toml
├── api
│   ├── openapi
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   └── postman
│       └── saas_photo_listing_collection.json
├── cmd
│   ├── api
│   │   └── main.go
│   ├── migrate
│   │   └── main.go
│   ├── tools
│   │   └── main.go
│   └── worker
│       └── main.go
├── CODE_OF_CONDUCT.md
├── deployments
│   ├── docker
│   │   ├── Dockerfile
│   │   ├── Dockerfile.worker
│   │   └── .dockerignore
│   ├── docker-compose
│   │   ├── docker-compose.prod.yml
│   │   ├── docker-compose.staging.yml
│   │   └── docker-compose.yml
│   └── kubernetes
│       ├── configmap.yaml
│       ├── deployment.yaml
│       ├── ingress.yaml
│       ├── secrets.yaml
│       └── service.yaml
├── docker-compose.yml
├── docs
│   ├── api.md
│   ├── architecture.md
│   ├── deployment.md
│   ├── development.md
│   └── diagrams
│       ├── architecture.png
│       └── event_flow.png
├── .env
├── .github
│   ├── CODEOWNERS
│   ├── CONTRIBUTING.md
│   ├── ISSUE_TEMPLATE
│   │   ├── bug.md
│   │   ├── feature.md
│   │   └── refactor.md
│   ├── PULL_REQUEST_TEMPLATE.md.md
│   └── workflows
│       ├── ci.yml
│       ├── deploy.yml
│       ├── docker.yml
│       ├── lint.yml
│       ├── migrate.yml
│       ├── release.yml
│       └── security.yml
├── .gitignore
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go
│   ├── domains
│   │   ├── audit
│   │   ├── auth
│   │   ├── gallery
│   │   ├── health
│   │   ├── media
│   │   ├── notification
│   │   ├── payment
│   │   ├── sharing
│   │   ├── subscription
│   │   └── tenant
│   ├── infrastructure
│   │   ├── cache
│   │   ├── database
│   │   ├── messaging
│   │   ├── observability
│   │   └── storage
│   ├── interfaces
│   │   ├── cli
│   │   ├── grpc
│   │   └── http
│   ├── shared
│   │   ├── constants
│   │   ├── errors
│   │   ├── events
│   │   └── valueobjects
│   └── workers
│       ├── cleanup_jobs.go
│       ├── email_sender.go
│       ├── registry.go
│       ├── thumbnail_generator.go
│       └── usage_calculator.go
├── LICENSE
├── Makefile
├── pkg
│   ├── collections
│   │   └── slices.go
│   ├── context
│   │   └── context.go
│   ├── crypto
│   │   └── crypto.go
│   ├── datetime
│   │   └── datetime.go
│   ├── errors
│   │   └── errors.go
│   ├── idgenerator
│   │   └── generator.go
│   ├── json
│   │   └── json.go
│   ├── middleware
│   │   └── helpers.go
│   ├── pagination
│   │   └── pagination.go
│   ├── ratelimit
│   │   └── ratelimit.go
│   ├── response
│   │   └── response.go
│   ├── retry
│   │   └── retry.go
│   ├── strings
│   │   └── stringutils.go
│   ├── testing
│   │   ├── fixtures.go
│   │   └── mocks.go
│   ├── types
│   │   └── common.go
│   └── validator
│       └── validator.go
├── README.md
├── scripts
│   ├── backup_db.sh
│   ├── build.sh
│   ├── deploy.sh
│   ├── migrate.sh
│   └── seed.sh
├── tmp
│   └── build-errors.log
└── tools
    └── tools.go
```

---

## Branching Strategy

Use the following branch naming conventions:

- `main` → stable, production-ready
- `dev` → integration branch
- `feat/*` → adding new features
- `fix/*` → fix bugs
- `refactor/*` → refactor structure or files

Examples:

   feat/auth-service

   feat/payment-domain

   fix/jwt-validation

   refactor/upload-service

---

## Commit Message Convention

Follow the following conventions when writing commit messages

Examples:

   feat(domain): add authentication middleware
   fix(domain): resolve postgres connection leak
   refactor(domain): improve upload service architecture
   chore: update dependencies

---

## Workflow
1. Create an issue first
2. Branch from `dev`
3. Commit small, focused changes
4. Open PR early
5. Review before merge
6. Squash & merge for clean history

---

## Code Review
- Respond to comments
- Review codes

---

## Architecture Principles -DDD
- Domain layer is independent
- No infrastructure dependencies in domain
- Aggregates enforce invariants
- Application layer orchestrates use cases only

---

## Running the Application

1. ## Start PostgreSQL
   
   ```bash
   docker compose up -d postgres
   ```

2. ## Run Migrations
      
   ```bash
   go run ./cmd/migrate up
   ```
   
3. ## Start API server
      
   ```bash
   go run ./cmd/api
   ```

---

## Running Tests
      
   ```bash
   go test ./...
   ```

---

## Linting
      
   ```bash
   golangci-lint run ./...
   ```
---

## Pull Request Rules

Before opening a PR:

- Ensure tests pass
- Ensure lint passes
- Ensure migrations are reversible
- Keep PRs focused and small
- Update documentation if needed

---

## Database Migration Rules

Migration files must:

- Be reversible
- Avoid destructive operations when possible
- Be reviewed before merge
- Follow naming convention:

Example:

    00001_create_users_table.up.sql
    00001_create_users_table.down.sql

---

## Security Rules

Never:

- Commit secrets
- Commit .env files
- Hardcode credentials
- Expose internal tokens

Use GitHub Secret for storing credentials

---

_Updated on 16th May 2026 by Author: Enock Yator_



