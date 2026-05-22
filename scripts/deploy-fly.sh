#!/bin/bash
set -e

APP_NAME="saas-photo-listing-platform"
REGISTRY="ghcr.io"
IMAGE_TAG="${GITHUB_SHA:-latest}"

echo "Deploying to Fly.io..."

# Deploy API
flyctl deploy \
    --app $APP_NAME \
    --image $REGISTRY/$GITHUB_REPOSITORY-api:$IMAGE_TAG \
    --remote-only

# Run migrations
flyctl ssh console \
    --app $APP_NAME \
    --command "/app/migrate up"

echo "Deployment complete!"