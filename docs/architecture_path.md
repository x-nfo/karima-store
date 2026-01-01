# Architecture Path - Karima Store
## Comprehensive Development vs Production Guide

---

## 📋 Overview

Sistem Karima Store menggunakan **microservices architecture** dengan pemisahan domain untuk keamanan dan isolasi sesi.

---

## 🏗️ 1. Domain Strategy

### **Production Architecture**

```
┌──────────────────────────────────────────────────────────────┐
│                      PRODUCTION DOMAINS                       │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  👥 Customer Storefront                                      │
│  https://karima.com                                          │
│  └─ Cloudflare Pages (Static Frontend)                      │
│                                                               │
│  🔐 Admin Panel                                              │
│  https://admin.ks-backend.cloud                              │
│  └─ Cloudflare Pages + Zero Trust Access                    │
│                                                               │
│  🚀 API Backend                                              │
│  https://api.ks-backend.cloud                                │
│  └─ VPS Docker Container (Go API)                           │
│                                                               │
│  🔑 Authentication (Ory Kratos)                              │
│  https://auth.ks-backend.cloud                               │
│  └─ VPS Docker Container (Kratos)                           │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

### **Development Architecture (Current)**

```
┌──────────────────────────────────────────────────────────────┐
│                    DEVELOPMENT (localhost)                    │
├──────────────────────────────────────────────────────────────┤
│                                                               │
│  👥 Customer Storefront (Future)                             │
│  http://localhost:3000                                       │
│  └─ Next.js / React dev server                              │
│                                                               │
│  🔐 Kratos Self-Service UI                                   │
│  http://127.0.0.1:4455                                       │
│  └─ Ory Kratos UI Container                                 │
│                                                               │
│  🚀 API Backend                                              │
│  http://localhost:8080                                       │
│  └─ Podman Container (Go API)                               │
│                                                               │
│  🔑 Kratos Public API                                        │
│  http://127.0.0.1:4433                                       │
│  └─ Ory Kratos Container                                    │
│                                                               │
│  🔧 Kratos Admin API                                         │
│  http://127.0.0.1:4434                                       │
│  └─ Ory Kratos Admin (Internal)                             │
│                                                               │
└──────────────────────────────────────────────────────────────┘
```

---

## 🔐 2. Authentication Flow

### **Production Flow**

```
1. User → https://auth.ks-backend.cloud/login
   └─ Kratos UI for login/register

2. Kratos → Set session cookie
   Domain: .ks-backend.cloud
   Cookie: ory_kratos_session
   HttpOnly: true
   Secure: true
   SameSite: Lax

3. Frontend → https://api.ks-backend.cloud/api/v1/...
   └─ Cookie automatically sent (same domain)

4. API → Validate session with Kratos
   Internal call: http://kratos:4433/sessions/whoami
   
5. API → Return data or 401/403
```

### **Development Flow**

```
1. User → http://127.0.0.1:4455/login
   └─ Kratos Self-Service UI

2. Kratos → Set session cookie
   Domain: 127.0.0.1
   Cookie: ory_kratos_session
   HttpOnly: true
   Secure: false (HTTP in dev)
   SameSite: Lax

3. Browser → http://localhost:8080/api/v1/...
   ⚠️ ISSUE: Cookie won't be sent (different domain!)
   └─ Need to use Bearer token instead

4. API → Validate session
   Call: http://kratos:4433/sessions/whoami
   
5. API → Return data or 401/403
```

---

## ⚙️ 3. Configuration Differences

### **Environment Variables**

#### **Production (.env.production)**

```env
# Application
APP_ENV=production
APP_PORT=8080

# Database (Production)
DB_HOST=postgres.ks-backend.cloud
DB_PORT=5432
DB_USER=karima_prod
DB_PASSWORD=<strong-password>
DB_NAME=karima_prod
DB_SSL_MODE=require

# Redis (Production)
REDIS_HOST=redis.ks-backend.cloud
REDIS_PORT=6379
REDIS_PASSWORD=<strong-password>

# CORS - Production Domains
CORS_ORIGIN=https://karima.com,https://admin.ks-backend.cloud

# Ory Kratos (Production)
KRATOS_PUBLIC_URL=https://auth.ks-backend.cloud
KRATOS_ADMIN_URL=http://kratos:4434  # Internal Docker network
KRATOS_UI_URL=https://auth.ks-backend.cloud

# Midtrans (Production)
MIDTRANS_SERVER_KEY=<production-server-key>
MIDTRANS_CLIENT_KEY=<production-client-key>
MIDTRANS_API_BASE_URL=https://app.midtrans.com/snap/v1
MIDTRANS_IS_PRODUCTION=true

# Cloudflare R2 (Production)
R2_BUCKET_NAME=karima-production
R2_PUBLIC_URL=https://cdn.karima.com

# Security
JWT_SECRET=<strong-random-secret>
```

#### **Development (.env - Current)**

```env
# Application
APP_ENV=development
APP_PORT=8080

# Database (Local)
DB_HOST=db
DB_PORT=5432
DB_USER=karima_store
DB_PASSWORD=lokal
DB_NAME=karima_db
DB_SSL_MODE=disable

# Redis (Local)
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# CORS - Local Development
CORS_ORIGIN=http://localhost:3000

# Ory Kratos (Local)
KRATOS_PUBLIC_URL=http://127.0.0.1:4433
KRATOS_ADMIN_URL=http://127.0.0.1:4434
KRATOS_UI_URL=http://127.0.0.1:4455

# Midtrans (Sandbox)
MIDTRANS_SERVER_KEY=YOUR_SERVER_KEY
MIDTRANS_CLIENT_KEY=YOUR_CLIENT_KEY
MIDTRANS_API_BASE_URL=https://app.sandbox.midtrans.com/snap/v1
MIDTRANS_IS_PRODUCTION=false

# Cloudflare R2 (Dev)
R2_BUCKET_NAME=dev-testing

# Security (Dev - Weak is OK)
JWT_SECRET=super_secret_key
```

---

## 🌐 4. CORS Configuration

### **Production CORS**

```go
// internal/middleware/cors.go (Production)
app.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "https://karima.com",
        "https://admin.ks-backend.cloud",
    },
    AllowMethods: []string{
        "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
        "X-Requested-With",
        "Cookie",
    },
    ExposeHeaders: []string{
        "Content-Length",
        "Content-Type",
        "Set-Cookie",
    },
    AllowCredentials: true,  // ✅ CRITICAL for HttpOnly cookies
    MaxAge:           86400,
}))
```

### **Development CORS**

```go
// internal/middleware/cors.go (Development)
app.Use(cors.New(cors.Config{
    AllowOrigins: []string{
        "http://localhost:3000",
        "http://127.0.0.1:3000",
        "http://localhost:4455", // Kratos UI
    },
    AllowMethods: []string{
        "GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
    },
    AllowHeaders: []string{
        "Origin",
        "Content-Type",
        "Accept",
        "Authorization",
        "X-Requested-With",
        "Cookie",
        "X-Session-Token", // For development Bearer token
    },
    ExposeHeaders: []string{
        "Content-Length",
        "Content-Type",
        "Set-Cookie",
    },
    AllowCredentials: true,
    MaxAge:           86400,
}))
```

---

## 🐳 5. Docker/Podman Setup

### **Development (podman-compose)**

```yaml
# docker-compose.yml
services:
  backend:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - DB_HOST=db
      - KRATOS_PUBLIC_URL=http://kratos:4433
    depends_on:
      - db
      - redis
      - kratos

  kratos:
    image: oryd/kratos:v1.1.0
    ports:
      - "4433:4433"  # Public API
      - "4434:4434"  # Admin API
    environment:
      - DSN=postgres://karima_store:lokal@db:5432/kratos?sslmode=disable

  kratos-ui:
    image: oryd/kratos-selfservice-ui-node:v1.1.0
    ports:
      - "4455:4455"  # Self-service UI
```

**Start Development:**
```bash
# Start all services
make kratos-up

# Or manually:
podman-compose up -d
```

### **Production (docker-compose.prod.yml)**

```yaml
# docker-compose.prod.yml
services:
  backend:
    image: karima_store_backend:${VERSION}
    restart: always
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=production
      - DB_HOST=postgres.ks-backend.cloud
      - KRATOS_PUBLIC_URL=https://auth.ks-backend.cloud
    env_file:
      - .env.production

  kratos:
    image: oryd/kratos:v1.1.0
    restart: always
    ports:
      - "4433:4433"
    environment:
      - DSN=postgres://user:pass@postgres:5432/kratos?sslmode=require
    volumes:
      - ./deploy/kratos/kratos.prod.yml:/etc/config/kratos/kratos.yml

  # No kratos-ui in production (use Cloudflare Pages frontend)
```

**Deploy Production:**
```bash
# Build for production
docker build -t karima_store_backend:v1.0.0 .

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

---

## 🧪 6. Testing in Development

### **Option 1: Using Kratos UI (Recommended for Development)**

```bash
# 1. Start services
make kratos-up

# 2. Open Kratos UI
http://127.0.0.1:4455

# 3. Register/Login
# Session cookie will be set for 127.0.0.1

# 4. Test API (Cookie approach won't work across localhost domains)
# Use Bearer token instead:

# Get session from Kratos
curl http://127.0.0.1:4433/sessions/whoami \
  -H "Cookie: ory_kratos_session=<your_session>"

# Use token with API
curl http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <session_token>"
```

### **Option 2: Using Bearer Token (For API Testing)**

```bash
# 1. Register user via Kratos API
curl -X POST http://127.0.0.1:4433/self-service/registration/api

# 2. Login and get session
curl -X POST http://127.0.0.1:4433/self-service/login/api \
  -H "Content-Type: application/json" \
  -d '{"identifier":"user@example.com","password":"secret"}'

# 3. Use session token
curl -X POST http://localhost:8080/api/v1/checkout \
  -H "Authorization: Bearer <session_token>" \
  -H "Content-Type: application/json" \
  -d '{...}'
```

---

## 🔒 7. Security Considerations

### **Development (Relaxed)**

✅ **OK in Development:**
- HTTP (not HTTPS)
- Weak passwords
- CORS allowing localhost
- Sandbox API keys
- Disabled SSL for database
- Verbose logging
- No rate limiting

### **Production (Strict)**

🔐 **REQUIRED in Production:**
- ✅ HTTPS only (TLS certificates)
- ✅ Strong passwords/secrets
- ✅ Strict CORS (only production domains)
- ✅ Production API keys
- ✅ Database SSL required
- ✅ Minimal logging (no sensitive data)
- ✅ Rate limiting enabled
- ✅ Cloudflare Zero Trust for admin panel
- ✅ Regular security updates
- ✅ Database backups
- ✅ Monitoring and alerts

---

## 📦 8. Deployment Checklist

### **Pre-Production Checklist**

- [ ] Update all environment variables to production values
- [ ] Change database credentials
- [ ] Enable SSL for database
- [ ] Update CORS origins to production domains
- [ ] Switch Midtrans to production keys
- [ ] Update R2 bucket to production
- [ ] Configure domain DNS
- [ ] Setup SSL certificates (Let's Encrypt)
- [ ] Enable Cloudflare Zero Trust for admin
- [ ] Configure Kratos for production domain
- [ ] Setup monitoring (Grafana/Prometheus)
- [ ] Configure log aggregation
- [ ] Setup automated backups
- [ ] Test all critical flows
- [ ] Load testing
- [ ] Security audit

---

## 🚀 9. Quick Start Commands

### **Development**

```bash
# Clone and setup
git clone <repo>
cd karima_store
cp .env.example .env

# Start all services
make kratos-up

# Or step by step:
podman-compose up -d db redis
podman-compose up -d kratos kratos-ui
podman-compose up -d backend

# View logs
podman-compose logs -f backend

# Rebuild after code changes
podman-compose build backend
podman-compose up -d backend

# Stop all
podman-compose down
```

### **Production**

```bash
# Pull latest code
git pull origin main

# Build production image
docker build -t karima_store_backend:v1.0.0 .

#Deploy with zero downtime
docker-compose -f docker-compose.prod.yml up -d --no-deps backend

# Health check
curl https://api.ks-backend.cloud/health
```

---

## 📝 10. Notes

**Current Development Environment:**
- ✅ OS: Windows with WSL2 (Ubuntu/Debian)
- ✅ Container Engine: **Podman** (not Docker Desktop)
- ✅ Orchestration: `podman-compose`
- ✅ **IMPORTANT**: All commands must use `podman` or `podman-compose`
- ✅ Podman runs rootless - ensure volume permissions are correct

**Production Environment:**
- VPS with Docker
- Cloudflare Pages for frontends
- Cloudflare Zero Trust for admin access
- Let's Encrypt for SSL
- PostgreSQL with SSL
- Redis with password

---

## 🆘 Troubleshooting

### **Cookie Issues in Development**

**Problem:** Session cookie not sent across `127.0.0.1:4455` → `localhost:8080`

**Solution:** Use Bearer token for development API testing:
```bash
# Get session token from Kratos
# Use it as: Authorization: Bearer <token>
```

### **CORS Errors**

**Problem:** CORS error when calling API from frontend

**Solution:** 
- Check `CORS_ORIGIN` in `.env`
- Verify frontend URL is allowed
- Ensure `AllowCredentials: true` for cookies

### **Container Build Issues**

**Problem:** Code changes not reflected

**Solution:**
```bash
# Force rebuild
podman-compose build --no-cache backend
podman-compose up -d backend
```

---

**Last Updated:** 2026-01-02
**Architecture Version:** 2.0 (with Ory Kratos)