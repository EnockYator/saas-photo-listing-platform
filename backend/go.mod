module saas-photo-listing-platform

go 1.25.5

require (
    // Web Framework
    github.com/gin-gonic/gin v1.9.1
    
    // Database
    github.com/jackc/pgx/v5 v5.5.0
    
    // Database Tools
    github.com/golang-migrate/migrate/v4 v4.17.0
    github.com/sqlc-dev/sqlc v1.24.0
    
    // Cache (Optional, for later)
    github.com/redis/go-redis/v9 v9.4.0
    
    // Cloud Storage
    github.com/aws/aws-sdk-go-v2 v1.24.0
    github.com/aws/aws-sdk-go-v2/config v1.26.1
    github.com/aws/aws-sdk-go-v2/service/s3 v1.47.7
    github.com/minio/minio-go/v7 v7.0.66 // Alternative to AWS S3
    
    // Authentication & Security
    github.com/golang-jwt/jwt/v5 v5.2.0
    golang.org/x/crypto v0.17.0 // For bcrypt/argon2
    
    // Configuration
    github.com/joho/godotenv v1.5.1
    
    // Image Processing
    github.com/disintegration/imaging v1.6.2
    github.com/davidbyttow/govips/v2 v2.14.0 // High-performance libvips bindings
    golang.org/x/image v0.14.0
    
    // Rate Limiting
    github.com/ulule/limiter/v3 v3.11.2
    
    // Background Jobs/Queues
    github.com/hibiken/asynq v0.24.1 // Redis-based task queue
    github.com/riverqueue/river v0.0.12 // PostgreSQL-based job queue (alternative)
    
    // Email & Notifications
    github.com/wneessen/go-mail v0.4.1 // Email sending
    
    // Monitoring & Observability
    go.opentelemetry.io/otel v1.21.0
    go.opentelemetry.io/otel/trace v1.21.0
    github.com/prometheus/client_golang v1.18.0 // Metrics
    
    // Utilities
    github.com/google/uuid v1.5.0
    github.com/rs/zerolog v1.31.0 // Structured logging
    github.com/stretchr/testify v1.8.4 // Testing
    golang.org/x/sync v0.5.0 // For worker pools
)

require (
    // Indirect dependencies from main requires
    github.com/bytedance/sonic v1.10.1 // indirect
    github.com/cespare/xxhash/v2 v2.2.0 // indirect
    github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
    github.com/klauspost/cpuid/v2 v2.2.5 // indirect
    github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
    golang.org/x/arch v0.5.0 // indirect
    
    // sqlc generator dependencies (development tools)
    github.com/antlr/antlr4/runtime/Go/antlr/v4 v4.7.1 // indirect
    github.com/google/cel-go v0.12.6 // indirect
    github.com/iancoleman/strcase v0.2.0 // indirect
    github.com/inconshreveable/mousetrap v1.1.0 // indirect
    github.com/jinzhu/inflection v1.0.0 // indirect
    github.com/mattn/go-colorable v0.1.13 // indirect
    github.com/mattn/go-isatty v0.0.19 // indirect
    github.com/pganalyze/pg_query_go/v5 v5.1.0 // indirect
    github.com/spf13/cobra v1.7.0 // indirect
    github.com/spf13/pflag v1.0.5 // indirect
    github.com/stoewer/go-strcase v1.2.0 // indirect
    google.golang.org/protobuf v1.31.0 // indirect
)