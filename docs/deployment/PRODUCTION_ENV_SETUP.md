# Production Environment Setup Guide

## 📋 Overview

Panduan lengkap untuk mengkonfigurasi environment variables production untuk Karima Store.

## 🔐 Security First

**PENTING**: Jangan pernah commit file `.env.production` yang sudah terisi ke Git!

### Checklist Keamanan
- [ ] Gunakan password yang kuat dan unik untuk setiap service
- [ ] Aktifkan SSL/TLS untuk semua koneksi eksternal
- [ ] Simpan secrets di password manager atau secret management tool
- [ ] Rotate secrets secara berkala (90 hari)
- [ ] Gunakan environment variables injection, bukan hardcode
- [ ] Audit akses ke production secrets

## 🛠️ Step-by-Step Setup

### 1. Generate Secure Passwords

Jalankan script generator:

```bash
chmod +x scripts/generate-env-secrets.sh
./scripts/generate-env-secrets.sh
```

Script ini akan generate:
- Database password (32 karakter)
- Redis password (32 karakter)
- JWT secret (64 karakter hex)
- API keys
- Session secrets
- Encryption keys

### 2. Konfigurasi Database (PostgreSQL)

```bash
# Database Configuration
DB_HOST=your_vps_ip_or_domain          # Contoh: 192.168.1.100 atau db.ks-backend.cloud
DB_PORT=5432
DB_USER=karima_prod_user               # Gunakan user khusus, bukan 'postgres'
DB_PASSWORD=<generated_password>       # Dari script generator
DB_NAME=karima_db_prod
DB_SSL_MODE=require                    # WAJIB untuk production!
```

**Tips Database:**
- Buat user database khusus dengan privilege terbatas
- Jangan gunakan user `postgres` atau `root`
- Aktifkan SSL/TLS connection
- Setup database backup otomatis
- Monitor database performance

### 3. Konfigurasi Redis

```bash
# Redis Configuration
REDIS_HOST=your_redis_host             # Contoh: 192.168.1.100 atau redis.ks-backend.cloud
REDIS_PORT=6379
REDIS_PASSWORD=<generated_password>    # Dari script generator
```

**Tips Redis:**
- Aktifkan password authentication
- Bind ke IP internal saja (jangan 0.0.0.0)
- Gunakan Redis ACL untuk granular permissions
- Setup persistence (RDB + AOF)
- Monitor memory usage

### 4. Konfigurasi Ory Kratos

```bash
# Ory Kratos Configuration
KRATOS_PUBLIC_URL=https://auth.yourdomain.com
KRATOS_ADMIN_URL=https://auth.yourdomain.com/admin
KRATOS_UI_URL=https://auth.yourdomain.com/ui
```

**Tips Kratos:**
- Gunakan subdomain khusus untuk auth (contoh: auth.ks-backend.cloud)
- Setup SSL certificate (Let's Encrypt)
- Konfigurasi SMTP untuk email verification
- Setup session management yang aman
- Review Kratos security checklist

### 5. Konfigurasi JWT

```bash
# JWT Configuration
JWT_SECRET=<generated_hex_64_chars>    # Dari script generator
JWT_EXPIRATION=24h                     # Sesuaikan dengan kebutuhan
```

**Tips JWT:**
- Gunakan minimum 64 karakter untuk secret
- Set expiration sesuai security policy (1h - 24h)
- Implement refresh token mechanism
- Rotate JWT secret secara berkala
- Gunakan strong signing algorithm (HS256 minimum, RS256 recommended)

### 6. Konfigurasi Payment Gateway (Midtrans)

```bash
# Midtrans Configuration
MIDTRANS_SERVER_KEY=<your_production_server_key>
MIDTRANS_CLIENT_KEY=<your_production_client_key>
MIDTRANS_IS_PRODUCTION=true
MIDTRANS_API_BASE_URL=https://app.midtrans.com/snap/v1
```

**Cara Mendapatkan Production Keys:**
1. Login ke [Midtrans Dashboard](https://dashboard.midtrans.com)
2. Pilih environment "Production"
3. Navigate ke Settings → Access Keys
4. Copy Server Key dan Client Key
5. Lengkapi verifikasi bisnis jika belum

**Tips Midtrans:**
- Jangan gunakan sandbox keys di production!
- Setup notification URL untuk payment callback
- Implement proper error handling
- Log semua transaksi
- Setup monitoring untuk failed payments

### 7. Konfigurasi Shipping (RajaOngkir/Komerce)

```bash
# Shipping Configuration
RAJAONGKIR_API_KEY=<your_production_api_key>
RAJAONKIR_API_KEY_SHIPPING_DELIVERY=<your_production_shipping_key>
RAJAONGKIR_BASE_URL=https://pro.rajaongkir.com/api
```

**Cara Mendapatkan Production Keys:**
1. Login ke [RajaOngkir](https://rajaongkir.com)
2. Upgrade ke paket Pro atau Business
3. Navigate ke API Key section
4. Generate production API key

**Tips RajaOngkir:**
- Gunakan Pro/Business account untuk production
- Cache hasil shipping calculation
- Implement retry mechanism
- Monitor API quota usage
- Setup fallback shipping options

### 8. Konfigurasi Storage (Cloudflare R2)

```bash
# Cloudflare R2 Configuration
FILE_STORAGE=r2
FILE_UPLOAD_MAX_SIZE=10MB
R2_ACCOUNT_ID=<your_cloudflare_account_id>
R2_ENDPOINT=https://<account_id>.r2.cloudflarestorage.com
R2_ACCESS_KEY_ID=<your_r2_access_key>
R2_SECRET_ACCESS_KEY=<your_r2_secret_key>
R2_BUCKET_NAME=karima-media-prod
R2_PUBLIC_URL=https://media.yourdomain.com
R2_REGION=auto
```

**Cara Setup Cloudflare R2:**
1. Login ke [Cloudflare Dashboard](https://dash.cloudflare.com)
2. Navigate ke R2 → Create Bucket
3. Buat bucket: `karima-media-prod`
4. Generate API Token:
   - R2 → Manage R2 API Tokens
   - Create API Token
   - Set permissions: Read & Write
5. Setup Custom Domain:
   - Bucket Settings → Public Access
   - Add custom domain: media.yourdomain.com
   - Configure DNS CNAME

**Tips R2:**
- Setup lifecycle policies untuk auto-delete old files
- Enable versioning untuk backup
- Configure CORS untuk web access
- Use CDN untuk faster delivery
- Monitor storage usage dan costs

### 9. Konfigurasi Email (SMTP)

```bash
# Email Configuration
EMAIL_HOST=smtp.gmail.com
EMAIL_PORT=587
EMAIL_USER=noreply@yourdomain.com
EMAIL_PASSWORD=<app_specific_password>
```

**Cara Setup Gmail SMTP:**
1. Login ke Google Account
2. Enable 2-Factor Authentication
3. Generate App Password:
   - Account Settings → Security
   - 2-Step Verification → App passwords
   - Select "Mail" and "Other"
   - Generate password
4. Gunakan app password (bukan password akun)

**Alternatif Email Providers:**
- **SendGrid**: Reliable, good free tier
- **AWS SES**: Cost-effective untuk volume tinggi
- **Mailgun**: Developer-friendly
- **Postmark**: Excellent deliverability

**Tips Email:**
- Gunakan dedicated email untuk transactional emails
- Setup SPF, DKIM, dan DMARC records
- Monitor email deliverability
- Implement email templates
- Setup bounce handling

### 10. Konfigurasi WhatsApp (Fonnte)

```bash
# Fonnte Configuration
FONNTE_TOKEN=<your_production_token>
FONNTE_URL=https://api.fonnte.com/send
```

**Cara Mendapatkan Fonnte Token:**
1. Login ke [Fonnte](https://fonnte.com)
2. Navigate ke Dashboard
3. Copy API Token
4. Verify WhatsApp number

**Tips Fonnte:**
- Verify business WhatsApp number
- Setup message templates
- Monitor API quota
- Implement rate limiting
- Handle webhook callbacks

### 11. Konfigurasi Logging & Monitoring

```bash
# Logging Configuration
LOG_LEVEL=warn                         # production: warn atau error
LOG_FILE=logs/app.log
```

**Log Levels:**
- `debug`: Development only
- `info`: Staging
- `warn`: Production (recommended)
- `error`: Production (minimal logging)

**Tips Logging:**
- Jangan log sensitive data (passwords, tokens)
- Implement log rotation
- Setup centralized logging (ELK, Loki)
- Monitor error rates
- Setup alerting untuk critical errors

### 12. Konfigurasi CORS

```bash
# CORS Configuration
CORS_ORIGIN=https://yourdomain.com,https://www.yourdomain.com
```

**Tips CORS:**
- List semua production domains
- Jangan gunakan wildcard (*) di production
- Include www dan non-www variants
- Test CORS dari browser
- Monitor CORS errors

### 13. Konfigurasi Rate Limiting

```bash
# Rate Limiting Configuration
RATE_LIMIT_WINDOW=1m
RATE_LIMIT_LIMIT=100                   # Sesuaikan dengan traffic
```

**Tips Rate Limiting:**
- Adjust berdasarkan expected traffic
- Implement different limits untuk different endpoints
- Whitelist internal IPs
- Monitor rate limit hits
- Implement graceful degradation

## 🚀 Deployment Checklist

### Pre-Deployment
- [ ] Generate semua passwords dan secrets
- [ ] Update `.env.production` dengan values yang benar
- [ ] Verify semua API keys valid
- [ ] Test koneksi ke semua external services
- [ ] Backup existing production data
- [ ] Review security configurations

### Deployment
- [ ] Copy `.env.production` ke server (gunakan SCP atau secret manager)
- [ ] Set proper file permissions: `chmod 600 .env.production`
- [ ] Verify environment variables loaded correctly
- [ ] Run database migrations
- [ ] Test application startup
- [ ] Verify all integrations working

### Post-Deployment
- [ ] Monitor application logs
- [ ] Test critical user flows
- [ ] Verify payment processing
- [ ] Test email delivery
- [ ] Check monitoring dashboards
- [ ] Setup alerting
- [ ] Document deployment

## 🔒 Secret Management Best Practices

### Option 1: Environment Variables (Basic)
```bash
# Di server, set environment variables
export DB_PASSWORD="your_secure_password"
export JWT_SECRET="your_jwt_secret"
```

### Option 2: Docker Secrets (Recommended untuk Docker)
```bash
# Create secrets
echo "your_db_password" | docker secret create db_password -
echo "your_jwt_secret" | docker secret create jwt_secret -

# Reference di docker-compose.yml
secrets:
  - db_password
  - jwt_secret
```

### Option 3: HashiCorp Vault (Enterprise)
```bash
# Store secrets
vault kv put secret/karima/db password="your_password"
vault kv put secret/karima/jwt secret="your_secret"

# Retrieve in application
vault kv get -field=password secret/karima/db
```

### Option 4: Cloud Provider Secret Managers
- **AWS Secrets Manager**
- **Google Cloud Secret Manager**
- **Azure Key Vault**

## 🔄 Secret Rotation Strategy

### Quarterly Rotation (Every 90 Days)
1. Generate new secrets
2. Update secret manager
3. Deploy new configuration
4. Verify application working
5. Revoke old secrets
6. Document rotation

### Emergency Rotation (Breach Detected)
1. Immediately generate new secrets
2. Deploy emergency update
3. Revoke compromised secrets
4. Audit access logs
5. Investigate breach
6. Update security procedures

## 📊 Monitoring & Alerting

### Key Metrics to Monitor
- Application uptime
- Response times
- Error rates
- Database connections
- Redis memory usage
- API rate limits
- Payment success rate
- Email delivery rate

### Recommended Tools
- **Prometheus + Grafana**: Metrics & dashboards
- **ELK Stack**: Centralized logging
- **Sentry**: Error tracking
- **UptimeRobot**: Uptime monitoring
- **PagerDuty**: Incident management

## 🆘 Troubleshooting

### Database Connection Issues
```bash
# Test database connection
psql -h $DB_HOST -U $DB_USER -d $DB_NAME

# Check SSL mode
psql "postgresql://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=require"
```

### Redis Connection Issues
```bash
# Test Redis connection
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD ping

# Check Redis info
redis-cli -h $REDIS_HOST -p $REDIS_PORT -a $REDIS_PASSWORD info
```

### API Integration Issues
```bash
# Test Midtrans API
curl -X GET https://api.midtrans.com/v2/status \
  -H "Authorization: Basic $(echo -n $MIDTRANS_SERVER_KEY: | base64)"

# Test RajaOngkir API
curl -X GET "https://pro.rajaongkir.com/api/province" \
  -H "key: $RAJAONGKIR_API_KEY"
```

## 📚 Additional Resources

- [Ory Kratos Production Checklist](https://www.ory.sh/docs/kratos/production)
- [PostgreSQL Security Best Practices](https://www.postgresql.org/docs/current/security.html)
- [Redis Security](https://redis.io/docs/management/security/)
- [OWASP Security Guidelines](https://owasp.org/www-project-top-ten/)
- [12 Factor App Methodology](https://12factor.net/)

## 🤝 Support

Jika ada pertanyaan atau issues:
1. Check troubleshooting section
2. Review application logs
3. Consult team documentation
4. Contact DevOps team

---

**Last Updated**: 2026-01-04
**Version**: 1.0.0
