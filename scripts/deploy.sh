#!/bin/bash

# Deployment script for Karima Store
# This script should be run on the VPS server

set -e

echo "🚀 Starting deployment..."

# Configuration
APP_DIR="/opt/karima-store"
DOCKER_COMPOSE_FILE="docker-compose.yml"
BACKUP_DIR="/opt/backups/karima-store"

# Create backup directory if it doesn't exist
mkdir -p "$BACKUP_DIR"

# Backup current database
echo "📦 Creating database backup..."
BACKUP_FILE="$BACKUP_DIR/db_backup_$(date +%Y%m%d_%H%M%S).sql"
docker exec karima_postgres pg_dump -U karima_store karima_db > "$BACKUP_FILE"
echo "✅ Database backed up to $BACKUP_FILE"

# Navigate to app directory
cd "$APP_DIR"

# Pull latest changes
echo "📥 Pulling latest Docker images..."
docker-compose pull

# Stop current containers
echo "🛑 Stopping current containers..."
docker-compose down

# Start new containers
echo "🚀 Starting new containers..."
docker-compose up -d

# Wait for services to be healthy
echo "⏳ Waiting for services to be healthy..."
sleep 10

# Health check
echo "🏥 Performing health check..."
if curl -f http://localhost:8080/api/v1/health; then
    echo "✅ Health check passed!"
else
    echo "❌ Health check failed!"
    echo "🔄 Rolling back..."
    docker-compose down
    # Restore from backup if needed
    exit 1
fi

# Clean up old images
echo "🧹 Cleaning up old Docker images..."
docker image prune -f

# Keep only last 7 backups
echo "🧹 Cleaning up old backups..."
cd "$BACKUP_DIR"
ls -t | tail -n +8 | xargs -r rm

echo "✅ Deployment completed successfully!"
