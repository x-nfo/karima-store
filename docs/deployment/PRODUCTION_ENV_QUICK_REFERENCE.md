# Production Environment - Quick Reference

## 🚀 Quick Setup (5 Minutes)

### 1. Generate Secrets
```bash
chmod +x scripts/generate-env-secrets.sh
./scripts/generate-env-secrets.sh
```

### 2. Copy Template
```bash
cp .env.production.template .env.production
```

### 3. Fill Required Values

Open `.env.production` and replace these **CRITICAL** values:

#### Database (PostgreSQL)
```bash
DB_HOST=192.168.1.100              # Your VPS IP or domain
DB_USER=karima_prod_user
DB_PASSWORD=<from generator>        # 32+ characters
DB_NAME=karima_db_prod
```

#### Redis
```bash
REDIS_HOST=192.168.1.100           # Your VPS IP or domain
REDIS_PASSWORD=<from generator>     # 32+ characters
```

#### JWT
```bash
JWT_SECRET=<from generator>         # 64 characters hex
```

#### Midtrans (Payment)
```bash
MIDTRANS_SERVER_KEY=<from dashboard>
MIDTRANS_CLIENT_KEY=<from dashboard>
MIDTRANS_IS_PRODUCTION=true
```

#### Cloudflare R2 (Storage)
```bash
R2_ACCOUNT_ID=<from cloudflare>
R2_ACCESS_KEY_ID=<from cloudflare>
R2_SECRET_ACCESS_KEY=<from cloudflare>
R2_BUCKET_NAME=karima-media-prod
R2_PUBLIC_URL=https://media.yourdomain.com
```

#### Email
```bash
EMAIL_USER=noreply@yourdomain.com
EMAIL_PASSWORD=<gmail app password>
```

#### CORS
```bash
CORS_ORIGIN=https://yourdomain.com,https://www.yourdomain.com
```

### 4. Verify Configuration
```bash
chmod +x scripts/verify-env.sh
./scripts/verify-env.sh .env.production
```

### 5. Deploy
```bash
# Set secure permissions
chmod 600 .env.production

# Copy to server
scp .env.production user@server:/app/.env.production

# SSH to server and verify
ssh user@server
cd /app
chmod 600 .env.production
```

## 📝 Where to Get Production Keys

### Midtrans
1. Login: https://dashboard.midtrans.com
2. Switch to **Production** mode (top right)
3. Go to: Settings → Access Keys
4. Copy: Server Key & Client Key

### RajaOngkir
1. Login: https://rajaongkir.com
2. Upgrade to Pro/Business plan
3. Go to: API Key section
4. Copy production API key

### Cloudflare R2
1. Login: https://dash.cloudflare.com
2. Navigate to: R2 → Create Bucket
3. Create: `karima-media-prod`
4. Go to: Manage R2 API Tokens
5. Create token with Read & Write permissions
6. Copy: Account ID, Access Key, Secret Key

### Gmail App Password
1. Enable 2FA: https://myaccount.google.com/security
2. Go to: App passwords
3. Select: Mail → Other (Custom name)
4. Generate and copy password

### Fonnte (WhatsApp)
1. Login: https://fonnte.com
2. Go to: Dashboard
3. Copy: API Token
4. Verify WhatsApp number

## ⚡ Common Commands

```bash
# Generate secrets
make generate-secrets

# Verify environment
make verify-env

# Deploy to production
make deploy-prod

# Check production status
make prod-status

# View production logs
make prod-logs

# Rollback deployment
make prod-rollback
```

## 🔒 Security Checklist

- [ ] All passwords are 32+ characters
- [ ] JWT secret is 64+ characters
- [ ] File permissions: `chmod 600 .env.production`
- [ ] SSL enabled: `DB_SSL_MODE=require`
- [ ] CORS restricted to production domains
- [ ] Log level: `warn` or `error`
- [ ] Midtrans production mode enabled
- [ ] No localhost in CORS_ORIGIN
- [ ] Secrets stored in password manager

## 🆘 Troubleshooting

### Database Connection Failed
```bash
# Test connection
psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Check SSL
psql "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require"
```

### Redis Connection Failed
```bash
# Test connection
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD ping
```

### Midtrans Not Working
- Verify you're using **Production** keys (not Sandbox)
- Check `MIDTRANS_IS_PRODUCTION=true`
- Test API: https://dashboard.midtrans.com/settings/config_info

### Email Not Sending
- Verify Gmail App Password (not account password)
- Check 2FA is enabled
- Test SMTP: `telnet smtp.gmail.com 587`

## 📚 Full Documentation

- Complete Setup: `docs/deployment/PRODUCTION_ENV_SETUP.md`
- Security Guide: `docs/security/SECURITY_BEST_PRACTICES.md`
- Deployment Guide: `docs/deployment/DEPLOYMENT_GUIDE.md`

## 🔄 Secret Rotation (Every 90 Days)

```bash
# 1. Generate new secrets
./scripts/generate-env-secrets.sh

# 2. Update .env.production with new values

# 3. Verify configuration
./scripts/verify-env.sh .env.production

# 4. Deploy with zero downtime
make deploy-prod

# 5. Verify application working
make prod-health-check

# 6. Revoke old secrets
```

## 📞 Support

- Documentation: `docs/`
- Issues: Check logs first
- Emergency: Contact DevOps team

---
**Last Updated**: 2026-01-04
