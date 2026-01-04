# Production Deployment Checklist

## 📋 Pre-Deployment Checklist

### 1. Environment Configuration ✓
- [ ] Generate secure passwords: `make generate-secrets`
- [ ] Create `.env.production`: `make setup-prod-env`
- [ ] Fill all required values in `.env.production`
- [ ] Verify configuration: `make verify-env`
- [ ] Set file permissions: `chmod 600 .env.production`
- [ ] Store secrets in password manager (1Password, Vault, etc.)

### 2. Database Setup ✓
- [ ] PostgreSQL installed and running on VPS
- [ ] Create production database: `karima_db_prod`
- [ ] Create dedicated user: `karima_prod_user`
- [ ] Grant appropriate permissions (not superuser!)
- [ ] Enable SSL/TLS: `ssl = on` in postgresql.conf
- [ ] Configure `pg_hba.conf` for SSL connections
- [ ] Test connection: `psql "postgresql://user:pass@host:5432/db?sslmode=require"`
- [ ] Setup automated backups (pg_dump + cron)
- [ ] Configure backup retention policy

### 3. Redis Setup ✓
- [ ] Redis installed and running on VPS
- [ ] Enable password authentication in redis.conf
- [ ] Bind to internal IP only (not 0.0.0.0)
- [ ] Enable persistence: RDB + AOF
- [ ] Configure maxmemory policy
- [ ] Test connection: `redis-cli -h host -p 6379 -a password ping`
- [ ] Setup monitoring (memory usage, hit rate)

### 4. Ory Kratos Setup ✓
- [ ] Kratos deployed and accessible
- [ ] SSL certificate installed (Let's Encrypt)
- [ ] Domain configured: `auth.yourdomain.com`
- [ ] SMTP configured for email verification
- [ ] Identity schemas configured
- [ ] Session management configured
- [ ] Cookie settings configured
- [ ] Test registration flow
- [ ] Test login flow
- [ ] Test password recovery

### 5. External Services ✓

#### Midtrans (Payment Gateway)
- [ ] Production account verified
- [ ] Business verification completed
- [ ] Production Server Key obtained
- [ ] Production Client Key obtained
- [ ] Notification URL configured
- [ ] Test payment flow in sandbox
- [ ] Switch to production mode
- [ ] Test small production payment

#### RajaOngkir (Shipping)
- [ ] Upgraded to Pro/Business plan
- [ ] Production API key obtained
- [ ] Test shipping calculation
- [ ] Verify supported couriers
- [ ] Configure fallback options

#### Cloudflare R2 (Storage)
- [ ] Bucket created: `karima-media-prod`
- [ ] API token generated with R/W permissions
- [ ] Custom domain configured: `media.yourdomain.com`
- [ ] CORS configured for web access
- [ ] Lifecycle policies configured
- [ ] Test file upload
- [ ] Test file download via public URL

#### Email (SMTP)
- [ ] Production email account setup
- [ ] 2FA enabled (if Gmail)
- [ ] App password generated
- [ ] SPF record configured
- [ ] DKIM configured
- [ ] DMARC configured
- [ ] Test email delivery
- [ ] Check spam score

#### Fonnte (WhatsApp)
- [ ] Business WhatsApp verified
- [ ] API token obtained
- [ ] Message templates configured
- [ ] Test message sending
- [ ] Webhook configured (if needed)

### 6. Application Code ✓
- [ ] All tests passing: `make test`
- [ ] Test coverage ≥ 80%: `make test-coverage`
- [ ] No critical security vulnerabilities
- [ ] Code reviewed and approved
- [ ] Version tagged in Git
- [ ] Changelog updated
- [ ] API documentation updated: `make swagger`

### 7. Infrastructure ✓
- [ ] VPS provisioned and accessible
- [ ] Firewall configured (UFW/iptables)
  - [ ] Port 22 (SSH) - restricted IPs
  - [ ] Port 80 (HTTP) - open
  - [ ] Port 443 (HTTPS) - open
  - [ ] Port 5432 (PostgreSQL) - internal only
  - [ ] Port 6379 (Redis) - internal only
- [ ] SSL certificates installed
- [ ] Nginx/Caddy configured as reverse proxy
- [ ] systemd service configured
- [ ] Log rotation configured
- [ ] Disk space monitored
- [ ] Backup storage configured

### 8. Security ✓
- [ ] SSH key-based authentication only
- [ ] Root login disabled
- [ ] Fail2ban installed and configured
- [ ] Security updates enabled
- [ ] File permissions reviewed
- [ ] Secrets not in Git history
- [ ] API rate limiting configured
- [ ] CORS properly restricted
- [ ] Security headers configured
- [ ] SQL injection protection verified
- [ ] XSS protection verified

### 9. Monitoring & Logging ✓
- [ ] Application logging configured
- [ ] Log level set to `warn` or `error`
- [ ] Log rotation configured
- [ ] Centralized logging setup (optional)
- [ ] Error tracking setup (Sentry, etc.)
- [ ] Uptime monitoring (UptimeRobot, etc.)
- [ ] Performance monitoring
- [ ] Alerting configured
- [ ] Dashboard created

### 10. Documentation ✓
- [ ] Deployment runbook created
- [ ] Architecture diagram updated
- [ ] API documentation published
- [ ] Environment variables documented
- [ ] Troubleshooting guide created
- [ ] Incident response plan documented
- [ ] Rollback procedure documented
- [ ] Team access documented

---

## 🚀 Deployment Steps

### Step 1: Prepare Application
```bash
# On local machine
cd /home/xnfo/projects/karima_store

# Run tests
make test

# Generate production env
make setup-prod-env

# Edit .env.production with real values
nano .env.production

# Verify configuration
make verify-env

# Build binary (if deploying binary)
GOOS=linux GOARCH=amd64 go build -o karima-store ./cmd/api
```

### Step 2: Prepare VPS
```bash
# SSH to VPS
ssh user@your-vps-ip

# Create application directory
sudo mkdir -p /opt/karima-store
sudo chown $USER:$USER /opt/karima-store

# Create logs directory
mkdir -p /opt/karima-store/logs

# Install dependencies
sudo apt update
sudo apt install -y postgresql-client redis-tools curl
```

### Step 3: Deploy Application
```bash
# On local machine - Copy files
scp .env.production user@vps:/opt/karima-store/.env.production
scp karima-store user@vps:/opt/karima-store/
scp -r migrations user@vps:/opt/karima-store/

# Or use Git (recommended)
# On VPS
cd /opt/karima-store
git clone https://github.com/yourusername/karima-store.git .
git checkout production  # or specific tag

# Copy environment file
scp .env.production user@vps:/opt/karima-store/
```

### Step 4: Configure Service
```bash
# On VPS - Create systemd service
sudo nano /etc/systemd/system/karima-store.service
```

```ini
[Unit]
Description=Karima Store API
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=karima
Group=karima
WorkingDirectory=/opt/karima-store
EnvironmentFile=/opt/karima-store/.env.production
ExecStart=/opt/karima-store/karima-store
Restart=always
RestartSec=10
StandardOutput=append:/opt/karima-store/logs/app.log
StandardError=append:/opt/karima-store/logs/error.log

# Security
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/karima-store/logs

[Install]
WantedBy=multi-user.target
```

```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable service
sudo systemctl enable karima-store

# Start service
sudo systemctl start karima-store

# Check status
sudo systemctl status karima-store
```

### Step 5: Configure Nginx
```bash
# Create Nginx config
sudo nano /etc/nginx/sites-available/karima-store
```

```nginx
upstream karima_backend {
    server 127.0.0.1:8080;
}

server {
    listen 80;
    server_name api.yourdomain.com;
    
    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.yourdomain.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/api.yourdomain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/api.yourdomain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security Headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Logging
    access_log /var/log/nginx/karima-access.log;
    error_log /var/log/nginx/karima-error.log;

    # Client body size
    client_max_body_size 10M;

    location / {
        proxy_pass http://karima_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # Health check endpoint
    location /health {
        proxy_pass http://karima_backend/health;
        access_log off;
    }
}
```

```bash
# Enable site
sudo ln -s /etc/nginx/sites-available/karima-store /etc/nginx/sites-enabled/

# Test configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### Step 6: Run Migrations
```bash
# On VPS
cd /opt/karima-store

# Run migrations
./karima-store migrate up

# Or if using separate migration tool
migrate -path ./migrations -database "postgresql://user:pass@localhost:5432/karima_db_prod?sslmode=require" up
```

### Step 7: Verify Deployment
```bash
# Check service status
sudo systemctl status karima-store

# Check logs
tail -f /opt/karima-store/logs/app.log

# Test health endpoint
curl https://api.yourdomain.com/health

# Test API endpoint
curl https://api.yourdomain.com/api/v1/products
```

---

## ✅ Post-Deployment Checklist

### Immediate (Within 1 hour)
- [ ] Application is running
- [ ] Health check endpoint responding
- [ ] Database connection working
- [ ] Redis connection working
- [ ] All API endpoints accessible
- [ ] SSL certificate valid
- [ ] Logs are being written
- [ ] No error spikes in logs

### Critical Flows (Within 4 hours)
- [ ] User registration working
- [ ] User login working
- [ ] Password reset working
- [ ] Product listing working
- [ ] Product search working
- [ ] Cart operations working
- [ ] Checkout flow working
- [ ] Payment processing working
- [ ] Email notifications working
- [ ] WhatsApp notifications working
- [ ] File upload working

### Monitoring (Within 24 hours)
- [ ] Uptime monitoring active
- [ ] Error tracking active
- [ ] Performance metrics collected
- [ ] Alerts configured
- [ ] Dashboard accessible
- [ ] Backup job ran successfully
- [ ] Log rotation working

### Documentation (Within 48 hours)
- [ ] Deployment documented
- [ ] Issues encountered documented
- [ ] Team notified
- [ ] Runbook updated
- [ ] Changelog published

---

## 🔄 Rollback Procedure

If deployment fails:

```bash
# Stop new version
sudo systemctl stop karima-store

# Restore previous version
cd /opt/karima-store
git checkout previous-tag

# Or restore binary
cp karima-store.backup karima-store

# Rollback migrations (if needed)
migrate -path ./migrations -database "postgresql://..." down 1

# Start service
sudo systemctl start karima-store

# Verify
curl https://api.yourdomain.com/health
```

---

## 📊 Success Metrics

After deployment, monitor:
- **Uptime**: Should be > 99.9%
- **Response Time**: p95 < 500ms
- **Error Rate**: < 0.1%
- **Database Connections**: Stable
- **Memory Usage**: < 80%
- **CPU Usage**: < 70%
- **Disk Usage**: < 80%

---

## 🆘 Emergency Contacts

- **DevOps Lead**: [Contact]
- **Backend Lead**: [Contact]
- **On-Call Engineer**: [Contact]
- **Database Admin**: [Contact]

---

## 📚 Related Documentation

- [Production Environment Setup](./PRODUCTION_ENV_SETUP.md)
- [Quick Reference](./PRODUCTION_ENV_QUICK_REFERENCE.md)
- [Troubleshooting Guide](./TROUBLESHOOTING.md)
- [Security Best Practices](../security/SECURITY_BEST_PRACTICES.md)

---

**Last Updated**: 2026-01-04  
**Version**: 1.0.0  
**Deployment Type**: VPS Production
