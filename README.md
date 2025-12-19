# Photo Listing SaaS

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

A **production-grade Photo Listing SaaS platform** for uploading, organizing, watermarking, and securely sharing photos at scale.
Built with **Go (Golang), PostgreSQL, SQLC, and React**, the system is designed for **multi-tenant environments, role-based access control, and high-performance media handling**.
This system is production-ready and follows **clean architecture principles** with strict separation of concerns, ensuring maintainability, testability, and long-term scalability.

---

## Table of Contents

1. [Project Overview](#project-overview)  
2. [Core Features](#core-features)
3. [Architecture & Design Principles](#architecture--design-principles)
4. [Tech Stack](#tech-stack)  
5. [Project Structure](#project-structure)  
6. [Setup & Installation](#setup--installation)  
7. [Database & SQLC](#database--sqlc)  
8. [API Documentation](#api-documentation)  
9. [Environment Variables](#environment-variables)  
10. [Developer Workflow](#developer-workflow)  
11. [Contributing](#contributing)  
12. [License](#license)  

---

## Project Overview

**Photo Listing SaaS** enables individuals and organizations to manage large collections of photos efficiently while enforcing access control and subscription limits.

### Key Capabilities

- Upload and manage photos securely
- Group photos into structured listings
- Apply text or image-based watermarks
- Share photos publicly or privately
- Enforce subscription-based upload limits
- Expose an API-first backend for web and mobile clients

The platform is designed to scale horizontally and is suitable for:

- Photography portfolios
- Real estate listings
- Event photo delivery
- Media agencies
- Internal enterprise photo systems

---

## Core Features

- **Multi-Tenant Architecture**: Isolated data per tenant 
- **Role-Based Access Control**: Guest, Tenant, Viewer organize listings 
- **Photo Management**: Upload, update, delete, and organize photos into listings
- **Watermarking**: Protect media with configurable watermarks
- **Subscription Enforcement**: Upload and storage limits per plan
- **Cloud Object Storage**: S3-compatible storage (AWS S3 / MinIO)
- **RESTful API**: JSON-based, stateless endpoints
- **Async Workers**: Background image processing and long-running tasks
- **Docker-First Setup**: Consistent environments across development & production

---

## Architecture & Design Principles

The backend strictly follows **Clean Architecture** with clear separation of concerns.

### High-Level Flow

```pqsql
HTTP (Handlers)
      ↓
Service Layer (Use Cases)
      ↓
Domain (Business Rules & Entities)
      ↓
Storage (Postgres / S3)
      ↓
Database / Object Storage
```

### Layer responsibilities

- **HTTP Handlers**

    - Request validation
    - Response formatting
    - Authentication & authorization hooks

- **Service Layer**

    - Business logic and workflows
    - Coordinates domain entities and storage interfaces
    - No framework or infrastructure dependencies

- **Domain**

    - Core business entities
    - Business rules and invariants
    - Framework-agnostic

- **Storage**

    - Database persistence via SQLC
    - Object storage integrations (S3-compatible)

This structure ensures:
- Low coupling
- High testability
- Easy refactoring
- Long-term maintainability

---

## Tech Stack

**Backend**:

- [**Go (Golang**)](https://golang.org/) – High-performance, statically typed backend language
- [**Gin**](https://gin-gonic.com/) – Fast and minimalist HTTP framework
- [**SQLC**](https://sqlc.dev/) – Type-safe Go code generation from raw SQL queries
- [**PostgreSQL**](https://www.postgresql.org/) – Relational database for persistent data storage
- [**S3 / MinIO**](https://min.io/) – Object storage for photos and media assets
- [**Docker**](https://www.docker.com/) – Containerization for consistent runtime environments
- [**Docker Compose**](https://docs.docker.com/compose/) – Local multi-container orchestration 

**Frontend**:

- [**React**](https://react.dev/) – Component-based UI framework
- **Mantine / TailwindCSS** – UI styling 
- **Redux Toolkit Query (RTK Query)** – State management & API handling  

**Infrastructure & Tooling**:

- [**NGINX**](https://nginx.org/) – Reverse proxy and API gateway  
- **Makefile** – Developer automation and shortcut  

---

## Project Structure

```text
saas-photo-listing-platform/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go               # Application entry point
│   │
│   ├── deployments/
│   │   ├── docker-compose.yml        # Local orchestration (API, DB, etc.)
│   │   └── Dockerfile                # Backend container image
│   │
│   ├── go.mod                        # Go module definition
│   ├── go.sum                        # Go dependency checksums
│   │
│   ├── internal/
│   │   ├── db/
│   │   │   ├── schema/               # Database migrations
│   │   │   │   ├── 001_create_users.sql
│   │   │   │   └── 002_create_photos.sql
│   │   │   └── sqlc.yaml             # sqlc configuration
│   │   │
│   │   ├── domain/                   # Core business entities
│   │   │   ├── listing.go
│   │   │   ├── photo.go
│   │   │   ├── subscription.go
│   │   │   └── user.go
│   │   │
│   │   ├── http/
│   │   │   ├── handlers/             # HTTP request handlers
│   │   │   │   ├── auth_handler.go
│   │   │   │   ├── listing_handler.go
│   │   │   │   └── photo_handler.go
│   │   │   ├── router.go             # Route definitions
│   │   │   └── schemas/              # Request/response DTOs
│   │   │       ├── auth.go
│   │   │       ├── listing.go
│   │   │       └── photo.go
│   │   │
│   │   ├── middleware/               # HTTP middlewares
│   │   │   ├── auth.go               # Authentication middleware
│   │   │   └── middleware.go         # Common middleware utilities
│   │   │
│   │   ├── service/                   # Application use cases
│   │   │   ├── auth
│   │   │   │   └── service.go
│   │   │   ├── listing
│   │   │   │   └── service.go
│   │   │   ├── photo
│   │   │   │   └── service.go
│   │   │   └── subscription
│   │   │       └── service.go
│   │   ├── storage/
│   │   │   ├── postgres/             # PostgreSQL persistence
│   │   │   │   ├── db.go
│   │   │   │   ├── queries/           # SQL queries for sqlc
│   │   │   │   └── sqlc_generated/    # Generated Go code from sqlc
│   │   │   │
│   │   │   └── s3_storage/            # Object storage (S3-compatible)
│   │   │       └── s3_storage.go
│   │   │
│   │   ├── util/                     # Shared utilities
│   │   │   ├── helpers.go
│   │   │   └── validator.go
│   │   │
│   │   └── worker/                   # Background jobs & async processing
│   │
│   └── Makefile                     # Build & dev automation
│
├── frontend/                         # React frontend (UI)
│
├── nginx/
│   └── nginx.conf                   # Reverse proxy & static serving
│
└── README.md                        # Project documentation
```

---

## Setup & Installation

This project is **Docker-first**. All services (Backend, Frontend, PostgreSQL, NGINX) are containerized and orchestrated using **Docker Compose**.

---

### Prerequisites

- Docker
- Docker Compose

> You do **NOT** need Go, PostgreSQL, or Node.js installed locally.

---

### Quick Start

#### Clone the Repository

```bash
git clone https://github.com/EnockYator/saas-photo-listing-platform.git
cd saas-photo-listing-platform
```

---

#### Start all Services

```bash
docker-compose up --build
```

This command will:

- Build the Go backend

- Start PostgreSQL

- Apply database migrations

- Generate SQLC code

- Start the React frontend

- Configure NGINX as a reverse proxy

---

#### Access URLs

| Service       | URL                          | Port |
|---------------|------------------------------|------|
| Frontend      | http://localhost             | 80   |
| Backend API   | http://localhost/api         | 8080 |
| PostgreSQL    | localhost                    | 5432 |

---

#### Stopping the Application

```bash
docker compose down
docker compose down -v # Removes all volumes
```

---

## Database & SQLC

- All database schemas are defined using raw SQL

- SQLC generates type-safe Go code

- No ORM magic — queries are explicit and performant

This approach provides:

- Compile-time query validation

- Predictable SQL performance

- Clear ownership of database logic

---

## API Documentation

The API follows REST principles with standardized responses.

**Example: Create Photo**

```http
POST /photos
Authorization: Bearer <token>
Content-Type: application/json
```

Request Body:
```json
{
  "listing_id": 10,
  "title": "Wedding Shot"
}
```

Response (201 Created):
```json
{
  "id": 45,
  "title": "Wedding Shot",
  "listing_id": 10
}
```

All endpoints:

- Return consistent HTTP status codes

- Validate input data

- Provide structured error responses

---

## Environment Variables

Create a .env file in the root:

### Server
```ini
PORT=8080
ENV=development
```

### Database
```ini
POSTGRES_USER=your_postgres_username
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=saas_photo_listing_db

DB_HOST=postgres
DB_PORT=5432
DB_NAME=saas_photo_listing_db
DB_USER=postgres
DB_PASSWORD=your_postgres_password
```

### Object Storage
```ini
S3_ENDPOINT=https://s3.example.com
S3_ACCESS_KEY=your_s3_access_key
S3_SECRET_KEY=your_s3_secret_key
S3_BUCKET=photos
```

- **Important Note:** Never commit secrets to version control (Git).

---

## Developer Workflow

Data flow from client to database and back:

```text
Client (Web / Mobile)
   ↓
HTTP Handler
   ↓
Service Layer (Business Logic)
   ↓
Domain Models
   ↓
Storage (SQLC / S3)
   ↓
Database / Object Storage
```

**Rules enforced:**

    Handlers work with DTOs only

    Services operate on domain models

    Storage layers deal with persistence only

    All conversions happen at boundaries

This guarantees clean separation and long-term maintainability.

---

## Contributing

  1. Fork the repository

  2. Create a feature branch:
  ```bash
  git checkout -b feature/my-feature
  ```

  3. Commit changes: 
  ```bash
  git commit -m "Add your_feature_name"
  ```

  4. Push to branch:
  ```bash
  git push origin feature/your_feature_name
  ```

  5. Open a Pull Request

  ---

## License

This project is licensed under the **MIT License**.  
See the [LICENSE](LICENSE) file for details.

---

## Future Enhancements

- Swagger/OpenAPI integration for API docs
- CI/CD pipelines for automated testing and deployment
- Multi-region cloud deployment
- Role-based analytics dashboards


