# Ory Kratos Production Setup Guide

## Overview

This guide explains how to configure Ory Kratos for production deployment with security best practices. The configuration has been updated to address critical security vulnerabilities.

## Security Improvements Made

### 1. Database Configuration
**Before:** `dsn: memory` (in-memory storage, data lost on restart)
**After:** `dsn: ${KRATOS_DSN}` (PostgreSQL with environment variable)

### 2. Secret Management
**Before:** Hardcoded default secrets (CRITICAL SECURITY ISSUE)
```yaml
secrets:
  cookie:
    - PLEASE-CHANGE-ME-I-AM-VERY-INSECURE
  cipher:
    - 32-LONG-SECRET-MUST-BE-CHANGED
```

**After:** Environment variable-based secrets
```yaml
secrets:
  cookie:
    - ${KRATOS_SECRET_COOKIE}
  cipher:
    - ${KRATOS_SECRET_CIPHER}
```

### 3. CORS Configuration
**Before:** Wildcard origins (localhost only, but still insecure)
```yaml
cors:
  allowed_origins:
    - http://127.0.0.1:3000
    - http://localhost:3000
    - http://127.0.0.1:4455
    - http://localhost:4455
```

**After:** Production domains only
```yaml
cors:
  allowed_origins:
    - https://karima.com
    - https://admin.ks-backend.cloud
```

### 4. Logging Security
**Before:** Debug mode with sensitive value leakage
```yaml
log:
  level: debug
  leak_sensitive_values: true
```

**After:** Production-safe logging
```yaml
log:
  level: ${KRATOS_LOG_LEVEL:-info}
  format: ${KRATOS_LOG_FORMAT:-json}
  leak_sensitive_values: ${KRATOS_LEAK_SENSITIVE_VALUES:-false}
```

### 5. URL Configuration
**Before:** Localhost URLs only
```yaml
serve:
  public:
    base_url: http://127.0.0.1:4433/
selfservice:
  default_browser_return_url: http://127.0.0.1:4455/
```

**After:** Production URLs with HTTPS
```yaml
serve:
  public:
    base_url: ${KRATOS_PUBLIC_BASE_URL}
selfservice:
  default_browser_return_url: ${KRATOS_DEFAULT_RETURN_URL}
```

## Environment Variables Required

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `KRATOS_DSN` | PostgreSQL connection string | `postgresql://kratos:password@localhost:5432/kratos?sslmode=require` |
| `KRATOS_PUBLIC_BASE_URL` | Public API endpoint | `https://auth.ks-backend.cloud/` |
| `KRATOS_ADMIN_BASE_URL` | Admin API endpoint | `http://127.0.0.1:4434/` |
| `KRATOS_UI_URL` | Frontend application URL | `https://karima.com` |
| `KRATOS_DEFAULT_RETURN_URL` | Default return URL after auth | `https://karima.com/` |
| `KRATOS_SECRET_COOKIE` | Cookie signing secret (min 32 chars) | `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6` |
| `KRATOS_SECRET_CIPHER` | Cipher encryption secret (exactly 32 chars) | `a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6` |
| `KRATOS_SMTP_CONNECTION_URI` | SMTP server for emails | `smtps://user:pass@smtp.example.com:465/?skip_ssl_verify=false` |

### Optional Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `KRATOS_LOG_LEVEL` | Logging level | `info` |
| `KRATOS_LOG_FORMAT` | Log format (json/text) | `json` |
| `KRATOS_LEAK_SENSITIVE_VALUES` | Log sensitive data | `false` |

## Setup Instructions

### Step 1: Create Environment File

```bash
# Copy the example file
cp deploy/kratos/.env.kratos.example deploy/kratos/.env.kratos

# Edit the file with your actual values
nano deploy/kratos/.env.kratos
```

### Step 2: Generate Secure Secrets

Generate strong random secrets for production:

```bash
# Generate cookie secret (32+ characters)
openssl rand -base64 32

# Generate cipher secret (exactly 32 characters)
openssl rand -base64 32 | cut -c1-32
```

**IMPORTANT:** Store these secrets securely. Use a secrets manager (AWS Secrets Manager, HashiCorp Vault, etc.) in production.

### Step 3: Configure PostgreSQL

Create a dedicated database for Kratos:

```sql
-- Connect to PostgreSQL
psql -U postgres

-- Create database and user
CREATE DATABASE kratos;
CREATE USER kratos WITH PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE kratos TO kratos;
\q
```

Update `KRATOS_DSN` in your `.env.kratos` file with the correct connection string.

### Step 4: Configure SMTP

Set up email delivery for verification and recovery flows:

**For Gmail (using App Password):**
1. Enable 2FA on your Google account
2. Generate an App Password
3. Update `KRATOS_SMTP_CONNECTION_URI`:
   ```
   smtps://your-email@gmail.com:app-password@smtp.gmail.com:465/?skip_ssl_verify=false
   ```

**For SendGrid:**
```
smtps://apikey:SG.YOUR_API_KEY@smtp.sendgrid.net:465/?skip_ssl_verify=false
```

**For AWS SES:**
```
smtps://username:password@email-smtp.us-east-1.amazonaws.com:465/?skip_ssl_verify=false
```

### Step 5: Deploy with Docker/Podman

Update your `docker-compose.kratos.yml` to use environment variables:

```yaml
version: '3.8'

services:
  kratos:
    image: oryd/kratos:v0.10.1
    ports:
      - "4433:4433" # Public
      - "4434:4434" # Admin
    environment:
      - KRATOS_DSN=${KRATOS_DSN}
      - KRATOS_PUBLIC_BASE_URL=${KRATOS_PUBLIC_BASE_URL}
      - KRATOS_ADMIN_BASE_URL=${KRATOS_ADMIN_BASE_URL}
      - KRATOS_UI_URL=${KRATOS_UI_URL}
      - KRATOS_DEFAULT_RETURN_URL=${KRATOS_DEFAULT_RETURN_URL}
      - KRATOS_SECRET_COOKIE=${KRATOS_SECRET_COOKIE}
      - KRATOS_SECRET_CIPHER=${KRATOS_SECRET_CIPHER}
      - KRATOS_LOG_LEVEL=${KRATOS_LOG_LEVEL:-info}
      - KRATOS_LOG_FORMAT=${KRATOS_LOG_FORMAT:-json}
      - KRATOS_LEAK_SENSITIVE_VALUES=${KRATOS_LEAK_SENSITIVE_VALUES:-false}
      - KRATOS_SMTP_CONNECTION_URI=${KRATOS_SMTP_CONNECTION_URI}
    volumes:
      - ./deploy/kratos/kratos.yml:/etc/config/kratos/kratos.yml
      - ./deploy/kratos/identity.schema.json:/etc/config/kratos/identity.schema.json
    depends_on:
      - postgres
    restart: unless-stopped

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=kratos
      - POSTGRES_USER=kratos
      - POSTGRES_PASSWORD=${KRATOS_DB_PASSWORD}
    volumes:
      - kratos_data:/var/lib/postgresql/data
    restart: unless-stopped

volumes:
  kratos_data:
```

Start the services:

```bash
# Load environment variables
export $(cat deploy/kratos/.env.kratos | xargs)

# Start services
docker-compose -f docker-compose.kratos.yml up -d
```

### Step 6: Verify Deployment

Check that Kratos is running correctly:

```bash
# Check health endpoint
curl http://localhost:4434/health/alive

# Check public endpoint
curl https://auth.ks-backend.cloud/.well-known/ory/kratos/public
```

## Security Best Practices

### 1. HTTPS/SSL Configuration

Ensure your application is behind HTTPS (Cloudflare, Nginx, etc.). The cookie security settings are automatically configured by Ory Kratos when using HTTPS URLs:

- `secure: true` - Cookies only sent over HTTPS
- `http_only: true` - Cookies not accessible via JavaScript
- `same_site: strict` - CSRF protection

### 2. Secrets Management

**Never commit secrets to version control!** Use one of these approaches:

1. **Environment Variables** (simple, suitable for single-server deployments)
2. **Secrets Manager** (AWS Secrets Manager, HashiCorp Vault, Azure Key Vault)
3. **Kubernetes Secrets** (if using K8s)
4. **Docker Secrets** (if using Swarm)

### 3. Database Security

- Use `sslmode=require` in DSN for encrypted connections
- Create dedicated database user with minimal privileges
- Enable PostgreSQL connection encryption
- Regular database backups

### 4. Network Security

- Admin API (`KRATOS_ADMIN_BASE_URL`) should be internal-only
- Public API (`KRATOS_PUBLIC_BASE_URL`) should be behind Cloudflare/WAF
- Restrict database access to Kratos container only
- Use firewall rules to limit access

### 5. Monitoring & Logging

- Set `KRATOS_LOG_LEVEL=info` or `warning` in production
- Use `KRATOS_LOG_FORMAT=json` for log aggregation
- Monitor authentication failures and suspicious activity
- Set up alerts for security events

## Troubleshooting

### Issue: Database Connection Failed

**Solution:** Verify DSN format and PostgreSQL accessibility:
```bash
# Test connection
psql "${KRATOS_DSN}"
```

### Issue: SMTP Not Working

**Solution:** Test SMTP connection:
```bash
# Test with telnet
telnet smtp.gmail.com 465

# Check logs
docker logs kratos
```

### Issue: CORS Errors

**Solution:** Ensure `KRATOS_UI_URL` and allowed origins match your frontend domain exactly.

### Issue: Cookie Not Set

**Solution:** Verify:
1. Using HTTPS in production
2. Domain names match exactly
3. Secrets are set correctly
4. Browser is not blocking cookies

## Migration from Development to Production

1. **Backup existing data** (if any)
2. **Create production database**
3. **Set environment variables**
4. **Deploy with new configuration**
5. **Test authentication flows**
6. **Monitor logs and metrics**
7. **Update DNS to point to new instance**

## Additional Resources

- [Ory Kratos Documentation](https://www.ory.sh/docs/kratos)
- [Ory Kratos Security Best Practices](https://www.ory.sh/docs/kratos/self-service/flows/user-login)
- [Ory Kratos Configuration Reference](https://www.ory.sh/docs/kratos/reference/configuration)

## Checklist Before Going Live

- [ ] All environment variables set correctly
- [ ] Strong secrets generated and stored securely
- [ ] PostgreSQL database created and accessible
- [ ] SMTP server configured and tested
- [ ] HTTPS/SSL enabled
- [ ] CORS origins set to production domains
- [ ] Logging configured to `info` level
- [ ] `leak_sensitive_values` set to `false`
- [ ] Admin API not publicly accessible
- [ ] Database backups configured
- [ ] Monitoring and alerting set up
- [ ] Security audit completed
- [ ] Load testing performed
- [ ] Disaster recovery plan documented

---

**Last Updated:** 2026-01-02
**Version:** 1.0.0
