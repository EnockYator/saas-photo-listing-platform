# Photo Listing SaaS

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
[![CI](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml)

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org/)
[![Documentation](https://img.shields.io/badge/docs-overview-blue)](./docs/overview.md)
[![Documentation](https://img.shields.io/badge/docs-architecture-blue)](./ARCHITECTURE.md)

---

## **Introduction**

**Photo Listing SaaS** is a **multi-tenant**, **professional portfolio platform** built exclusively for **photographers**, **studios**, and **creative businesses**. Unlike generic photo-sharing platforms, we provide business-enabling tools with strict **tenant isolation**, subscription-based monetization, and enterprise-grade reliability all built with **Golang** for **maximum performance** and **maintainability**.

**Core Philosophy**: We don't compete with social networks. We provide the digital business backbone for photography professionals to showcase work, manage client deliveries, and grow their business through a polished, professional platform.

---

## Table of Contents

- [Photo Listing SaaS](#photo-listing-saas)
  - [Introduction](#introduction)
  - [Quick Start](#quick-start)
  - [Features](#features)
  - [Architecture](#architecture)
  - [Documentation](#documentation)
    - [Comprehensive Documentation](#comprehensive-documentation)
    - [Architecture Documentation](#architecture-documentation)
    - [Domain Layer Documentation](#domain-layer-documentation)
  - [Deployment](#deployment)
    - [Free Tier Deployment](#free-tier-deployment)
    - [Production Deployment](#production-deployment)
  - [Testing](#testing)
  [Contributing](#contributing)
  - [License](#license)
  - [Contact & Support](#contact-support)

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/EnockYator/photo-listing-saas.git
cd photo-listing-saas

# Start development environment
docker-compose up -d

# Run migrations
go run cmd/migrate/main.go up

# Start the server
go run cmd/server/main.go

```

### Access the services

- API http://localhost:8080
- MinIO Console: http://localhost:9001 (admin:admin)
- PostgreSQL: localhost:5432
- Redis CLI: docker exec -it photo-listing-redis redis-cli

---

## Features

### Multi-Tenancy
- Complete data isolation with Row-Level Security
- Tenant-aware resource allocation and limits
- Separate storage namespaces per tenant

### Professional Photography Tools
- Curated album management with lifecycle states
- Batch uploads with progress tracking
- Client proofing galleries with watermarking
- EXIF metadata preservation and search

### Business Features
- Three-tier subscription model (Free/Pro/Studio)
- Usage-based billing with automated enforcement
- Client engagement analytics and reporting
- White-label options for studios

### Enterprise Security
- JWT-based authentication with tenant context
- Signed URLs for media access
- Row-Level Security for data isolation
- Comprehensive audit logging

### Observability
- Structured logging with correlation IDs
- Prometheus metrics and Grafana dashboards
- Distributed tracing with OpenTelemetry
- Real-time performance monitoring

---

## Architecture

Photo Listing SaaS follows Clean Architecture and Domain-Driven Design principles:

Read [**artchitecture**](./ARCHITECTURE.md) documentation for more details on architecture.

```text

┌─────────────────────────────────────────┐
│      Presentation Layer (API/REST)      │
├─────────────────────────────────────────┤
│      Application Layer (Use Cases)      │
├─────────────────────────────────────────┤
│        Domain Layer (Business)          │
├─────────────────────────────────────────┤
│ Infrastructure Layer (DB/Storage/Queue) │
└─────────────────────────────────────────┘

```

### Key Technologies:

- **Backend**: [**Go (Golang**)](https://golang.org/), [**Gin**](https://gin-gonic.com/),
[**SQLC**](https://sqlc.dev/), **Validator**
- **Database**: [**PostgreSQL 18.1**](https://www.postgresql.org/) ([**Supabase**](https://supabase.com/)) with RLS
- **Image storage**: [**R2**](https://developers.cloudflare.com/r2/) ([S3](https://aws.amazon.com/s3/)-compatible)
- **Events**: [**NATS JetStream**](https://nats.io/) for reliable messaging
- **Infrastructure**: [**Fly.io**](https://fly.io/), [**Docker**](https://www.docker.com/), [**Cloudflare**](https://www.cloudflare.com/)

---

## Documentation

### Comprehensive Documentation

Comprehensive documentation is available in the /docs directory:

1. [**Overview**](./docs/overview.md) - Product vision and features
2. [**Setup Guide**](./docs/setup.md) - Development environment setup
3. [**Deployment**](./docs/deployment.md) - Production deployment guide
4. [**Security**](./docs/security.md) - Security architecture and practices
5. [**API Reference**](./docs/api.md) - REST API documentation
6. [**Architecture**](./ARCHITECTURE.md)

### Architecture Documentation

1. [**Architecture Context**](./docs/architecture/context.md) Context - Architectural decisions and patterns
2. [**Event System**](./docs/architecture/events.md) - Event-driven architecture details
3. [**Data Model**](./docs/architecture/data-model.md) - Database schema and design


### Domain Layer Documentation

[**Domain Documentation**](./backend/internal/domains/README.md) - Business domain models and rules

---

## Deployment

### Free Tier Deployment
```bash
# Deploy to Fly.io
fly deploy --app photo-listing-api

# Set environment secrets
fly secrets set DATABASE_URL="postgresql://..." \
                 R2_ACCESS_KEY="..." \
                 JWT_SECRET="..."

# Run migrations
fly ssh console -a photo-listing-api --command "./migrate up"

```

### Production Deployment

See [**deployment**](./docs/deployment.md) for complete deployment guide including:
- Multi-region setup
- Database backup strategies
- Monitoring and alerting configuration
- CI/CD pipeline setup

---

## Testing

```bash
# Run unit tests
go test ./... -v

# Run integration tests
docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Run end-to-end tests
go test ./tests/e2e -v

```

---

## Contributing

We welcome contributions! Please see our [**Contributing Guide**](./CONTRIBUTING.md) for details on:

- Code style and conventions
- Testing requirements
- Pull request process
- Development workflow

---

## License

This project is licensed under the [MIT License](https://opensource.org/license/mit) - see the [**LICENSE**](./LICENSE) file for details.

The SaaS service itself is proprietary, while the core infrastructure is open-source.

---

## Contact Support

- Website: https://photolisting.dev
- Email: support@photolisting.dev
- GitHub Issues: https://github.com/EnockYator/saas-photo-listing-platform/issues
- Documentation: https://photolisting.dev/docs

Built with ❤️ for photography professionals worldwide
