
## /docs/security.md

```markdown
# Security

This document outlines the security architecture, practices, and compliance measures for Photo Listing SaaS.

## 🛡️ Security Philosophy

### Defense in Depth
We implement multiple layers of security controls throughout the application stack:

┌─────────────────────────────────────┐
│ Business Logic & Validation │ ← Domain-driven security
├─────────────────────────────────────┤
│ Application Security Controls │ ← AuthZ, input validation
├─────────────────────────────────────┤
│ Infrastructure Security │ ← Network, container, OS
├─────────────────────────────────────┤
│ Cloud Provider Security │ ← IAM, networking, WAF
└─────────────────────────────────────┘
text


### Principle of Least Privilege
Every component operates with the minimum permissions necessary to perform its function.

## 🔐 Authentication & Authorization

### JWT-Based Authentication

```go
// JWT token structure
type Claims struct {
    jwt.RegisteredClaims
    TenantID   string   `json:"tenant_id"`
    UserID     string   `json:"user_id"`
    Role       string   `json:"role"`
    Permissions []string `json:"permissions,omitempty"`
}

// Token generation with secure defaults
func GenerateToken(userID, tenantID, role string) (string, error) {
    claims := &Claims{
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   userID,
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "photo-listing-saas",
            ID:        uuid.New().String(), // JWT ID for replay prevention
        },
        TenantID:    tenantID,
        UserID:      userID,
        Role:        role,
        Permissions: getPermissionsForRole(role),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}
```

## Authentication flow

![Authentication Flow](./images/authentication_flow.svg)

Refresh Token Strategy
go

// Secure refresh token implementation
type RefreshToken struct {
    ID           string    `json:"id"`
    UserID       string    `json:"user_id"`
    TenantID     string    `json:"tenant_id"`
    TokenHash    string    `json:"token_hash"` // bcrypt hash
    UserAgent    string    `json:"user_agent"`
    IPAddress    string    `json:"ip_address"`
    CreatedAt    time.Time `json:"created_at"`
    LastUsedAt   time.Time `json:"last_used_at"`
    ExpiresAt    time.Time `json:"expires_at"`
}

// Single-use refresh tokens
func RotateRefreshToken(oldTokenHash, newToken string) error {
    // Delete old token
    db.Exec("DELETE FROM refresh_tokens WHERE token_hash = $1", oldTokenHash)
    
    // Store new token hash
    hashedToken, err := bcrypt.GenerateFromPassword([]byte(newToken), bcrypt.DefaultCost)
    db.Exec("INSERT INTO refresh_tokens (id, user_id, token_hash) VALUES ($1, $2, $3)",
        uuid.New(), userID, string(hashedToken))
    
    return nil
}

🏢 Multi-Tenancy Security
Row-Level Security (RLS)
sql

-- PostgreSQL RLS policies
-- Enable RLS on all tenant-owned tables
ALTER TABLE albums ENABLE ROW LEVEL SECURITY;
ALTER TABLE photos ENABLE ROW LEVEL SECURITY;
ALTER TABLE users ENABLE ROW LEVEL SECURITY;

-- Tenant isolation policy
CREATE POLICY tenant_isolation ON albums
    FOR ALL
    USING (tenant_id = current_setting('app.current_tenant_id')::uuid)
    WITH CHECK (tenant_id = current_setting('app.current_tenant_id')::uuid);

-- Public read access for published albums
CREATE POLICY public_read ON albums
    FOR SELECT
    USING (
        visibility = 'public' 
        AND status = 'published'
    );

-- Function to set tenant context
CREATE OR REPLACE FUNCTION set_tenant_context(tenant_id UUID)
RETURNS void AS $$
BEGIN
    PERFORM set_config('app.current_tenant_id', tenant_id::text, false);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

Tenant Context Propagation
go

// Tenant middleware ensures context is set
func TenantMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Extract tenant from JWT claims
        claims, _ := c.Get("claims")
        tenantID := claims.(*Claims).TenantID
        
        // Set in context
        ctx := context.WithValue(c.Request.Context(), "tenant_id", tenantID)
        c.Request = c.Request.WithContext(ctx)
        
        // Set for RLS
        db.ExecContext(ctx, "SELECT set_config('app.current_tenant_id', $1, false)", tenantID)
        
        c.Next()
    }
}

🔒 Data Protection
Encryption at Rest
Database Encryption
sql

-- Encrypt sensitive columns
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE clients (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone TEXT,
    -- Encrypted sensitive data
    notes BYTEA, -- Encrypted with tenant-specific key
    metadata JSONB,
    encryption_key_id TEXT -- Reference to KMS key
);

-- Encrypt data
INSERT INTO clients (notes, encryption_key_id)
VALUES (
    pgp_sym_encrypt('Sensitive notes', current_setting('app.encryption_key')),
    'key-123'
);

Storage Encryption
go

// Server-side encryption with R2
type EncryptedStorage struct {
    r2Client *R2Client
    kmsKeyID string
}

func (s *EncryptedStorage) Upload(ctx context.Context, tenantID, key string, data []byte) error {
    // Generate data key
    dataKey, err := s.generateDataKey()
    
    // Encrypt data
    encryptedData, err := s.encrypt(data, dataKey)
    
    // Store encrypted data
    return s.r2Client.PutObject(ctx, 
        fmt.Sprintf("tenants/%s/%s", tenantID, key),
        encryptedData,
    )
}

Encryption in Transit
yaml

# Traefik TLS configuration
tls:
  certificates:
    - certFile: /certs/cert.pem
      keyFile: /certs/key.pem
  options:
    default:
      minVersion: VersionTLS13
      cipherSuites:
        - TLS_AES_256_GCM_SHA384
        - TLS_CHACHA20_POLY1305_SHA256
        - TLS_AES_128_GCM_SHA256

Signed URLs for Media Access
go

func GenerateSignedURL(tenantID, objectKey string, expiresIn time.Duration) (string, error) {
    // Generate presigned URL
    presignedURL, err := s.r2Client.PresignGetObject(
        ctx,
        s.bucketName,
        fmt.Sprintf("tenants/%s/%s", tenantID, objectKey),
        expiresIn,
        nil,
    )
    
    // Add security headers
    presignedURL.Header.Set("Content-Disposition", "attachment")
    presignedURL.Header.Set("X-Content-Type-Options", "nosniff")
    
    return presignedURL.URL, nil
}

🚫 Input Validation & Sanitization
Structured Validation
go

// Request validation with custom rules
type CreateAlbumRequest struct {
    Title       string    `json:"title" validate:"required,min=3,max=255,alphanumspace"`
    Description string    `json:"description" validate:"max=1000"`
    Visibility  string    `json:"visibility" validate:"required,oneof=public private password client"`
    Password    string    `json:"password,omitempty" validate:"required_if=Visibility password,min=8"`
    Settings    JSONObject `json:"settings" validate:"validAlbumSettings"`
}

// Custom validation function
func validAlbumSettings(fl validator.FieldLevel) bool {
    settings := fl.Field().Interface().(JSONObject)
    
    // Validate watermark settings
    if wm, ok := settings["watermark"].(map[string]interface{}); ok {
        if opacity, ok := wm["opacity"].(float64); ok {
            return opacity >= 0 && opacity <= 1
        }
    }
    
    return true
}

SQL Injection Prevention
go

// Using sqlc for type-safe SQL
-- name: GetAlbum :one
SELECT * FROM albums
WHERE id = $1 AND tenant_id = $2;

// Generated Go code
func (q *Queries) GetAlbum(ctx context.Context, arg GetAlbumParams) (Album, error) {
    row := q.db.QueryRowContext(ctx, getAlbum, arg.ID, arg.TenantID)
    var i Album
    err := row.Scan(
        &i.ID,
        &i.TenantID,
        // ... other fields
    )
    return i, err
}

XSS Prevention
go

// Template auto-escaping
import "html/template"

func RenderAlbumPage(w http.ResponseWriter, album *Album) {
    tmpl := template.Must(template.New("album").Parse(`
        <h1>{{.Title}}</h1>
        <p>{{.Description}}</p>
        {{range .Photos}}
            <img src="{{.URL}}" alt="{{.AltText}}">
        {{end}}
    `))
    
    // Auto-escaped by html/template
    tmpl.Execute(w, album)
}

🛡️ Application Security
Rate Limiting
go

// Token bucket rate limiter per tenant
type RateLimiter struct {
    redis *redis.Client
}

func (rl *RateLimiter) Allow(ctx context.Context, tenantID, endpoint string) (bool, error) {
    key := fmt.Sprintf("ratelimit:%s:%s:%d", 
        tenantID, 
        endpoint,
        time.Now().Unix()/60, // New key per minute
    )
    
    // Use Redis INCR with expiry
    current, err := rl.redis.Incr(ctx, key).Result()
    if err != nil {
        return false, err
    }
    
    // Set expiry on first request
    if current == 1 {
        rl.redis.Expire(ctx, key, 65*time.Second)
    }
    
    // Get limits for tenant tier
    limit := rl.getLimit(tenantID)
    
    return current <= limit, nil
}

CSRF Protection
go

// Double-submit cookie pattern for SPA
func CSRFMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if c.Request.Method == "GET" {
            // Generate and set CSRF token
            token := generateCSRFToken()
            c.SetCookie("XSRF-TOKEN", token, 3600, "/", "", true, false)
            c.Set("csrf_token", token)
        } else {
            // Validate CSRF token
            headerToken := c.GetHeader("X-XSRF-TOKEN")
            cookieToken, _ := c.Cookie("XSRF-TOKEN")
            
            if headerToken == "" || headerToken != cookieToken {
                c.AbortWithStatusJSON(403, gin.H{"error": "CSRF token invalid"})
                return
            }
        }
        
        c.Next()
    }
}

Security Headers
go

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Header("Content-Security-Policy", 
            "default-src 'self'; "+
            "script-src 'self' 'unsafe-inline' https://cdn.jsdelivr.net; "+
            "style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; "+
            "img-src 'self' data: https:; "+
            "font-src 'self' https://fonts.gstatic.com; "+
            "connect-src 'self' https://api.photolisting.dev;")
        
        c.Header("X-Content-Type-Options", "nosniff")
        c.Header("X-Frame-Options", "DENY")
        c.Header("X-XSS-Protection", "1; mode=block")
        c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
        c.Header("Permissions-Policy", 
            "camera=(), microphone=(), geolocation=(), payment=()")
        
        c.Next()
    }
}

🏗️ Infrastructure Security
Network Security
yaml

# Docker Compose network segmentation
services:
  api:
    networks:
      - frontend
      - backend
  
  postgres:
    networks:
      - backend
  
  redis:
    networks:
      - backend

networks:
  frontend:
    driver: bridge
  
  backend:
    internal: true  # No external access
    driver: bridge
    driver_opts:
      com.docker.network.bridge.enable_icc: "false"

Container Security
dockerfile

# Security-hardened Dockerfile
FROM golang:1.21-alpine AS builder

# Add security scanning tools
RUN apk add --no-cache git gcc musl-dev

# Create non-root user
RUN adduser -D -g '' appuser

FROM scratch

# Copy only necessary files
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --chown=appuser:appuser photo-listing /app/

# Drop privileges
USER appuser

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["/app/healthcheck"]

EXPOSE 8080
ENTRYPOINT ["/app/photo-listing"]

Secrets Management
bash

# Using Fly.io secrets
fly secrets set JWT_SECRET=$(openssl rand -hex 32) \
                ENCRYPTION_KEY=$(openssl rand -hex 32) \
                DATABASE_URL="postgresql://..."

# Secret rotation script
#!/bin/bash
# rotate-secrets.sh
NEW_JWT_SECRET=$(openssl rand -hex 32)
fly secrets set JWT_SECRET=$NEW_JWT_SECRET

# Notify running instances to refresh
fly deploy --app photo-listing-api --strategy immediate

📜 Compliance & Auditing
Audit Logging
go

// Structured audit logging
type AuditLog struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time              `json:"timestamp"`
    TenantID    string                 `json:"tenant_id,omitempty"`
    UserID      string                 `json:"user_id,omitempty"`
    Action      string                 `json:"action"`
    Resource    string                 `json:"resource"`
    ResourceID  string                 `json:"resource_id,omitempty"`
    Changes     map[string]interface{} `json:"changes,omitempty"`
    IPAddress   string                 `json:"ip_address"`
    UserAgent   string                 `json:"user_agent"`
    RequestID   string                 `json:"request_id"`
}

// Log all mutating operations
func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Capture request details
        auditLog := AuditLog{
            ID:         uuid.New().String(),
            Timestamp:  time.Now(),
            Action:     c.Request.Method,
            Resource:   c.Request.URL.Path,
            IPAddress:  c.ClientIP(),
            UserAgent:  c.Request.UserAgent(),
            RequestID:  middleware.GetRequestID(c),
        }
        
        // Capture user context if available
        if claims, exists := c.Get("claims"); exists {
            auditLog.TenantID = claims.(*Claims).TenantID
            auditLog.UserID = claims.(*Claims).UserID
        }
        
        // Store audit log
        go auditRepository.Save(c.Request.Context(), auditLog)
        
        c.Next()
    }
}

GDPR Compliance
Data Portability
go

func ExportUserData(ctx context.Context, userID string) ([]byte, error) {
    data := map[string]interface{}{
        "profile":   getUserProfile(ctx, userID),
        "albums":    getUserAlbums(ctx, userID),
        "photos":    getUserPhotos(ctx, userID),
        "activity":  getUserActivity(ctx, userID),
    }
    
    return json.MarshalIndent(data, "", "  ")
}

Right to Erasure
go

func DeleteUserData(ctx context.Context, userID string) error {
    // Soft delete with anonymization
    return db.InTransaction(ctx, func(tx *sql.Tx) error {
        // Anonymize user data
        tx.Exec(`
            UPDATE users 
            SET email = 'deleted@user.com',
                name = 'Deleted User',
                avatar_url = NULL,
                metadata = '{}'
            WHERE id = $1
        `, userID)
        
        // Mark for permanent deletion
        tx.Exec(`
            INSERT INTO deletion_requests 
            (user_id, requested_at, schedule_for) 
            VALUES ($1, NOW(), NOW() + INTERVAL '30 days')
        `, userID)
        
        return nil
    })
}

Compliance Documentation
markdown

# Compliance Matrix

## GDPR
- [x] Data portability (export)
- [x] Right to erasure
- [x] Data processing agreements
- [x] Privacy policy
- [x] Cookie consent

## CCPA/CPRA
- [x] Right to know
- [x] Right to delete
- [x] Right to opt-out
- [x] "Do Not Sell" mechanism

## SOC 2
- [ ] Security controls documentation
- [ ] Regular audits
- [ ] Incident response plan
- [ ] Vendor management

🚨 Incident Response
Incident Response Plan
yaml

# incident-response.yaml
stages:
  identification:
    - Monitor alerts and logs
    - Confirm incident
    - Classify severity (P1-P4)
  
  containment:
    - Isolate affected systems
    - Disable compromised accounts
    - Preserve evidence
  
  eradication:
    - Remove malicious code
    - Patch vulnerabilities
    - Rotate credentials
  
  recovery:
    - Restore from backups
    - Verify integrity
    - Monitor for recurrence
  
  post_mortem:
    - Document root cause
    - Implement prevention
    - Update runbooks

Security Monitoring
go

// Anomaly detection
type AnomalyDetector struct {
    threshold float64
}

func (ad *AnomalyDetector) DetectAnomalies(ctx context.Context, metrics []Metric) []Anomaly {
    var anomalies []Anomaly
    
    for _, metric := range metrics {
        baseline := ad.getBaseline(metric.Name)
        
        // Check for statistical anomalies
        if math.Abs(metric.Value-baseline.Mean) > baseline.StdDev*ad.threshold {
            anomalies = append(anomalies, Anomaly{
                Metric:     metric.Name,
                Value:      metric.Value,
                Baseline:   baseline.Mean,
                Timestamp:  time.Now(),
            })
        }
    }
    
    return anomalies
}

Alert Configuration
yaml

# alertmanager.yml
route:
  group_by: ['alertname', 'tenant_id']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 12h
  receiver: 'telegram-alerts'
  
  routes:
  - match:
      severity: critical
    receiver: 'pagerduty'
    
  - match:
      severity: warning
    receiver: 'email'

receivers:
- name: 'telegram-alerts'
  telegram_configs:
  - bot_token: '{{ TELEGRAM_BOT_TOKEN }}'
    chat_id: '{{ TELEGRAM_CHAT_ID }}'
    message: '{{ .GroupLabels.alertname }}'
    
- name: 'pagerduty'
  pagerduty_configs:
  - service_key: '{{ PAGERDUTY_KEY }}'
  
- name: 'email'
  email_configs:
  - to: 'alerts@photolisting.dev'
    from: 'alerts@photolisting.dev'
    smarthost: 'smtp.resend.com:587'
    auth_username: 'apikey'
    auth_password: '{{ RESEND_API_KEY }}'

🧪 Security Testing
Automated Security Testing
yaml

# GitHub Actions security workflow
name: Security Scan
on: [push, pull_request]

jobs:
  gosec:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: securego/gosec@master
        with:
          args: ./...
  
  trivy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          scan-type: 'fs'
          scan-ref: '.'
          format: 'sarif'
          output: 'trivy-results.sarif'
  
  sast:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run Semgrep
        uses: returntocorp/semgrep-action@v1

Penetration Testing Checklist
markdown

# Penetration Testing Scope

## Authentication & Session Management
- [ ] Password policy enforcement
- [ ] Brute force protection
- [ ] Session fixation
- [ ] JWT security

## Authorization
- [ ] Horizontal privilege escalation
- [ ] Vertical privilege escalation
- [ ] IDOR vulnerabilities

## Input Validation
- [ ] SQL injection
- [ ] XSS attacks
- [ ] CSRF vulnerabilities
- [ ] File upload security

## Business Logic
- [ ] Workflow bypass
- [ ] Race conditions
- [ ] Business limit circumvention

## Infrastructure
- [ ] SSRF vulnerabilities
- [ ] File inclusion
- [ ] Command injection

Bug Bounty Program
markdown

# Security Bug Bounty

## Scope
- api.photolisting.dev
- app.photolisting.dev
- *.photolisting.dev

## Rewards
| Severity | Reward |
|----------|--------|
| Critical | $1,000 |
| High     | $500   |
| Medium   | $250   |
| Low      | $100   |

## Out of Scope
- DDoS attacks
- Social engineering
- Physical security
- Theoretical issues without proof

📚 Security Resources
Developer Security Training
markdown

# Required Security Training

## All Developers
- OWASP Top 10
- Secure coding practices
- GDPR basics

## Backend Developers
- SQL injection prevention
- Authentication security
- API security

## Frontend Developers
- XSS prevention
- CSRF protection
- Content Security Policy

## DevOps
- Container security
- Secrets management
- Infrastructure hardening

Security Tools
bash

# Development tools
go install github.com/securego/gosec/v2/cmd/gosec@latest
go install github.com/google/yamlfmt/cmd/yamlfmt@latest

# CI/CD tools
- Trivy for container scanning
- Semgrep for SAST
- OWASP ZAP for DAST
- GitGuardian for secrets detection

Security Documentation
bash

# Generate security documentation
make security-docs

# Update compliance matrix
make compliance-matrix

# Generate audit report
make audit-report

🔄 Continuous Security
Security as Code
yaml

# security-policies.yaml
policies:
  - name: "no-wildcard-permissions"
    description: "Prevent wildcard permissions in IAM policies"
    resource: "aws.iam.policy"
    filters:
      - "PolicyDocument.Statement[].Action": "*"
    actions:
      - type: "mark-for-op"
        op: "delete"
        days: 7
  
  - name: "encrypted-storage"
    description: "Ensure all storage is encrypted"
    resource: "aws.s3"
    filters:
      - type: "value"
        key: "ServerSideEncryption"
        value: "absent"
    actions:
      - type: "set-sse"
        crypto: "AES256"

Security Metrics Dashboard
json

{
  "security_metrics": {
    "vulnerabilities": {
      "open": 12,
      "critical": 1,
      "high": 3,
      "medium": 8
    },
    "compliance": {
      "gdpr": 95,
      "soc2": 85,
      "pci_dss": 70
    },
    "incidents": {
      "last_30_days": 2,
      "mttr_minutes": 45,
      "mttd_minutes": 15
    }
  }
}

🆘 Security Contact
Reporting Security Issues

Email: security@photolisting.dev
PGP Key: Download
Emergency Contact
yaml

# security-contacts.yaml
primary:
  name: "Security Team"
  email: "security@photolisting.dev"
  phone: "+1-XXX-XXX-XXXX"

backup:
  name: "CTO"
  email: "cto@photolisting.dev"
  phone: "+1-XXX-XXX-XXXX"

incident_response:
  - name: "Lead Developer"
    email: "dev@photolisting.dev"
  - name: "Operations"
    email: "ops@photolisting.dev"

This security documentation is a living document. Last updated: $(date)
All team members must review and acknowledge understanding of security policies.