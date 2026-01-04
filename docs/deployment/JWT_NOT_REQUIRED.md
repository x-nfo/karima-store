# JWT Configuration - Not Required

## ❓ Apakah JWT Diperlukan?

**Jawaban: TIDAK** ❌

Karima Store menggunakan **Ory Kratos** sepenuhnya untuk authentication dan session management, sehingga JWT **tidak diperlukan** dan **tidak digunakan** di aplikasi.

## 🔍 Analisis Kode

### Yang Ditemukan:
1. ✅ `JWT_SECRET` didefinisikan di `internal/config/config.go`
2. ✅ Validasi JWT ada di production config (sekarang sudah di-comment)
3. ❌ **TIDAK ADA** penggunaan `config.JWTSecret` di seluruh codebase
4. ❌ **TIDAK ADA** JWT middleware atau handler
5. ❌ **TIDAK ADA** JWT token generation atau validation

### Kesimpulan:
JWT hanya ada di config untuk **backward compatibility** atau **future use**, tapi **tidak pernah digunakan** di aplikasi.

## 🔐 Authentication Flow Saat Ini

Karima Store menggunakan **Ory Kratos** untuk:
- ✅ User registration
- ✅ User login
- ✅ Session management
- ✅ Password recovery
- ✅ Email verification
- ✅ Multi-factor authentication (jika enabled)

## ⚙️ Perubahan Yang Dilakukan

### 1. Config Validation (`internal/config/config.go`)
```go
// Sebelum:
if c.JWTSecret == "" {
    errors = append(errors, "JWT_SECRET is required in production")
}

// Sesudah (di-comment):
// Note: JWT_SECRET is optional since we use Ory Kratos for authentication
// If you plan to use JWT in the future, uncomment the validation below:
// if c.JWTSecret == "" {
//     errors = append(errors, "JWT_SECRET is required in production")
// }
```

### 2. Environment Template (`.env.production.template`)
```bash
# Sebelum:
# JWT CONFIGURATION
# CRITICAL: Use strong random secret (minimum 64 characters)
JWT_SECRET=REPLACE_WITH_GENERATED_HEX_64_CHARS

# Sesudah:
# JWT CONFIGURATION (OPTIONAL)
# NOTE: This application uses Ory Kratos for authentication, so JWT is NOT required.
# You can safely leave this empty or set a value if you plan to use JWT in the future.
JWT_SECRET=
```

### 3. Verification Script (`scripts/verify-env.sh`)
```bash
# Sebelum:
check_required "JWT_SECRET" "REPLACE_|super_secret|your_|placeholder"
check_password_strength "JWT_SECRET" 64

# Sesudah:
# JWT is optional since we use Ory Kratos for authentication
if [ -n "$JWT_SECRET" ] && [ "$JWT_SECRET" != "REPLACE_WITH_GENERATED_HEX_64_CHARS" ]; then
    check_password_strength "JWT_SECRET" 64
else
    echo -e "${YELLOW}⚠ JWT_SECRET not set (optional - using Ory Kratos)${NC}"
fi
```

## 📝 Cara Mengisi .env.production

### Option 1: Kosongkan (Recommended)
```bash
# JWT CONFIGURATION (OPTIONAL)
JWT_SECRET=
JWT_EXPIRATION=24h
```

### Option 2: Isi untuk Future Use
Jika Anda berencana menggunakan JWT di masa depan:
```bash
# JWT CONFIGURATION (OPTIONAL)
JWT_SECRET=5c4ae3579c017a321538faddcaedce572aacbf47357b2b0d3c886258136dcc76
JWT_EXPIRATION=24h
```

## 🚀 Deployment Checklist Update

### Yang TIDAK Perlu Dilakukan:
- ❌ Generate JWT_SECRET
- ❌ Set JWT_SECRET di production
- ❌ Validate JWT_SECRET strength
- ❌ Rotate JWT_SECRET

### Yang TETAP Perlu Dilakukan:
- ✅ Setup Ory Kratos dengan benar
- ✅ Configure Kratos URLs (PUBLIC, ADMIN, UI)
- ✅ Setup Kratos session management
- ✅ Configure Kratos SMTP untuk email
- ✅ Test Kratos authentication flow

## 🔮 Jika Ingin Menggunakan JWT di Masa Depan

Jika suatu saat Anda ingin menambahkan JWT (misalnya untuk API keys atau service-to-service auth):

### 1. Uncomment Validation di `config.go`:
```go
if c.JWTSecret == "" {
    errors = append(errors, "JWT_SECRET is required in production")
}
```

### 2. Buat JWT Middleware:
```go
// internal/middleware/jwt.go
package middleware

import (
    "github.com/golang-jwt/jwt/v5"
    // ... implement JWT middleware
)
```

### 3. Gunakan di Routes:
```go
// internal/routes/routes.go
router.Use(middleware.JWTAuth())
```

### 4. Generate & Validate Tokens:
```go
// internal/auth/jwt.go
package auth

func GenerateToken(userID string) (string, error) {
    // ... implement token generation
}

func ValidateToken(tokenString string) (*Claims, error) {
    // ... implement token validation
}
```

## 📚 Referensi

- **Ory Kratos Documentation**: https://www.ory.sh/docs/kratos
- **Kratos Session Management**: https://www.ory.sh/docs/kratos/session-management
- **Kratos Production Checklist**: https://www.ory.sh/docs/kratos/production

## ✅ Summary

| Item | Status | Keterangan |
|------|--------|------------|
| JWT Required? | ❌ NO | Using Ory Kratos |
| JWT Used in Code? | ❌ NO | Not implemented |
| JWT_SECRET Required? | ❌ NO | Optional |
| Kratos Required? | ✅ YES | Primary auth system |
| Kratos URLs Required? | ✅ YES | Must be configured |

---

**Last Updated**: 2026-01-04  
**Author**: DevOps Team  
**Status**: ✅ Confirmed - JWT Not Required
