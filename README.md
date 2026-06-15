# Saas Photo Listing Platform

[![Go Version](https://img.shields.io/badge/Go-1.23.0-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Docker](https://img.shields.io/badge/Docker-✓-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Release](https://img.shields.io/github/v/release/EnockYator/saas-photo-listing-platform)](https://github.com/EnockYator/saas-photo-listing-platform/releases)

[![Overview](https://img.shields.io/badge/docs-overview-blue)](./docs/overview.md)
[![Architecture](https://img.shields.io/badge/docs-architecture-blue)](./docs/architecture.md)

---

**NOTE: This application is in active development stage. Expect docummentation to be incomplete**

---

## CI/CD Pipeline

| Stage | Status |
|-------|--------|
| CI (Build & Test) | [![CI](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/ci.yml) |
| Linting | [![Lint](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/lint.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/lint.yml) |
| Security Scan | [![Security Scan](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/security.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/security.yml) |
| Docker Image Build | [![Docker Build](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/docker-build.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/docker-build.yml) |
| Database Migrations | [![Migrations](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/migrate.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/migrate.yml) |
| Staging Deployment | [![Deploy Staging](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/deploy-staging.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/deploy-staging.yml) |
| Production Deployment | [![Deploy Production](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/deploy-production.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/deploy-production.yml) |
| Release | [![Release](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/release.yml/badge.svg)](https://github.com/EnockYator/saas-photo-listing-platform/actions/workflows/release.yml) |

---

## Project Stats

[![Last Commit](https://img.shields.io/github/last-commit/EnockYator/saas-photo-listing-platform)](https://github.com/EnockYator/saas-photo-listing-platform/commits/main)
[![Open Issues](https://img.shields.io/github/issues/EnockYator/saas-photo-listing-platform)](https://github.com/EnockYator/saas-photo-listing-platform/issues)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)

---

## Introduction

**SaaS Photo Listing** is a **multi-tenant**, **professional portfolio platform** built exclusively for **photographers**, **creative businesses**, and **individuals** who want to showcase their work to clients. Unlike generic photo-sharing platforms, we provide business-enabling tools with strict **tenant isolation**, subscription-based monetization, and enterprise-grade reliability all built with **Golang** for **maximum performance** and **maintainability**.

**Core Philosophy**: We don't compete with social networks. We provide the digital business backbone for photography professionals to showcase work, manage client deliveries, and grow their business through a reliable, and professional platform.

---

## Table of Contents

- [Saas Photo Listing Platform](#saas-photo-listing-platform)
  - [CI/CD Pipeline](#cicd-pipeline)
  - [Project Stats](#project-stats)
  - [Introduction](#introduction)
  - [Table of Contents](#table-of-contents)
  - [Quick Start](#quick-start)
    - [Access the services](#access-the-services)
  - [Features](#features)
    - [Multi-Tenancy](#multi-tenancy)
    - [Professional Photography Tools](#professional-photography-tools)
    - [Business Features](#business-features)
    - [Enterprise Security](#enterprise-security)
    - [Observability](#observability)
  - [Architecture](#architecture)
    - [Key Technologies:](#key-technologies)
  - [Documentation](#documentation)
    - [Comprehensive Documentation](#comprehensive-documentation)
    - [Architecture Documentation](#architecture-documentation)
    - [Domain Layer Documentation](#domain-layer-documentation)
  - [Deployment](#deployment)
    - [Free Tier Deployment](#free-tier-deployment)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)
  - [Contact Support](#contact-support)
  - [Update expected soon](#update-expected-soon)

---

## Quick Start

```bash
# Clone the repository
git clone https://github.com/EnockYator/saas-photo-listing-platform.git
cd saas-photo-listing-platform

# Start development environment

# Run migrations

# Start the server

```

### Access the services

- API http://localhost:8080
- PostgreSQL: localhost:5432

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

SaaS Photo Listing follows Clean Architecture and Domain-Driven Design principles:

Read [**artchitecture**](./docs/architecture.md) documentation for more details on architecture.

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
- **Infrastructure**: [**Render**](https://render.com/), [**Docker**](https://www.docker.com/), [**Cloudflare**](https://www.cloudflare.com/)

---

## Documentation

### Comprehensive Documentation

Comprehensive documentation is available in the /docs directory:

1. [**Overview**](./docs/overview.md) - Product vision and features
2. [**Setup Guide**](./docs/setup.md) - Development environment setup
3. [**Deployment**](./docs/deployment.md) - Production deployment guide
4. [**Security**](./docs/security.md) - Security architecture and practices
5. [**API Reference**](./docs/api.md) - REST API documentation
6. [**Architecture**](./docs/architecture.md)

### Architecture Documentation

1. [**Architecture Context**](./docs/architecture/context.md) Context - Architectural decisions and patterns
2. [**Event System**](./docs/architecture/events.md) - Event-driven architecture details
3. [**Data Model**](./docs/architecture/data-model.md) - Database schema and design


### Domain Layer Documentation

[**Domain Documentation**](./backend/internal/domains/README.md) - Business domain models and rules

---

## Deployment

### Free Tier Deployment

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

- Service: https://saas-photo-listing-platform-production.onrender.com
- Email: ekyator02@gmail.com
- GitHub Issues: https://github.com/EnockYator/saas-photo-listing-platform/issues
- Documentation: https://saas-photo-listing-platform-production.onrender.com/docs

Built with ❤️ for photography professionals worldwide

## Update expected soon

_Last updated on 15th June 2026 by Author: Enock Yator_
