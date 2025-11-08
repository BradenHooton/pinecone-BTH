# Deployment Guide
## Pinecone Recipe Management System

**Version:** 1.0  
**Date:** 2025-11-08  
**Target Environment:** Production (VPS)

---

## Table of Contents
1. [Prerequisites](#prerequisites)
2. [Server Setup](#server-setup)
3. [Initial Deployment](#initial-deployment)
4. [CI/CD Pipeline](#cicd-pipeline)
5. [Environment Configuration](#environment-configuration)
6. [Database Management](#database-management)
7. [Monitoring & Logs](#monitoring--logs)
8. [Backup & Recovery](#backup--recovery)
9. [Rollback Procedure](#rollback-procedure)
10. [Troubleshooting](#troubleshooting)

---

## Prerequisites

### Hardware Requirements

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| **CPU** | 2 cores | 4 cores |
| **RAM** | 2GB | 4GB |
| **Storage** | 20GB SSD | 50GB SSD |
| **Bandwidth** | 1TB/month | 2TB/month |

### Software Requirements

| Software | Version |
|----------|---------|
| **Ubuntu** | 24.04 LTS |
| **Docker** | 24+ |
| **Docker Compose** | 2.20+ |
| **Git** | 2.40+ |

### Domain & DNS

- Domain name (e.g., `pinecone.example.com`)
- DNS A record pointing to server IP
- Open ports: 80 (HTTP), 443 (HTTPS)

---

## Server Setup

### Step 1: Initial Server Configuration

```bash
# SSH into server
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Set hostname
hostnamectl set-hostname pinecone-prod

# Create deployment user
adduser deploy
usermod -aG sudo deploy
usermod -aG docker deploy

# Switch to deploy user
su - deploy
```

### Step 2: Install Docker

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# Verify installation
docker --version
docker-compose --version
```

### Step 3: Configure Firewall

```bash
# Install UFW (if not already installed)
sudo apt install ufw

# Allow SSH (change 22 if using custom port)
sudo ufw allow 22/tcp

# Allow HTTP and HTTPS
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# Enable firewall
sudo ufw enable

# Check status
sudo ufw status
```

### Step 4: Set Up SSH Keys

```bash
# On your local machine, generate SSH key if needed
ssh-keygen -t ed25519 -C "deploy@pinecone"

# Copy public key to server
ssh-copy-id deploy@your-server-ip

# Test passwordless login
ssh deploy@your-server-ip
```

### Step 5: Create Project Directory

```bash
# As deploy user
mkdir -p /var/www/pinecone
cd /var/www/pinecone

# Create required directories
mkdir -p uploads backups logs config
```

---

## Initial Deployment

### Step 1: Clone Repositories

```bash
cd /var/www/pinecone

# Clone backend
git clone https://github.com/bhooton/pinecone-api.git

# Clone frontend
git clone https://github.com/bhooton/pinecone-web.git
```

### Step 2: Configure Environment Variables

```bash
cd /var/www/pinecone/pinecone-api

# Create production .env file
nano .env.prod
```

**`.env.prod` Template:**
```bash
# Database
DB_USER=pinecone_prod
DB_PASSWORD=<STRONG_PASSWORD_HERE>
DATABASE_URL=postgres://pinecone_prod:<PASSWORD>@db:5432/pinecone

# JWT
JWT_SECRET=<GENERATE_256_BIT_SECRET>
JWT_EXPIRY_HOURS=24

# USDA API
USDA_API_KEY=<YOUR_USDA_API_KEY>
USDA_API_BASE_URL=https://api.nal.usda.gov/fdc/v1

# Server
SERVER_PORT=8080
SERVER_HOST=0.0.0.0
UPLOAD_DIR=/app/uploads
MAX_UPLOAD_SIZE_MB=5
ALLOWED_ORIGINS=https://pinecone.example.com
LOG_LEVEL=info

# Sentry (optional)
SENTRY_DSN=<YOUR_SENTRY_DSN>
```

**Generate Secrets:**
```bash
# Generate JWT secret (256-bit)
openssl rand -base64 32

# Generate strong DB password
openssl rand -base64 24
```

### Step 3: Create Production Docker Compose File

```bash
cd /var/www/pinecone/pinecone-api

nano docker-compose.prod.yml
```

**`docker-compose.prod.yml`:**
```yaml
version: '3.8'

services:
  db:
    image: postgres:16-alpine
    container_name: pinecone_db
    environment:
      POSTGRES_DB: pinecone
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - /var/www/pinecone/backups:/backups
    ports:
      - "127.0.0.1:5432:5432"
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - pinecone_network

  api:
    image: ghcr.io/bhooton/pinecone-api:latest
    container_name: pinecone_api
    env_file:
      - .env.prod
    volumes:
      - /var/www/pinecone/uploads:/app/uploads
      - /var/www/pinecone/config:/app/config:ro
    ports:
      - "127.0.0.1:8080:8080"
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
    networks:
      - pinecone_network

  caddy:
    image: caddy:2-alpine
    container_name: pinecone_caddy
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - /var/www/pinecone/pinecone-web/dist:/var/www/html:ro
      - caddy_data:/data
      - caddy_config:/config
      - /var/www/pinecone/logs:/var/log/caddy
    depends_on:
      - api
    restart: unless-stopped
    networks:
      - pinecone_network

volumes:
  postgres_data:
  caddy_data:
  caddy_config:

networks:
  pinecone_network:
    driver: bridge
```

### Step 4: Create Caddyfile

```bash
nano Caddyfile
```

**`Caddyfile`:**
```caddyfile
pinecone.example.com {
    # Reverse proxy API requests
    handle /api/* {
        reverse_proxy api:8080
    }

    # Serve uploaded images
    handle /uploads/* {
        root * /var/www/html
        file_server
    }

    # Serve React app
    handle {
        root * /var/www/html
        try_files {path} /index.html
        file_server
    }

    # Security headers
    header {
        X-Content-Type-Options "nosniff"
        X-Frame-Options "DENY"
        Referrer-Policy "no-referrer-when-downgrade"
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
    }

    # Enable compression
    encode gzip

    # Logging
    log {
        output file /var/log/caddy/access.log
        format json
    }
}
```

### Step 5: Run Database Migrations

```bash
cd /var/www/pinecone/pinecone-api

# Start database only
docker-compose -f docker-compose.prod.yml up -d db

# Wait for database to be ready
sleep 10

# Install Goose (if not on server)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# Verify
goose -dir internal/store/migrations postgres "$DATABASE_URL" status
```

### Step 6: Build and Deploy

```bash
# Pull latest images (or build if not using CI/CD yet)
docker-compose -f docker-compose.prod.yml pull

# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f

# Verify health
curl http://localhost:8080/health
curl https://pinecone.example.com/health
```

### Step 7: Create First User

```bash
# Open psql
docker exec -it pinecone_db psql -U pinecone_prod pinecone

# Insert first user (use bcrypt to hash password first)
# Example with bcrypt cost 12: $2a$12$...
INSERT INTO users (email, password_hash, name) 
VALUES ('admin@example.com', '<BCRYPT_HASH>', 'Admin User');
```

---

## CI/CD Pipeline

### GitHub Actions Workflow

**File:** `.github/workflows/cd.yml`

```yaml
name: Deploy to Production

on:
  push:
    branches:
      - main
    tags:
      - 'v*'

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Log in to GitHub Container Registry
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push Docker image
        uses: docker/build-push-action@v5
        with:
          context: .
          push: true
          tags: |
            ghcr.io/bhooton/pinecone-api:latest
            ghcr.io/bhooton/pinecone-api:${{ github.sha }}

      - name: Deploy to server
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.SSH_HOST }}
          username: ${{ secrets.SSH_USER }}
          key: ${{ secrets.SSH_PRIVATE_KEY }}
          script: |
            cd /var/www/pinecone/pinecone-api
            
            # Pull latest code
            git pull origin main
            
            # Pull new images
            docker-compose -f docker-compose.prod.yml pull
            
            # Restart services (only recreates changed containers)
            docker-compose -f docker-compose.prod.yml up -d
            
            # Prune old images
            docker image prune -f
            
            # Health check
            sleep 10
            curl -f http://localhost:8080/health || exit 1

      - name: Notify on failure
        if: failure()
        run: echo "Deployment failed! Check logs."
```

### GitHub Secrets Configuration

1. Go to GitHub repository → Settings → Secrets and variables → Actions
2. Add the following secrets:

| Secret | Value |
|--------|-------|
| `SSH_HOST` | Server IP or domain |
| `SSH_USER` | `deploy` |
| `SSH_PRIVATE_KEY` | Private SSH key for deploy user |
| `DOCKER_USERNAME` | GitHub username |
| `DOCKER_PASSWORD` | GitHub Personal Access Token (with `write:packages` scope) |

---

## Environment Configuration

### Staging Environment

**Purpose:** Mirror of production for UAT

**Setup:**
1. Use same server with different ports or subdomain (`staging.pinecone.example.com`)
2. Separate Docker Compose file: `docker-compose.staging.yml`
3. Separate `.env.staging` file
4. Deploy on every merge to `main`

### Production Environment

**Purpose:** Live user-facing application

**Setup:**
1. Deploy on Git tag (e.g., `v1.0.0`)
2. Manual approval required (GitHub Environments)
3. Rollback plan ready

---

## Database Management

### Regular Maintenance

```bash
# Vacuum database (weekly)
docker exec pinecone_db psql -U pinecone_prod pinecone -c "VACUUM ANALYZE;"

# Check database size
docker exec pinecone_db psql -U pinecone_prod pinecone -c "
  SELECT pg_size_pretty(pg_database_size('pinecone'));"

# Check table sizes
docker exec pinecone_db psql -U pinecone_prod pinecone -c "
  SELECT schemaname, tablename, 
         pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) 
  FROM pg_tables WHERE schemaname = 'public' 
  ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;"
```

### Schema Migrations in Production

```bash
# Always test in staging first!

# 1. Create backup before migration
./backup.sh

# 2. Run migration
cd /var/www/pinecone/pinecone-api
goose -dir internal/store/migrations postgres "$DATABASE_URL" up

# 3. Verify migration
goose -dir internal/store/migrations postgres "$DATABASE_URL" status

# 4. If migration fails, rollback
goose -dir internal/store/migrations postgres "$DATABASE_URL" down
```

---

## Monitoring & Logs

### Application Logs

```bash
# View API logs
docker-compose -f docker-compose.prod.yml logs -f api

# View Caddy logs
docker-compose -f docker-compose.prod.yml logs -f caddy

# View database logs
docker-compose -f docker-compose.prod.yml logs -f db

# Search logs for errors
docker-compose -f docker-compose.prod.yml logs api | grep -i error
```

### Sentry Integration

**Setup:**
1. Create Sentry project at https://sentry.io
2. Add `SENTRY_DSN` to `.env.prod`
3. Restart services
4. Verify errors are being tracked

### Health Checks

```bash
# API health endpoint
curl https://pinecone.example.com/api/v1/health

# Expected response:
# {"status": "ok", "timestamp": "2025-11-08T12:00:00Z"}

# Database health
docker exec pinecone_db pg_isready -U pinecone_prod

# Container status
docker ps --filter name=pinecone
```

### Resource Monitoring

```bash
# Container resource usage
docker stats

# Disk usage
df -h

# Memory usage
free -h

# CPU usage
top
```

---

## Backup & Recovery

### Automated Daily Backups

**Script:** `/var/www/pinecone/backup.sh`

```bash
#!/bin/bash
set -e

BACKUP_DIR="/var/www/pinecone/backups"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/pinecone_$TIMESTAMP.sql"

# Create backup
docker exec pinecone_db pg_dump -U pinecone_prod pinecone > "$BACKUP_FILE"

# Compress
gzip "$BACKUP_FILE"

# Keep only last 7 days
find "$BACKUP_DIR" -name "pinecone_*.sql.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_FILE.gz"
```

**Crontab Entry:**
```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * /var/www/pinecone/backup.sh >> /var/www/pinecone/logs/backup.log 2>&1
```

### Manual Backup

```bash
# Create backup
./backup.sh

# Or manually
docker exec pinecone_db pg_dump -U pinecone_prod pinecone > backup_$(date +%Y%m%d).sql
gzip backup_$(date +%Y%m%d).sql
```

### Restore from Backup

```bash
# Stop API to prevent writes
docker-compose -f docker-compose.prod.yml stop api

# Restore database
gunzip -c /var/www/pinecone/backups/pinecone_20250108_020000.sql.gz | \
  docker exec -i pinecone_db psql -U pinecone_prod pinecone

# Start API
docker-compose -f docker-compose.prod.yml start api

# Verify
curl https://pinecone.example.com/api/v1/health
```

---

## Rollback Procedure

### Scenario: New Deployment Breaks Production

**Quick Rollback:**
```bash
cd /var/www/pinecone/pinecone-api

# Revert to previous image
docker-compose -f docker-compose.prod.yml down
docker pull ghcr.io/bhooton/pinecone-api:<PREVIOUS_TAG>
docker-compose -f docker-compose.prod.yml up -d

# Or checkout previous Git commit
git log --oneline  # Find last working commit
git checkout <COMMIT_HASH>
docker-compose -f docker-compose.prod.yml up -d --build
```

### Scenario: Database Migration Failed

```bash
# Rollback migration
cd /var/www/pinecone/pinecone-api
goose -dir internal/store/migrations postgres "$DATABASE_URL" down

# Restore from backup
gunzip -c /var/www/pinecone/backups/pinecone_<TIMESTAMP>.sql.gz | \
  docker exec -i pinecone_db psql -U pinecone_prod pinecone

# Restart services
docker-compose -f docker-compose.prod.yml restart
```

---

## Troubleshooting

### Issue: 502 Bad Gateway

**Cause:** API container not responding

**Solution:**
```bash
# Check API container status
docker ps | grep pinecone_api

# Check logs
docker logs pinecone_api

# Restart API
docker-compose -f docker-compose.prod.yml restart api
```

### Issue: Database Connection Refused

**Cause:** Database container not running or wrong credentials

**Solution:**
```bash
# Check database container
docker ps | grep pinecone_db

# Check database logs
docker logs pinecone_db

# Verify .env.prod DATABASE_URL is correct

# Restart database
docker-compose -f docker-compose.prod.yml restart db
```

### Issue: SSL Certificate Error

**Cause:** Caddy failed to obtain Let's Encrypt certificate

**Solution:**
```bash
# Check Caddy logs
docker logs pinecone_caddy

# Common fixes:
# 1. Ensure DNS A record points to server
# 2. Ensure ports 80 and 443 are open
# 3. Restart Caddy
docker-compose -f docker-compose.prod.yml restart caddy
```

### Issue: Disk Space Full

**Solution:**
```bash
# Check disk usage
df -h

# Clean up Docker
docker system prune -a

# Clean up logs
sudo journalctl --vacuum-time=3d

# Clean up old backups
find /var/www/pinecone/backups -name "*.sql.gz" -mtime +30 -delete
```

---

## Security Checklist

- [ ] Firewall configured (UFW)
- [ ] SSH key-based authentication only
- [ ] Strong database password
- [ ] JWT secret is 256-bit random
- [ ] HTTPS enforced via Caddy
- [ ] `.env.prod` file permissions: `chmod 600`
- [ ] Regular security updates: `apt update && apt upgrade`
- [ ] Sentry error tracking enabled
- [ ] Database backups automated
- [ ] Fail2ban configured (optional but recommended)

---

## Post-Deployment Checklist

- [ ] All services running (`docker ps`)
- [ ] Health endpoint returns 200 (`curl https://pinecone.example.com/api/v1/health`)
- [ ] Can register and login
- [ ] Can create recipe
- [ ] Can plan meal
- [ ] Can generate grocery list
- [ ] Sentry receiving errors (test with 404)
- [ ] Backups running daily
- [ ] Logs accessible and readable

---

**Deployment Complete! 🚀**

For support or questions, refer to the project documentation or open a GitHub issue.
