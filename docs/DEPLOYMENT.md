# Pinecone Deployment Guide

This guide covers deploying the Pinecone Recipe Management System in various environments.

## Table of Contents

1. [Local Development](#local-development)
2. [Production Deployment](#production-deployment)
3. [Environment Variables](#environment-variables)
4. [Database Setup](#database-setup)
5. [Security Considerations](#security-considerations)
6. [Backup & Recovery](#backup--recovery)
7. [Monitoring](#monitoring)

---

## Local Development

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 16
- npm or yarn

### Database Setup

1. **Install PostgreSQL 16**

2. **Create Database**
```bash
psql -U postgres
CREATE DATABASE pinecone;
CREATE USER pinecone_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE pinecone TO pinecone_user;
\q
```

3. **Set Database URL**
```bash
export DATABASE_URL="postgres://pinecone_user:your_password@localhost:5432/pinecone?sslmode=disable"
```

4. **Run Migrations**
```bash
cd backend
go install github.com/pressly/goose/v3/cmd/goose@latest
goose -dir internal/store/migrations postgres $DATABASE_URL up
```

### Backend Setup

```bash
cd backend
go mod download
go run cmd/server/main.go
```

Server starts at: `http://localhost:8080`

### Frontend Setup

```bash
cd frontend
npm install
npm run dev
```

Frontend starts at: `http://localhost:5173`

---

## Production Deployment

### Docker Compose (Recommended)

1. **Create `docker-compose.yml`**
```yaml
version: '3.8'

services:
  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_DB: pinecone
      POSTGRES_USER: pinecone_user
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U pinecone_user -d pinecone"]
      interval: 10s
      timeout: 5s
      retries: 5

  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      DATABASE_URL: postgres://pinecone_user:${DB_PASSWORD}@postgres:5432/pinecone?sslmode=disable
      JWT_SECRET: ${JWT_SECRET}
      PORT: 8080
      ALLOWED_ORIGINS: ${ALLOWED_ORIGINS}
      USDA_API_KEY: ${USDA_API_KEY}
    ports:
      - "8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
    restart: unless-stopped

  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
      args:
        VITE_API_URL: ${API_URL}
    ports:
      - "80:80"
    depends_on:
      - backend
    restart: unless-stopped

volumes:
  postgres_data:
```

2. **Create Backend Dockerfile**

Create `backend/Dockerfile`:
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /pinecone-api cmd/server/main.go

FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates postgresql-client

# Copy binary and migrations
COPY --from=builder /pinecone-api .
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY --from=builder /app/internal/store/migrations ./migrations

# Create uploads directory
RUN mkdir -p /app/uploads

EXPOSE 8080

# Run migrations then start server
CMD goose -dir ./migrations postgres $DATABASE_URL up && ./pinecone-api
```

3. **Create Frontend Dockerfile**

Create `frontend/Dockerfile`:
```dockerfile
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build with production API URL
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM nginx:alpine

# Copy built files
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx config
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
```

4. **Create Nginx Config**

Create `frontend/nginx.conf`:
```nginx
server {
    listen 80;
    server_name _;

    root /usr/share/nginx/html;
    index index.html;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # SPA routing
    location / {
        try_files $uri $uri/ /index.html;
    }

    # Cache static assets
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

5. **Create `.env` file**
```env
# Database
DB_PASSWORD=your_secure_password_here

# JWT
JWT_SECRET=your_jwt_secret_here_min_32_chars

# API
ALLOWED_ORIGINS=http://localhost,https://yourdomain.com
API_URL=http://localhost:8080/api/v1

# USDA (optional)
USDA_API_KEY=your_usda_api_key_here
```

6. **Deploy**
```bash
docker-compose up -d
```

### Manual Production Deployment

#### Backend

1. **Build binary**
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -o pinecone-api cmd/server/main.go
```

2. **Run migrations**
```bash
goose -dir internal/store/migrations postgres $DATABASE_URL up
```

3. **Create systemd service**

Create `/etc/systemd/system/pinecone-api.service`:
```ini
[Unit]
Description=Pinecone Recipe API
After=network.target postgresql.service

[Service]
Type=simple
User=pinecone
WorkingDirectory=/opt/pinecone
Environment="DATABASE_URL=postgres://pinecone_user:password@localhost:5432/pinecone"
Environment="JWT_SECRET=your_jwt_secret_here"
Environment="PORT=8080"
Environment="ALLOWED_ORIGINS=https://yourdomain.com"
ExecStart=/opt/pinecone/pinecone-api
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

4. **Enable and start**
```bash
sudo systemctl enable pinecone-api
sudo systemctl start pinecone-api
```

#### Frontend

1. **Build**
```bash
cd frontend
VITE_API_URL=https://api.yourdomain.com/api/v1 npm run build
```

2. **Deploy to web server**
```bash
# Copy dist/ to nginx web root
sudo cp -r dist/* /var/www/pinecone/
```

3. **Configure Nginx**

Create `/etc/nginx/sites-available/pinecone`:
```nginx
server {
    listen 80;
    server_name yourdomain.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name yourdomain.com;

    # SSL certificates (use certbot)
    ssl_certificate /etc/letsencrypt/live/yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/yourdomain.com/privkey.pem;

    root /var/www/pinecone;
    index index.html;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self' https://api.yourdomain.com;" always;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
```

---

## Environment Variables

### Backend

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DATABASE_URL` | PostgreSQL connection string | Yes | - |
| `JWT_SECRET` | Secret key for JWT signing (min 32 chars) | Yes | - |
| `PORT` | Server port | No | `8080` |
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | No | `http://localhost:5173` |
| `USDA_API_KEY` | USDA FoodData Central API key | No | Uses stub |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | No | `info` |
| `MAX_UPLOAD_SIZE` | Max image upload size in bytes | No | `5242880` (5MB) |

### Frontend

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `VITE_API_URL` | Backend API base URL | Yes | `http://localhost:8080/api/v1` |

---

## Database Setup

### Initial Setup

```bash
# Connect to PostgreSQL
psql -U postgres

# Create database and user
CREATE DATABASE pinecone;
CREATE USER pinecone_user WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE pinecone TO pinecone_user;

# Grant schema permissions
\c pinecone
GRANT ALL ON SCHEMA public TO pinecone_user;
GRANT ALL ON ALL TABLES IN SCHEMA public TO pinecone_user;
GRANT ALL ON ALL SEQUENCES IN SCHEMA public TO pinecone_user;
```

### Migrations

```bash
# Check migration status
goose -dir internal/store/migrations postgres $DATABASE_URL status

# Migrate up
goose -dir internal/store/migrations postgres $DATABASE_URL up

# Migrate down one version
goose -dir internal/store/migrations postgres $DATABASE_URL down

# Rollback to specific version
goose -dir internal/store/migrations postgres $DATABASE_URL down-to VERSION
```

### Migration Files

1. `001_create_users.sql` - Users and authentication
2. `002_create_recipes.sql` - Recipes, ingredients, instructions
3. `003_create_meal_plans.sql` - Meal planning tables
4. `004_create_grocery_lists.sql` - Grocery lists and items
5. `005_create_cookbooks.sql` - Cookbooks and cookbook-recipe junction
6. `006_create_nutrition_cache.sql` - Nutrition data caching

---

## Security Considerations

### Production Checklist

- [ ] Use strong, unique `JWT_SECRET` (min 32 characters)
- [ ] Use strong database passwords
- [ ] Enable SSL/TLS for database connections (`sslmode=require`)
- [ ] Use HTTPS for all production traffic
- [ ] Configure proper CORS origins (no wildcards)
- [ ] Keep dependencies updated (`go get -u`, `npm update`)
- [ ] Review npm audit warnings before deployment
- [ ] Enable firewall rules (only expose ports 80, 443)
- [ ] Use reverse proxy (nginx) with rate limiting
- [ ] Implement backup strategy (see below)
- [ ] Enable database connection pooling
- [ ] Set secure cookie attributes (HttpOnly, Secure, SameSite)
- [ ] Implement monitoring and alerting
- [ ] Regular security audits with `go vet` and `npm audit`

### Rate Limiting (Backend)

The application includes built-in rate limiting:
- 100 requests per minute per IP address
- Returns 429 Too Many Requests when exceeded
- Configurable in middleware

### HTTPS Setup

```bash
# Install certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d yourdomain.com

# Auto-renewal is configured automatically
```

---

## Backup & Recovery

### Database Backups

#### Automated Daily Backup

Create `/usr/local/bin/backup-pinecone-db.sh`:
```bash
#!/bin/bash
BACKUP_DIR="/var/backups/pinecone"
DATE=$(date +%Y%m%d_%H%M%S)
FILENAME="pinecone_backup_$DATE.sql.gz"

mkdir -p $BACKUP_DIR

# Dump and compress
pg_dump -h localhost -U pinecone_user -d pinecone | gzip > "$BACKUP_DIR/$FILENAME"

# Keep only last 30 days
find $BACKUP_DIR -name "pinecone_backup_*.sql.gz" -mtime +30 -delete

echo "Backup completed: $FILENAME"
```

Add to crontab:
```bash
# Daily backup at 2 AM
0 2 * * * /usr/local/bin/backup-pinecone-db.sh
```

#### Manual Backup

```bash
# Backup
pg_dump -h localhost -U pinecone_user -d pinecone > backup.sql

# Backup with compression
pg_dump -h localhost -U pinecone_user -d pinecone | gzip > backup.sql.gz
```

#### Restore

```bash
# Restore from plain SQL
psql -h localhost -U pinecone_user -d pinecone < backup.sql

# Restore from compressed
gunzip -c backup.sql.gz | psql -h localhost -U pinecone_user -d pinecone
```

### File Backups (Uploads)

```bash
# Backup uploads directory
tar -czf uploads_backup_$(date +%Y%m%d).tar.gz /app/uploads

# Restore
tar -xzf uploads_backup_20251109.tar.gz -C /
```

---

## Monitoring

### Health Check Endpoints

The backend provides health check endpoints:

```bash
# Basic health check
curl http://localhost:8080/health

# Database connection check
curl http://localhost:8080/health/db
```

### Application Logs

Logs are output in JSON format to stdout:

```bash
# View backend logs (systemd)
sudo journalctl -u pinecone-api -f

# View backend logs (Docker)
docker-compose logs -f backend

# Filter for errors only
docker-compose logs backend | grep '"level":"error"'
```

### Database Monitoring

```sql
-- Active connections
SELECT count(*) FROM pg_stat_activity WHERE datname = 'pinecone';

-- Table sizes
SELECT
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index usage
SELECT
    schemaname,
    tablename,
    indexname,
    idx_scan as index_scans
FROM pg_stat_user_indexes
WHERE schemaname = 'public'
ORDER BY idx_scan DESC;

-- Slow queries (if pg_stat_statements enabled)
SELECT
    query,
    calls,
    mean_exec_time,
    max_exec_time
FROM pg_stat_statements
ORDER BY mean_exec_time DESC
LIMIT 10;
```

### Performance Metrics

- **Backend compilation**: Clean build with no errors
- **Frontend TypeScript**: Zero type errors
- **Frontend bundle size**: 398.99 kB (115.34 kB gzipped)
- **Database indexes**: 24 indexes for query optimization
- **API response times**: Target <100ms for simple queries, <500ms for complex aggregations

### Recommended Monitoring Tools

- **Application**: Prometheus + Grafana
- **Database**: pg_stat_statements, pgAdmin
- **Logs**: ELK Stack (Elasticsearch, Logstash, Kibana)
- **Uptime**: UptimeRobot, Pingdom
- **Errors**: Sentry (optional)

---

## Troubleshooting

### Common Issues

**Database connection fails**
```bash
# Check PostgreSQL is running
sudo systemctl status postgresql

# Check connection
psql -h localhost -U pinecone_user -d pinecone

# Check logs
sudo journalctl -u postgresql -n 50
```

**Migrations fail**
```bash
# Check migration status
goose -dir internal/store/migrations postgres $DATABASE_URL status

# Manual rollback
goose -dir internal/store/migrations postgres $DATABASE_URL down

# Re-run migrations
goose -dir internal/store/migrations postgres $DATABASE_URL up
```

**Frontend can't connect to backend**
```bash
# Check VITE_API_URL is correct
# Check CORS allowed origins includes frontend domain
# Check backend is accessible
curl http://localhost:8080/api/v1/health
```

**JWT authentication fails**
```bash
# Ensure JWT_SECRET is set and consistent
# Check cookie settings (Secure flag requires HTTPS)
# Verify token expiry (24 hour default)
```

---

## Scaling Considerations

### Database

- **Connection pooling**: Configure `max_connections` in PostgreSQL
- **Read replicas**: For read-heavy workloads
- **Partitioning**: Partition large tables (recipes, meal_plans) by user_id
- **Archiving**: Archive old data (e.g., meal plans older than 1 year)

### Application

- **Horizontal scaling**: Run multiple backend instances behind load balancer
- **Caching**: Add Redis for session storage and frequently accessed data
- **CDN**: Serve static frontend files via CDN
- **Image optimization**: Resize uploaded images, serve WebP format

### Load Balancing

Example nginx load balancer config:
```nginx
upstream pinecone_backend {
    least_conn;
    server backend1:8080;
    server backend2:8080;
    server backend3:8080;
}

server {
    listen 80;
    location /api/ {
        proxy_pass http://pinecone_backend;
    }
}
```

---

## Maintenance

### Regular Tasks

- **Weekly**: Review logs for errors
- **Monthly**: Apply security updates, review database performance
- **Quarterly**: Test backup restoration, review storage usage
- **Annually**: Security audit, dependency updates

### Updating the Application

```bash
# Pull latest code
git pull origin main

# Backend
cd backend
go mod tidy
go build -o pinecone-api cmd/server/main.go
sudo systemctl restart pinecone-api

# Frontend
cd frontend
npm install
npm run build
sudo cp -r dist/* /var/www/pinecone/
```

---

For additional help, refer to:
- [README.md](../README.md) - Project overview and features
- [API.md](./API.md) - Complete API documentation
- [Backend README](../backend/README.md) - Backend-specific documentation
