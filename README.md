# Photo Listing SaaS

A **Photo Listing SaaS platform** built with **Go (Golang), PostgreSQL, SQLC, and React**, designed for multi-tenant, role-based access, and scalable photo management. This backend is production-ready and follows **clean architecture principles** with strict separation of concerns.

---

## Table of Contents

1. [Project Overview](#project-overview)  
2. [Features](#features)  
3. [Tech Stack](#tech-stack)  
4. [Project Structure](#project-structure)  
5. [Setup & Installation](#setup--installation)  
6. [Database & SQLC](#database--sqlc)  
7. [API Documentation](#api-documentation)  
8. [Environment Variables](#environment-variables)  
9. [Developer Workflow](#developer-workflow)  
10. [Contributing](#contributing)  
11. [License](#license)  

---

## Project Overview

This SaaS platform allows users to:

- Upload and manage photos
- Organize photos into listings
- Share photos publicly or privately
- Apply subscriptions and access limits
- Use an API-first backend for web/mobile clients

The backend follows **layered architecture**:

HTTP → Handlers → Services → Storage → Database


- **Handlers**: Handle requests/responses  
- **Services**: Contain business logic  
- **Storage**: Database access via SQLC  
- **Domain Models**: Pure business objects  

---

## Features

- **User Roles**: Admin, Tutor, Student  
- **Photo Management**: Upload, delete, organize listings  
- **Subscriptions**: Enforce upload limits based on plan  
- **Cloud Storage**: S3-compatible storage for media files  
- **REST API**: JSON request/response with proper status codes  
- **Async Workers**: Background processing for image tasks  

---

## Tech Stack

**Backend**:

- [Go (Golang)](https://golang.org/) – Server-side language  
- [SQLC](https://sqlc.dev/) – Type-safe SQL generation  
- PostgreSQL – Database  
- S3 / MinIO – File storage  
- Docker & Docker Compose – Containerization  

**Frontend**:

- React + Mantine / TailwindCSS  
- Redux Toolkit Query (RTK Query) – State management & API handling  

**Other Tools**:

- Nginx – Reverse proxy  
- Makefile – Task automation  

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

---

## Setup & Installation

### Prerequisites

- Go >= 1.21  
- PostgreSQL >= 15  
- Docker & Docker Compose  
- Node.js & npm (for frontend)

### Backend Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/saas-photo-listing-platform.git
cd saas-photo-listing-platform

# Install dependencies
go mod tidy

# Set up SQLC
sqlc generate

# Start the backend
go run cmd/api/main.go

Frontend Setup

(Assuming frontend is in frontend/ folder)

cd frontend
npm install
npm start

Database & SQLC

    Create a PostgreSQL database:

CREATE DATABASE saas_photo;

    Run migrations:

psql -U <username> -d saas_photo -f db/schema/001_create_users.sql
psql -U <username> -d saas_photo -f db/schema/002_create_photos.sql
psql -U <username> -d saas_photo -f db/schema/003_create_listings.sql

    Generate SQLC code:

sqlc generate

    SQLC generates type-safe Go structs and queries in internal/storage/postgres/sqlc_generated/.

API Documentation

All endpoints follow REST principles. Responses are in JSON.

Example: Create Photo

POST /photos
Authorization: Bearer <token>
Content-Type: application/json

Request Body:

{
  "listing_id": 10,
  "title": "Wedding Shot"
}

Response (201 Created):

{
  "id": 45,
  "title": "Wedding Shot",
  "listing_id": 10
}

    Status codes, validation, and error messages are standardized across all endpoints.

Environment Variables

Create a .env file in the root:

# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=saas_photo

# S3 Storage
S3_ENDPOINT=https://s3.example.com
S3_ACCESS_KEY=your_access_key
S3_SECRET_KEY=your_secret_key
S3_BUCKET=photos

Developer Workflow

Data flow from client to database and back:

Client (React/mobile)
   ↓
JSON Request
   ↓
HTTP Handler (http/handlers/photo_handler.go)
   └─ Converts JSON → API schema → Domain model
   ↓
Service Layer (service/photo_service.go)
   └─ Business logic on Domain model
   ↓
Storage Layer (storage/postgres/photo_repo.go + SQLC)
   └─ Domain model → SQLC Params → Database
Database
Database → SQLC Struct → Domain Model → Service → API Schema → JSON Response
   ↑
Client receives response

    Handlers: deal with API schemas only

    Services: operate on domain models only

    Storage: SQLC-generated structs only

    Conversions happen at layer boundaries

    This ensures clean separation of concerns and maintainable code.

Contributing

    Fork the repository

    Create a feature branch: git checkout -b feature/my-feature

    Commit changes: git commit -m "Add my feature"

    Push to branch: git push origin feature/my-feature

    Open a Pull R
