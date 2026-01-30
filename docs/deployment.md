
## /docs/deployment.md

```markdown
# Deployment Guide

This guide covers deploying Photo Listing SaaS to production environments.

## 🎯 Deployment Options

| Environment | Use Case | Complexity | Cost |
|------------|----------|------------|------|
| **Fly.io Free** | Development, Testing, MVP | Low | $0/month |
| **Fly.io Paid** | Production, Scaling | Medium | ~$50-200/month |
| **Kubernetes** | Enterprise, Custom Infrastructure | High | Variable |
| **Docker Swarm** | Small Production Cluster | Medium | ~$30-100/month |

## 🚀 Fly.io Deployment (Recommended)

### Prerequisites

1. **Fly.io Account**: [Sign up](https://fly.io/app/sign-up)
2. **Fly CLI**: [Installation guide](https://fly.io/docs/hands-on/install-flyctl/)
3. **Docker**: For building images
4. **GitHub Account**: For CI/CD (optional)

### 1. Initial Setup

```bash
# Login to Fly.io
fly auth login

# Create app
fly launch --name photo-listing-api --region iad --no-deploy

# Set secrets
fly secrets set APP_ENV=production \
                JWT_SECRET=$(openssl rand -hex 32) \
                ENCRYPTION_KEY=$(openssl rand -hex 32)

# Create volumes for persistent data
fly volumes create nats_data --region iad --size 1
fly volumes create redis_data --region iad --size 1

2. External Service Configuration
Supabase (Database)

    Create project at supabase.com

    Get connection string from Settings → Database

    Enable Row-Level Security

    Create initial tables with migrations

bash

# Set database secret
fly secrets set DATABASE_URL="postgresql://postgres.password@db.supabase.co:5432/postgres"

Cloudflare R2 (Storage)

    Create R2 bucket at Cloudflare Dashboard

    Generate API token with R2 read/write permissions

    Configure CORS for your domain

bash

# Set R2 secrets
fly secrets set R2_ACCESS_KEY="your-access-key" \
                R2_SECRET_KEY="your-secret-key" \
                R2_BUCKET_NAME="photo-listing" \
                R2_ENDPOINT="https://<account-id>.r2.cloudflarestorage.com"

Cloudflare (DNS & CDN)

    Add domain to Cloudflare

    Create CNAME record pointing to Fly.io

    Enable SSL/TLS (Full)

    Configure WAF rules

Resend (Email)

    Create account at resend.com

    Verify sending domain

    Get API key

bash

fly secrets set RESEND_API_KEY="re_123456789"

3. Deploy Applications
bash

# Deploy main API
fly deploy --app photo-listing-api

# Deploy NATS cluster
fly deploy --app photo-listing-nats

# Deploy Redis
fly deploy --app photo-listing-redis

# Deploy worker for background jobs
fly deploy --app photo-listing-worker

4. Run Database Migrations
bash

# SSH into app and run migrations
fly ssh console -a photo-listing-api --command "cd /app && ./migrate up"

# Or use a separate migration job
fly deploy --app photo-listing-migrate --command "./migrate up"

5. Verify Deployment
bash

# Check app status
fly status --app photo-listing-api

# Check logs
fly logs --app photo-listing-api

# Test health endpoint
curl https://api.yourdomain.com/health

# Test API endpoint
curl -H "Authorization: Bearer test" https://api.yourdomain.com/api/v1/health

🐳 Docker Swarm Deployment
1. Initialize Swarm
bash

# Initialize swarm on manager node
docker swarm init --advertise-addr <MANAGER-IP>

# Get join token for worker nodes
docker swarm join-token worker

# Deploy stack
docker stack deploy -c docker-compose.production.yml photo-listing

2. Production Docker Compose
yaml

# docker-compose.production.yml
version: '3.8'

services:
  traefik:
    image: traefik:v2.9
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/etc/traefik/traefik.yml
      - ./certs:/certs
    deploy:
      placement:
        constraints:
          - node.role == manager
    networks:
      - web

  api:
    image: yourregistry/photo-listing-api:${TAG:-latest}
    environment:
      - DATABASE_URL=postgresql://user:pass@postgres:5432/photo_listing
      - REDIS_URL=redis://redis:6379/0
      - NATS_URL=nats://nats:4222
    deploy:
      replicas: 3
      update_config:
        parallelism: 1
        delay: 10s
      restart_policy:
        condition: on-failure
    labels:
      - "traefik.http.routers.api.rule=Host(`api.example.com`)"
      - "traefik.http.services.api.loadbalancer.server.port=8080"
    networks:
      - web
      - internal

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: photo_listing
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./postgres.conf:/etc/postgresql/postgresql.conf
    deploy:
      placement:
        constraints:
          - node.role == manager
    networks:
      - internal

  # ... other services

☸️ Kubernetes Deployment
1. Cluster Setup
bash

# Using minikube for local testing
minikube start --cpus=4 --memory=8192

# Enable ingress
minikube addons enable ingress

# Set up PostgreSQL operator (optional)
kubectl apply -f https://raw.githubusercontent.com/zalando/postgres-operator/master/manifests/postgres-operator.yaml

2. Kubernetes Manifests
yaml

# k8s/api-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: photo-listing-api
  labels:
    app: photo-listing
    tier: api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: photo-listing
      tier: api
  template:
    metadata:
      labels:
        app: photo-listing
        tier: api
    spec:
      containers:
      - name: api
        image: yourregistry/photo-listing-api:1.0.0
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: photo-listing-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: photo-listing-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
# k8s/api-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: photo-listing-api
spec:
  selector:
    app: photo-listing
    tier: api
  ports:
  - port: 80
    targetPort: 8080
  type: ClusterIP
---
# k8s/ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: photo-listing-ingress
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
spec:
  tls:
  - hosts:
    - api.photolisting.dev
    secretName: photo-listing-tls
  rules:
  - host: api.photolisting.dev
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: photo-listing-api
            port:
              number: 80

3. Deploy to Kubernetes
bash

# Apply all manifests
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/configmaps.yaml
kubectl apply -f k8s/postgres.yaml
kubectl apply -f k8s/redis.yaml
kubectl apply -f k8s/nats.yaml
kubectl apply -f k8s/api-deployment.yaml
kubectl apply -f k8s/api-service.yaml
kubectl apply -f k8s/ingress.yaml

# Check deployment status
kubectl get pods -n photo-listing
kubectl get services -n photo-listing
kubectl get ingress -n photo-listing

🔐 Secrets Management
Environment Variables
bash

# Create .env.production
cat > .env.production << EOF
APP_ENV=production
DATABASE_URL=postgresql://...
JWT_SECRET=$(openssl rand -hex 32)
R2_ACCESS_KEY=...
R2_SECRET_KEY=...
RESEND_API_KEY=...
EOF

# Load environment
set -a && source .env.production && set +a

Using External Secret Managers
HashiCorp Vault
bash

# Store secret
vault kv put secret/photo-listing/production jwt_secret=$(openssl rand -hex 32)

# In deployment
vault read -field=jwt_secret secret/photo-listing/production

AWS Secrets Manager
yaml

# In deployment.yaml
env:
  - name: JWT_SECRET
    valueFrom:
      secretKeyRef:
        name: photo-listing-secrets
        key: jwt-secret

📊 Monitoring Setup
Grafana Cloud (Free Tier)

    Sign up at grafana.com

    Create API key with metrics push permissions

    Configure Prometheus remote write

bash

# Set Grafana secrets
fly secrets set GRAFANA_API_KEY="your-api-key" \
                GRAFANA_URL="https://your-instance.grafana.net"

# Deploy with monitoring
fly deploy --app photo-listing-api --build-arg MONITORING=true

Custom Prometheus Stack
yaml

# docker-compose.monitoring.yml
services:
  prometheus:
    image: prom/prometheus
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - prom_data:/prometheus
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
  
  grafana:
    image: grafana/grafana
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana_data:/var/lib/grafana
      - ./grafana/provisioning:/etc/grafana/provisioning
    ports:
      - "3000:3000"

🔄 CI/CD Pipeline
GitHub Actions
yaml

# .github/workflows/deploy.yml
name: Deploy to Production
on:
  push:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      - run: make test-coverage
      - uses: codecov/codecov-action@v3

  build:
    needs: test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Build and push Docker image
        uses: docker/build-push-action@v4
        with:
          push: true
          tags: |
            ghcr.io/${{ github.repository }}:${{ github.sha }}
            ghcr.io/${{ github.repository }}:latest
          secrets: |
            GITHUB_TOKEN=${{ github.token }}

  deploy:
    needs: build
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: superfly/flyctl-actions/setup-flyctl@master
      - run: flyctl deploy --app photo-listing-production --image ghcr.io/${{ github.repository }}:${{ github.sha }}
        env:
          FLY_API_TOKEN: ${{ secrets.FLY_API_TOKEN }}

GitLab CI
yaml

# .gitlab-ci.yml
stages:
  - test
  - build
  - deploy

test:
  stage: test
  image: golang:1.21
  script:
    - make test-coverage

build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

deploy:
  stage: deploy
  image: alpine:latest
  script:
    - apk add --no-cache curl
    - curl -L https://fly.io/install.sh | sh
    - flyctl deploy --app photo-listing-production --image $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

📈 Scaling Configuration
Horizontal Scaling
bash

# Scale API instances
fly scale count 3 --app photo-listing-api

# Scale memory
fly scale memory 512 --app photo-listing-api

# Scale to multiple regions
fly regions add ams fra sin --app photo-listing-api
fly scale count 2 --app photo-listing-api --region ams

Database Scaling
sql

-- Create read replicas
CREATE SUBSCRIPTION replica_subscription
CONNECTION 'host=replica.example.com port=5432 dbname=photo_listing user=replicator'
PUBLICATION master_publication;

-- Configure connection pooling
-- In pgbouncer.ini
[databases]
photo_listing = host=postgres.example.com port=5432 dbname=photo_listing

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 20

Storage Optimization
go

// Implement CDN caching
func GetPhotoURL(photoID, tenantID string) (string, error) {
    // Check CDN cache first
    if url := cdnCache.Get(photoID); url != "" {
        return url, nil
    }
    
    // Generate signed URL
    url := storage.GenerateSignedURL(tenantID, photoID, 24*time.Hour)
    
    // Cache in CDN
    cdnCache.Set(photoID, url, 24*time.Hour)
    
    return url, nil
}

🔒 Security Hardening
SSL/TLS Configuration
bash

# Generate SSL certificates with Let's Encrypt
fly certs create api.photolisting.dev --app photo-listing-api

# Check certificate status
fly certs show api.photolisting.dev --app photo-listing-api

# Force HTTPS redirect
fly.toml:
[[services]]
  http_options = { "h2_backend" = true, "response" = { "headers" = { "Strict-Transport-Security" = "max-age=31536000; includeSubDomains; preload" } } }

Network Security
yaml

# Docker network configuration
networks:
  frontend:
    driver: bridge
    driver_opts:
      com.docker.network.bridge.enable_icc: "false"
  
  backend:
    internal: true  # No external access

Container Security
dockerfile

# Security-hardened Dockerfile
FROM golang:1.21-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --chown=1000:1000 photo-listing /app/
USER 1000
EXPOSE 8080
ENTRYPOINT ["/app/photo-listing"]

📋 Health Checks
Application Health Endpoints
go

// Health check implementation
func HealthHandler(c *gin.Context) {
    checks := make(map[string]string)
    
    // Database check
    if err := db.Ping(); err != nil {
        checks["database"] = "unhealthy"
        c.JSON(503, gin.H{"status": "unhealthy", "checks": checks})
        return
    }
    checks["database"] = "healthy"
    
    // Redis check
    if err := redis.Ping(); err != nil {
        checks["redis"] = "unhealthy"
        c.JSON(503, gin.H{"status": "unhealthy", "checks": checks})
        return
    }
    checks["redis"] = "healthy"
    
    c.JSON(200, gin.H{"status": "healthy", "checks": checks})
}

External Monitoring
bash

# Use uptime robot or similar
# Configure alerts for:
# - HTTP 200 status
# - Response time < 500ms
# - SSL certificate expiry
# - Domain expiration

🗑️ Backup & Disaster Recovery
Database Backups
bash

#!/bin/bash
# backup.sh
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="backup_$DATE.sql"

# Backup PostgreSQL
pg_dump $DATABASE_URL --format=custom --no-owner > $BACKUP_FILE

# Upload to R2
rclone copy $BACKUP_FILE r2:photo-listing-backups/database/

# Cleanup old backups (keep 30 days)
find . -name "backup_*.sql" -mtime +30 -delete

Storage Backups
bash

# Sync storage to backup bucket
rclone sync r2:photo-listing/tenants/ r2:photo-listing-backups/storage/ \
  --checksum \
  --transfers=16 \
  --checkers=32

Recovery Procedures
bash

# Database recovery
pg_restore --clean --if-exists --dbname=photo_listing latest_backup.sql

# Storage recovery
rclone sync r2:photo-listing-backups/storage/ r2:photo-listing/tenants/

🚨 Emergency Procedures
Rollback Deployment
bash

# Fly.io rollback
fly deployments --app photo-listing-api
fly rollback --app photo-listing-api v123

# Kubernetes rollback
kubectl rollout history deployment/photo-listing-api
kubectl rollout undo deployment/photo-listing-api

# Docker Swarm rollback
docker service update --rollback photo-listing_api

Database Recovery
sql

-- Point-in-time recovery
-- In recovery.conf:
restore_command = 'cp /backup/wal/%f %p'
recovery_target_time = '2024-01-15 14:30:00'

Incident Response

    Identify: Check monitoring alerts and logs

    Contain: Scale down if needed, disable features

    Resolve: Apply fix or rollback

    Communicate: Update status page, notify users

    Post-mortem: Document root cause and prevention

📊 Cost Optimization
Free Tier Limits
Service	Free Tier	Monitoring
Fly.io	3 VMs × 256MB	Monitor usage with fly scale show
Supabase	500MB database	Check usage in dashboard
Cloudflare R2	10GB storage	Set up billing alerts
Resend	100 emails/day	Monitor sending volume
Grafana Cloud	10K metrics	Use efficient metric naming
Cost Monitoring
bash

# Fly.io cost tracking
fly billing show

# Set up cost alerts
fly alerts add cost --value 50 --currency usd

Optimization Tips

    Use CDN caching for static assets

    Enable compression for API responses

    Implement request coalescing

    Use efficient data formats (WebP, Protocol Buffers)

    Schedule resource-intensive jobs during off-peak

📝 Maintenance Tasks
Regular Maintenance
bash

# Weekly
- Check and apply security updates
- Review logs for anomalies
- Verify backups
- Clean up old data

# Monthly
- Rotate secrets (JWT, API keys)
- Update dependencies
- Review performance metrics
- Cost analysis

# Quarterly
- Security audit
- Disaster recovery test
- Infrastructure review
- Capacity planning

Update Procedures
bash

# Update Go version
go mod edit -go=1.22
go mod tidy

# Update Docker images
docker-compose pull
docker-compose up -d --force-recreate

# Update fly.toml for new features
fly deploy --app photo-listing-api

🆘 Troubleshooting
Common Issues

    High memory usage: Check for memory leaks, adjust VM size

    Database connection issues: Verify connection pool settings

    Slow uploads: Check network, storage provider, chunk size

    Authentication failures: Verify JWT configuration, clock sync

    Email delivery issues: Check SPF/DKIM/DMARC records

Debug Commands
bash

# Check application logs
fly logs --app photo-listing-api

# SSH into instance
fly ssh console -a photo-listing-api

# Check database connections
fly pg connections --app photo-listing-api

# Monitor metrics
fly dashboard --app photo-listing-api

📚 Additional Resources

    Fly.io Documentation

    Supabase Documentation

    Cloudflare R2 Documentation

    Kubernetes Documentation

    Docker Documentation

For production deployments, consider consulting with infrastructure experts or using managed services for critical components.