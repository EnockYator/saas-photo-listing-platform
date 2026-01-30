
## /docs/architecture/context.md

```markdown
# Architecture Context

This document describes the architectural decisions, patterns, and principles that guide the development of Photo Listing SaaS.

## 🎯 Architectural Goals

### Primary Objectives
1. **Multi-Tenancy**: Complete data isolation between tenants
2. **Scalability**: Horizontal scaling for both reads and writes
3. **Maintainability**: Clear separation of concerns and clean code
4. **Reliability**: High availability with graceful degradation
5. **Security**: Defense in depth with principle of least privilege

### Quality Attributes
| Attribute | Priority | Strategy |
|-----------|----------|----------|
| **Availability** | High | Multi-region deployment, auto-scaling, circuit breakers |
| **Performance** | High | Caching, CDN, optimized queries, asynchronous processing |
| **Security** | Critical | RLS, JWT, encryption, audit logging, regular penetration testing |
| **Scalability** | High | Stateless services, sharding, read replicas |
| **Maintainability** | High | Clean architecture, comprehensive tests, documentation |
| **Portability** | Medium | Containerization, environment abstraction |

## 🏛️ Architectural Style

### Clean Architecture

We follow Robert C. Martin's Clean Architecture principles:

┌─────────────────────────────────────────┐
│ Frameworks & Drivers (External) │
├─────────────────────────────────────────┤
│ Interface Adapters (Controllers/Gateways)│
├─────────────────────────────────────────┤
│ Application Business Rules (Use Cases)│
├─────────────────────────────────────────┤
│ Enterprise Business Rules (Entities) │
└─────────────────────────────────────────┘
text


**Key Principles:**
1. **Independence from Frameworks**: Business logic doesn't depend on external frameworks
2. **Testable**: Business rules can be tested without UI, database, or external services
3. **Independent of UI**: UI can change without changing business rules
4. **Independent of Database**: Business rules not bound to database or ORM
5. **Independent of External Agencies**: Business rules don't know about external interfaces

### Domain-Driven Design

We apply DDD strategic and tactical patterns:

#### Bounded Contexts

┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│ Portfolio │ │ Client │ │ Billing │
│ Management │◄──┤ Proofing │◄──┤ & Payments │
│ │ │ │ │ │
│ - Albums │ │ - Selections │ │ - Subscriptions │
│ - Photos │ │ - Comments │ │ - Invoices │
│ - Media │ │ - Watermarks │ │ - Usage Tracking│
└─────────────────┘ └─────────────────┘ └─────────────────┘
text


#### Ubiquitous Language
- **Album**: A curated collection of photos with metadata and settings
- **Photo**: A media asset with multiple versions and metadata
- **Tenant**: A photography business with isolated data and resources
- **Client**: External user with limited access to specific albums
- **Proofing**: Process where clients select photos from an album

## 🏢 System Context

### System Landscape

```mermaid
C4Context
    title System Context Diagram
    
    Person(photographer, "Photographer", "Professional photographer managing portfolio and clients")
    Person(client, "Client", "Client viewing and selecting photos")
    Person(admin, "Platform Admin", "Administrator managing the SaaS platform")
    
    System(photo_listing, "Photo Listing SaaS", "Multi-tenant portfolio platform")
    
    System_Ext(supabase, "Supabase", "PostgreSQL database with auth")
    System_Ext(cloudflare, "Cloudflare", "CDN, DNS, and R2 storage")
    System_Ext(stripe, "Stripe", "Payment processing")
    System_Ext(resend, "Resend", "Email delivery")
    
    Rel(photographer, photo_listing, "Manages portfolio, uploads photos, shares with clients")
    Rel(client, photo_listing, "Views albums, selects photos, provides feedback")
    Rel(admin, photo_listing, "Manages platform, monitors performance")
    
    Rel(photo_listing, supabase, "Reads/Writes tenant data, authenticates users")
    Rel(photo_listing, cloudflare, "Stores media, caches content, routes traffic")
    Rel(photo_listing, stripe, "Processes payments, manages subscriptions")
    Rel(photo_listing, resend, "Sends transactional emails")
    ```
    ## Container Diagram

    ![container diagram](./../images/container_diagram.svg)

    Technology Decisions
Why Go?

Pros:

    Excellent performance and concurrency support

    Strong typing and compile-time safety

    Minimal runtime with small binary sizes

    Excellent standard library and tooling

    Growing ecosystem with strong community

Cons:

    Less mature web framework ecosystem than some languages

    Generics are relatively new (but we use Go 1.21+)

Why PostgreSQL with RLS?

Pros:

    Row-Level Security provides database-level tenant isolation

    JSONB support for flexible schemas

    Strong consistency and ACID compliance

    Excellent performance with proper indexing

    Mature ecosystem with proven reliability

Cons:

    More complex than NoSQL for some use cases

    Requires careful schema design for multi-tenancy

Why Cloudflare R2?

Pros:

    Zero egress fees (critical for media-heavy application)

    S3-compatible API

    Built-in CDN with Cloudflare

    Competitive pricing

    Strong durability guarantees

Cons:

    Less mature than AWS S3

    Some advanced features missing

Why NATS JetStream?

Pros:

    High performance message streaming

    At-least-once delivery guarantees

    Stream persistence and replay

    Simple operational model

    Go-native implementation

Cons:

    Less enterprise features than Kafka

    Smaller ecosystem

🏗️ Component Architecture
API Gateway Pattern
go

// API Gateway responsibilities:
// 1. Request routing
// 2. Authentication and authorization
// 3. Rate limiting
// 4. Request/response transformation
// 5. Caching
// 6. Monitoring and logging

type APIGateway struct {
    router        *gin.Engine
    authService   AuthService
    rateLimiter   RateLimiter
    cache         Cache
    metrics       MetricsCollector
}

func (g *APIGateway) RegisterRoutes() {
    // Public routes
    g.router.GET("/health", g.healthCheck)
    g.router.GET("/public/albums/:id", g.getPublicAlbum)
    
    // Authenticated routes
    auth := g.router.Group("/api/v1")
    auth.Use(g.authMiddleware)
    auth.Use(g.rateLimitMiddleware)
    auth.Use(g.tenantMiddleware)
    
    auth.GET("/albums", g.listAlbums)
    auth.POST("/albums", g.createAlbum)
    // ... other routes
}
```

## Sevrvice Decomposition

![service decomposition](./../images/service_decomposition.svg)

Event-Driven Communication
go

// Event-driven service communication
type EventBus struct {
    natsConn *nats.Conn
    js       nats.JetStreamContext
}

func (eb *EventBus) PublishAlbumCreated(ctx context.Context, album *domain.Album) error {
    event := domain.AlbumCreated{
        AlbumID:     album.ID,
        TenantID:    album.TenantID,
        Title:       album.Title,
        CreatedAt:   time.Now(),
    }
    
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }
    
    // Publish to NATS JetStream
    _, err = eb.js.Publish(ctx, "events.album.created", data)
    return err
}

// Event handlers in different services
func RegisterEventHandlers(eventBus *EventBus) {
    // Analytics service listens for album events
    eventBus.Subscribe("events.album.*", func(msg *nats.Msg) {
        // Update analytics
    })
    
    // Notification service listens for client events
    eventBus.Subscribe("events.client.*", func(msg *nats.Msg) {
        // Send notifications
    })
}

📦 Deployment Architecture
Multi-Region Deployment
yaml

# fly.toml - Multi-region configuration
app = "photo-listing-api"
primary_region = "iad"

[[services]]
  internal_port = 8080
  protocol = "tcp"
  
  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]
  
  [[services.http_checks]]
    interval = "15s"
    timeout = "10s"
    method = "GET"
    path = "/health"
  
  # Multi-region deployment
  [[services.regions]]
    code = "iad"
    count = 2
  [[services.regions]]
    code = "lhr"
    count = 2
  [[services.regions]]
    code = "sin"
    count = 2

Database Topology
sql

-- Primary database in IAD with read replicas
-- Primary: iad (writes)
-- Replica 1: lhr (reads for EU)
-- Replica 2: sin (reads for Asia)
-- Replica 3: iad (analytical queries)

-- Use read-only transactions for replicas
BEGIN TRANSACTION READ ONLY;
SET LOCAL read_only = true;
-- Query goes to read replica
COMMIT;
```

## CDN and Edge Caching

![cdn and edge catching](./../images/cdn_and_edge_catching.svg)

🔄 Data Flow Examples
Photo Upload Flow

![photo upload flow](./../images/photo_upload_flow.svg)

## Client Proofing Flow

![client proofing flow](./../images/client_proofing_flow.svg)

Design Patterns
Repository Pattern
go

// Generic repository interface
type Repository[T any] interface {
    FindByID(ctx context.Context, id string) (*T, error)
    Save(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id string) error
    List(ctx context.Context, filter Filter) ([]*T, *Pagination, error)
}

// Album repository implementation
type AlbumRepository struct {
    db *sql.DB
    queries *sqlc.Queries
}

func (r *AlbumRepository) Save(ctx context.Context, album *domain.Album) error {
    // Convert domain model to persistence model
    dbAlbum := toDBAlbum(album)
    
    // Use sqlc for type-safe SQL
    return r.queries.SaveAlbum(ctx, dbAlbum)
}

// Mock repository for testing
type MockAlbumRepository struct {
    albums map[string]*domain.Album
}

func (m *MockAlbumRepository) Save(ctx context.Context, album *domain.Album) error {
    m.albums[album.ID] = album
    return nil
}

Factory Pattern
go

// Domain factory for creating valid aggregates
type AlbumFactory struct {
    idGenerator IDGenerator
    clock       Clock
    validator   Validator
}

func (f *AlbumFactory) CreateAlbum(tenantID, title, description string) (*domain.Album, error) {
    // Validate input
    if err := f.validator.ValidateAlbumTitle(title); err != nil {
        return nil, err
    }
    
    // Create aggregate with invariants enforced
    album := &domain.Album{
        ID:          f.idGenerator.Generate(),
        TenantID:    tenantID,
        Title:       title,
        Description: description,
        Status:      domain.StatusDraft,
        Visibility:  domain.VisibilityPrivate,
        CreatedAt:   f.clock.Now(),
        UpdatedAt:   f.clock.Now(),
        Version:     1,
    }
    
    return album, nil
}

Strategy Pattern
go

// Storage strategy for different providers
type StorageStrategy interface {
    Upload(ctx context.Context, file io.Reader, key string) (string, error)
    Download(ctx context.Context, key string) (io.ReadCloser, error)
    Delete(ctx context.Context, key string) error
    GenerateSignedURL(ctx context.Context, key string, expires time.Duration) (string, error)
}

// R2 implementation
type R2Storage struct {
    client *s3.Client
    bucket string
}

// Local implementation for development
type LocalStorage struct {
    basePath string
}

// Factory to get appropriate strategy
func NewStorageStrategy(env string) StorageStrategy {
    switch env {
    case "production":
        return NewR2Storage()
    case "development":
        return NewLocalStorage()
    default:
        return NewMockStorage()
    }
}

Observer Pattern (Events)
go

// Event publisher
type EventPublisher struct {
    observers []EventObserver
}

func (p *EventPublisher) Subscribe(observer EventObserver) {
    p.observers = append(p.observers, observer)
}

func (p *EventPublisher) Publish(event Event) {
    for _, observer := range p.observers {
        observer.OnEvent(event)
    }
}

// Event observers
type AnalyticsObserver struct {
    analyticsClient AnalyticsClient
}

func (o *AnalyticsObserver) OnEvent(event Event) {
    switch e := event.(type) {
    case AlbumCreated:
        o.analyticsClient.Track("album_created", e.TenantID)
    case PhotoUploaded:
        o.analyticsClient.Track("photo_uploaded", e.TenantID)
    }
}

📐 Architecture Decisions Records (ADRs)
ADR-001: Multi-Tenancy with RLS

Status: Accepted
Date: 2024-01-15
Context: Need to isolate data between tenants in a SaaS application.
Decision: Use PostgreSQL Row-Level Security for data isolation.
Consequences:

    ✅ Database-level security guarantees

    ✅ Simple application code

    ✅ Standard PostgreSQL feature

    ❌ Requires careful session management

    ❌ More complex backups/restores

ADR-002: Event-Driven Architecture

Status: Accepted
Date: 2024-01-15
Context: Need loose coupling between services and reliable event processing.
Decision: Use NATS JetStream for event streaming with at-least-once delivery.
Consequences:

    ✅ Loose coupling between services

    ✅ Reliable event processing

    ✅ Stream replay capability

    ❌ Additional operational complexity

    ❌ Eventual consistency in some cases

ADR-003: Go for Backend

Status: Accepted
Date: 2024-01-15
Context: Need high-performance, maintainable backend services.
Decision: Use Go 1.21+ for all backend services.
Consequences:

    ✅ Excellent performance

    ✅ Strong typing and safety

    ✅ Small binary sizes

    ❌ Less mature web framework ecosystem

    ❌ Learning curve for some teams

ADR-004: Clean Architecture

Status: Accepted
Date: 2024-01-15
Context: Need maintainable codebase with clear separation of concerns.
Decision: Implement Clean Architecture with clear layer boundaries.
Consequences:

    ✅ Testable business logic

    ✅ Framework independence

    ✅ Clear separation of concerns

    ❌ More boilerplate code

    ❌ Steeper learning curve

🔮 Evolution & Scaling
Monolith to Microservices

Current (Phase 1): Modular Monolith
text

┌─────────────────────────────────────┐
│          Photo Listing API          │
│                                     │
│  ┌─────┐ ┌─────┐ ┌─────┐ ┌─────┐   │
│  │Album│ │Photo│ │User │ │Billing│  │
│  │Module││Module││Module││Module │  │
│  └─────┘ └─────┘ └─────┘ └─────┘   │
└─────────────────────────────────────┘

Future (Phase 2): Microservices
text

┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐
│ Album   │ │ Photo   │ │ User    │ │ Billing │
│ Service │ │ Service │ │ Service │ │ Service │
└─────────┘ └─────────┘ └─────────┘ └─────────┘
        ↘           ↙       ↘           ↙
        └─────────────────────────────┘
                API Gateway

Database Scaling Strategy

    Phase 1: Single PostgreSQL instance with RLS

    Phase 2: Read replicas for different regions

    Phase 3: Sharding by tenant (using Citus or custom)

    Phase 4: Specialized databases (TimescaleDB for analytics)

Caching Evolution
go

// Phase 1: Single Redis instance
type CacheV1 struct {
    redis *redis.Client
}

// Phase 2: Multi-level caching
type CacheV2 struct {
    local  *ristretto.Cache  // In-memory LRU
    redis  *redis.Client     // Shared cache
    cdn    CDNCache          // Edge cache
}

// Phase 3: Regional caching
type CacheV3 struct {
    regional map[string]*redis.Client  // Cache per region
    global   *redis.Client             // Global cache
}

🧪 Testing Strategy
Test Pyramid
text

        ┌─────────────────┐
        │   E2E Tests     │ ← 5% of tests
        │   (User Journeys)│
        └─────────────────┘
               │
        ┌─────────────────┐
        │ Integration     │ ← 15% of tests
        │ Tests           │
        │ (Service Boundaries)│
        └─────────────────┘
               │
        ┌─────────────────┐
        │   Unit Tests    │ ← 80% of tests
        │   (Components)  │
        └─────────────────┘

Test Environment
yaml

# docker-compose.test.yml
version: '3.8'
services:
  postgres-test:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: test_db
      POSTGRES_USER: test
      POSTGRES_PASSWORD: test
    
  redis-test:
    image: redis:7-alpine
    
  nats-test:
    image: nats:2.9-alpine
    command: -js
    
  test-runner:
    build: .
    command: go test ./... -v
    depends_on:
      - postgres-test
      - redis-test
      - nats-test
    environment:
      DATABASE_URL: postgresql://test:test@postgres-test:5432/test_db

📊 Monitoring & Observability
Three Pillars Implementation
go

// Logging
logger.Info("request completed",
    "method", req.Method,
    "path", req.URL.Path,
    "duration_ms", duration.Milliseconds(),
    "tenant_id", tenantID,
    "request_id", requestID,
)

// Metrics
metrics.RequestDuration.WithLabelValues(
    req.Method,
    req.URL.Path,
    tenantID,
).Observe(duration.Seconds())

// Tracing
ctx, span := tracer.Start(ctx, "HandleRequest")
defer span.End()
span.SetAttributes(
    attribute.String("http.method", req.Method),
    attribute.String("tenant.id", tenantID),
)

Key Dashboards

    Business Dashboard: Tenant growth, revenue, conversion rates

    Performance Dashboard: API latency, error rates, cache hit rates

    Infrastructure Dashboard: CPU, memory, disk, network

    Security Dashboard: Failed logins, suspicious activity, rate limits

🔄 Deployment Pipeline
yaml

# GitHub Actions workflow
name: Deploy
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make test-all
  
  security-scan:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make security-scan
  
  build:
    needs: [test, security-scan]
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - run: make docker-build
  
  deploy-staging:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --app photo-listing-staging
  
  deploy-production:
    needs: deploy-staging
    if: github.ref == 'refs/heads/production'
    runs-on: ubuntu-latest
    steps:
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --app photo-listing-production

📚 Further Reading

    Clean Architecture by Robert C. Martin

    Domain-Driven Design by Eric Evans

    Building Microservices by Sam Newman

    PostgreSQL Row-Level Security

    NATS JetStream Documentation

This document is a living document. Last updated: $(date)
For questions or suggestions, create an issue in the GitHub repository.


