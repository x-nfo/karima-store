# Deployment Configuration Files

This directory contains configuration files and templates for deploying Karima Store to production.

## 📁 Directory Structure

```
deploy/
├── nginx/              # Nginx reverse proxy configuration
│   └── karima-store.conf
├── systemd/            # Systemd service files
│   └── karima-store.service
└── kratos/             # Ory Kratos authentication configuration
    └── kratos.yml
```

## 🚀 Quick Start

### 1. Setup VPS
```bash
# Copy setup script to VPS
scp scripts/setup-vps.sh user@your-vps:~/

# SSH to VPS and run
ssh user@your-vps
chmod +x setup-vps.sh
sudo ./setup-vps.sh
```

### 2. Deploy Application
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o karima-store ./cmd/api

# Copy files to VPS
scp karima-store user@vps:/opt/karima-store/
scp .env.production user@vps:/opt/karima-store/
scp -r migrations user@vps:/opt/karima-store/

# Setup systemd service
scp deploy/systemd/karima-store.service user@vps:~/
ssh user@vps "sudo mv karima-store.service /etc/systemd/system/ && \
              sudo systemctl daemon-reload && \
              sudo systemctl enable karima-store && \
              sudo systemctl start karima-store"
```

### 3. Configure Nginx
```bash
# Copy Nginx config
scp deploy/nginx/karima-store.conf user@vps:~/

# Install on VPS
ssh user@vps "sudo mv karima-store.conf /etc/nginx/sites-available/ && \
              sudo ln -sf /etc/nginx/sites-available/karima-store.conf /etc/nginx/sites-enabled/ && \
              sudo nginx -t && \
              sudo systemctl reload nginx"
```

## 📋 Configuration Files

### Nginx Configuration
**File**: `nginx/karima-store.conf`

Features:
- SSL/TLS termination
- HTTP/2 support
- Security headers (HSTS, CSP, etc.)
- Rate limiting
- Gzip compression
- Reverse proxy to backend
- Static file serving
- Health check endpoint

**Customization Required**:
- Replace `api.yourdomain.com` with your actual domain
- Update SSL certificate paths
- Adjust rate limiting based on your traffic
- Configure IP whitelisting for sensitive endpoints

### Systemd Service
**File**: `systemd/karima-store.service`

Features:
- Automatic restart on failure
- Security hardening (NoNewPrivileges, ProtectSystem, etc.)
- Resource limits (CPU, Memory)
- Proper logging
- Environment file support

**Customization Required**:
- Verify `WorkingDirectory` path
- Adjust resource limits based on your VPS specs
- Update `User` and `Group` if different

### Ory Kratos
**File**: `kratos/kratos.yml`

Configuration for authentication service.

## 🔧 Customization Guide

### Update Domain Names

1. **Nginx Configuration**:
```bash
# Replace in nginx/karima-store.conf
sed -i 's/api.yourdomain.com/api.example.com/g' deploy/nginx/karima-store.conf
```

2. **Environment Variables**:
```bash
# Update in .env.production
KRATOS_PUBLIC_URL=https://auth.example.com
CORS_ORIGIN=https://example.com,https://www.example.com
```

### Adjust Resource Limits

Edit `systemd/karima-store.service`:
```ini
# For 2GB RAM VPS
MemoryMax=1G
MemoryHigh=800M
CPUQuota=100%

# For 4GB RAM VPS
MemoryMax=2G
MemoryHigh=1.5G
CPUQuota=200%
```

### Configure Rate Limiting

Edit `nginx/karima-store.conf`:
```nginx
# Adjust based on expected traffic
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=100r/m;  # 100 requests per minute
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=20r/m;  # 20 requests per minute
```

## 🔒 Security Checklist

- [ ] SSL certificates installed and valid
- [ ] Security headers configured in Nginx
- [ ] Rate limiting enabled
- [ ] Firewall configured (UFW/iptables)
- [ ] SSH key-based authentication only
- [ ] Database SSL mode enabled
- [ ] Systemd security hardening enabled
- [ ] File permissions set correctly (600 for .env)
- [ ] Sensitive endpoints IP-restricted
- [ ] CORS properly configured

## 📊 Monitoring

### Check Service Status
```bash
# Application status
sudo systemctl status karima-store

# View logs
sudo journalctl -u karima-store -f

# Nginx status
sudo systemctl status nginx

# View Nginx logs
sudo tail -f /var/log/nginx/karima-access.log
sudo tail -f /var/log/nginx/karima-error.log
```

### Health Checks
```bash
# Application health
curl https://api.yourdomain.com/health

# SSL certificate expiry
sudo certbot certificates

# Disk space
df -h

# Memory usage
free -h

# Service resource usage
systemctl status karima-store
```

## 🔄 Updates and Rollbacks

### Deploy Update
```bash
# Build new version
GOOS=linux GOARCH=amd64 go build -o karima-store ./cmd/api

# Backup current version
ssh user@vps "cp /opt/karima-store/karima-store /opt/karima-store/karima-store.backup"

# Deploy new version
scp karima-store user@vps:/opt/karima-store/

# Restart service
ssh user@vps "sudo systemctl restart karima-store"

# Check status
ssh user@vps "sudo systemctl status karima-store"
```

### Rollback
```bash
# Restore previous version
ssh user@vps "cp /opt/karima-store/karima-store.backup /opt/karima-store/karima-store && \
              sudo systemctl restart karima-store"
```

## 🆘 Troubleshooting

### Service Won't Start
```bash
# Check logs
sudo journalctl -u karima-store -n 50 --no-pager

# Check configuration
sudo systemctl cat karima-store

# Verify binary
ls -la /opt/karima-store/karima-store
file /opt/karima-store/karima-store

# Check permissions
sudo -u karima /opt/karima-store/karima-store --version
```

### Nginx Issues
```bash
# Test configuration
sudo nginx -t

# Check error logs
sudo tail -f /var/log/nginx/error.log

# Verify upstream
curl http://127.0.0.1:8080/health
```

### SSL Certificate Issues
```bash
# Check certificate
sudo certbot certificates

# Renew certificate
sudo certbot renew

# Test renewal
sudo certbot renew --dry-run
```

## 📚 Related Documentation

- [Production Environment Setup](../docs/deployment/PRODUCTION_ENV_SETUP.md)
- [Deployment Checklist](../docs/deployment/DEPLOYMENT_CHECKLIST.md)
- [Quick Reference](../docs/deployment/PRODUCTION_ENV_QUICK_REFERENCE.md)

## 🤝 Support

For deployment issues:
1. Check service logs: `sudo journalctl -u karima-store -f`
2. Review Nginx logs: `sudo tail -f /var/log/nginx/karima-error.log`
3. Verify environment: `./scripts/verify-env.sh .env.production`
4. Consult troubleshooting guide
5. Contact DevOps team

---

**Last Updated**: 2026-01-04  
**Maintained by**: DevOps Team
