# 🚀 Production Setup - Quick Start (5 Menit)

## Step 1: Generate Passwords
```bash
make generate-secrets
```
**📝 Simpan output ke password manager!**

## Step 2: Buat File Production
```bash
make setup-prod-env
```

## Step 3: Edit .env.production

Buka file dan isi nilai berikut:

```bash
# Database
DB_HOST=192.168.1.100                    # IP VPS Anda
DB_PASSWORD=<dari step 1>

# Redis  
REDIS_HOST=192.168.1.100                 # IP VPS Anda
REDIS_PASSWORD=<dari step 1>

# JWT
JWT_SECRET=<dari step 1>

# Midtrans (dari dashboard.midtrans.com)
MIDTRANS_SERVER_KEY=<production key>
MIDTRANS_CLIENT_KEY=<production key>

# Cloudflare R2 (dari dash.cloudflare.com)
R2_ACCOUNT_ID=<account id>
R2_ACCESS_KEY_ID=<access key>
R2_SECRET_ACCESS_KEY=<secret key>

# Email
EMAIL_USER=noreply@yourdomain.com
EMAIL_PASSWORD=<gmail app password>

# CORS
CORS_ORIGIN=https://yourdomain.com
```

## Step 4: Verify
```bash
make verify-env
```

## Step 5: Deploy ke VPS

### Setup VPS (sekali saja):
```bash
scp scripts/setup-vps.sh user@vps:~/
ssh user@vps "chmod +x setup-vps.sh && sudo ./setup-vps.sh"
```

### Deploy Aplikasi:
```bash
# Build
GOOS=linux GOARCH=amd64 go build -o karima-store ./cmd/api

# Copy files
scp karima-store .env.production user@vps:/opt/karima-store/
scp -r migrations user@vps:/opt/karima-store/

# Setup service
scp deploy/systemd/karima-store.service user@vps:~/
ssh user@vps "sudo mv karima-store.service /etc/systemd/system/ && \
              sudo systemctl enable karima-store && \
              sudo systemctl start karima-store"

# Setup Nginx
scp deploy/nginx/karima-store.conf user@vps:~/
ssh user@vps "sudo mv karima-store.conf /etc/nginx/sites-available/karima-store && \
              sudo ln -sf /etc/nginx/sites-available/karima-store /etc/nginx/sites-enabled/ && \
              sudo nginx -t && sudo systemctl reload nginx"
```

## Step 6: Verify Deployment
```bash
# Check service
ssh user@vps "sudo systemctl status karima-store"

# Test API
curl https://api.yourdomain.com/health
```

## ✅ Done!

**📚 Dokumentasi Lengkap**:
- Setup Guide: `docs/deployment/PRODUCTION_ENV_SETUP.md`
- Checklist: `docs/deployment/DEPLOYMENT_CHECKLIST.md`
- Summary: `docs/deployment/PRODUCTION_SETUP_SUMMARY.md`

**🆘 Troubleshooting**:
```bash
# View logs
ssh user@vps "sudo journalctl -u karima-store -f"

# Verify config
make verify-env
```
