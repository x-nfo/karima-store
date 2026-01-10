# Panduan Testing dengan Docker Container (karima_postgres & karima_redis)

## Konfigurasi Docker Compose

Berdasarkan [`docker-compose.yml`](../docker-compose.yml), Anda menggunakan:

```yaml
services:
  db:
    container_name: karima_postgres
    ports:
      - "5432:5432"  # Expose ke host

  redis:
    container_name: karima_redis
    ports:
      - "6380:6379"  # Expose ke host (port 6380 di host)
```

## Apakah Perlu Ubah di .env?

**TIDAK PERLU** untuk testing jika Anda menjalankan test dari host (lokal), karena port sudah di-expose.

## Cara Menjalankan Test

### Langkah 1: Pastikan Container Berjalan

```bash
# Start database dan redis
docker-compose up -d db redis

# Cek status
docker ps
```

Output yang diharapkan:
```
CONTAINER ID   IMAGE                  STATUS         PORTS
abc123         postgres:15-alpine     Up 2 minutes   0.0.0.0:5432->5432/tcp
def456         redis:7-alpine         Up 2 minutes   0.0.0.0:6380->6379/tcp
```

### Langkah 2: Verifikasi Database Ada

```bash
# Masuk ke container PostgreSQL
docker exec -it karima_postgres psql -U karima_store -d postgres

# List semua database
\l

# Jika database kariman_db belum ada, buat:
CREATE DATABASE kariman_db;

# Exit
\q
```

### Langkah 3: Jalankan Test

**Opsi A: Jalankan semua test**
```bash
go test -v ./...
```

**Opsi B: Jalankan test spesifik**
```bash
go test -v ./internal/repository/user_repository_test.go
```

**Opsi C: Jalankan test dengan coverage**
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Konfigurasi Test

### File yang Dibuat

1. **[`.env.test`](../.env.test)** - Konfigurasi environment untuk testing
2. **[`internal/config/test_config.go`](../internal/config/test_config.go)** - Sudah diupdate untuk membaca dari environment variables

### Cara Kerja

Test config sekarang membaca dari environment variables dengan default values:

```go
DBHost:     getEnv("TEST_DB_HOST", "localhost"),
DBPort:     getEnv("TEST_DB_PORT", "5432"),
DBUser:     getEnv("TEST_DB_USER", "karima_store"),
DBPassword: getEnv("TEST_DB_PASSWORD", "lokal"),
DBName:     getEnv("TEST_DB_NAME", "karima_db"),
```

Jika Anda tidak set environment variables, test akan menggunakan default values:
- Database: `localhost:5432`
- Redis: `localhost:6380`

## Troubleshooting

### Error: `connection refused`

**Penyebab**: Container tidak berjalan atau port tidak expose

**Solusi**:
```bash
# Cek container berjalan
docker ps | grep karima

# Jika tidak berjalan, start container
docker-compose up -d db redis

# Cek port
docker port karima_postgres
docker port karima_redis
```

### Error: `database "karima_db" does not exist`

**Penyebab**: Database belum dibuat

**Solusi**:
```bash
# Buat database
docker exec -it karima_postgres psql -U karima_store -d postgres -c "CREATE DATABASE kariman_db;"

# Atau update .env.test untuk menggunakan database yang sudah ada
```

### Error: `authentication failed`

**Penyebab**: Username atau password salah

**Solusi**:
```bash
# Cek user di database
docker exec -it karima_postgres psql -U postgres -d postgres -c "\du"

# Jika user karima_store belum ada, buat:
docker exec -it karima_postgres psql -U postgres -d postgres -c "CREATE USER karima_store WITH PASSWORD 'lokal';"
docker exec -it karima_postgres psql -U postgres -d postgres -c "ALTER USER karima_store CREATEDB;"
```

### Error: `timeout`

**Penyebab**: Database belum siap

**Solusi**:
```bash
# Tunggu database siap
docker exec karima_postgres pg_isready -U karima_store

# Atau restart container
docker-compose restart db
```

## Menjalankan Test di Dalam Docker Network

Jika ingin test berjalan di dalam Docker network:

### Update .env.test

```bash
# Gunakan service name Docker, bukan localhost
TEST_DB_HOST=db
TEST_DB_PORT=5432
TEST_REDIS_HOST=redis
TEST_REDIS_PORT=6379
```

### Buat docker-compose.test.yml

```yaml
version: '3.8'

services:
  test:
    build: .
    command: go test -v ./...
    environment:
      - TEST_DB_HOST=db
      - TEST_DB_PORT=5432
      - TEST_DB_USER=karima_store
      - TEST_DB_PASSWORD=lokal
      - TEST_DB_NAME=kariman_db
      - TEST_REDIS_HOST=redis
      - TEST_REDIS_PORT=6379
    depends_on:
      db:
        condition: service_healthy
      redis:
        condition: service_healthy
```

### Jalankan Test di Docker

```bash
docker-compose -f docker-compose.yml -f docker-compose.test.yml up test --abort-on-container-exit
```

## Cek Koneksi Database

### Dari Host

```bash
# Test koneksi PostgreSQL
PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d kariman_db -c "SELECT version();"

# Test koneksi Redis
redis-cli -h localhost -p 6380 ping
```

### Dari Dalam Container

```bash
# Test koneksi PostgreSQL dari dalam container
docker exec -it karima_postgres psql -U karima_store -d karima_db -c "SELECT version();"

# Test koneksi Redis dari dalam container
docker exec -it karima_redis redis-cli ping
```

## Script Bantuan

Gunakan script diagnostik untuk mengecek status:

```bash
./scripts/check_test_db.sh
```

## Tips

1. **Selalu cek container status sebelum menjalankan test**
   ```bash
   docker ps | grep karima
   ```

2. **Gunakan database terpisah untuk testing**
   - Production: `karima_db`
   - Test: `karima_db_test` (opsional)

3. **Cleanup test data setelah testing**
   - Test setup sudah otomatis cleanup data
   - Tapi Anda bisa manual: `docker exec -it karima_postgres psql -U karima_store -d karima_db -c "TRUNCATE TABLE users CASCADE;"`

4. **Monitor log container jika ada error**
   ```bash
   docker logs -f karima_postgres
   docker logs -f karima_redis
   ```

## Quick Reference

### Start Containers
```bash
docker-compose up -d db redis
```

### Stop Containers
```bash
docker-compose down
```

### Cek Status
```bash
docker ps
docker logs karima_postgres
docker logs karima_redis
```

### Jalankan Test
```bash
go test -v ./...
```

### Test Koneksi
```bash
PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db -c "SELECT 1;"
redis-cli -h localhost -p 6380 ping
```

## Ringkasan

✅ **Tidak perlu ubah .env** untuk testing dari host
✅ **Port sudah di-expose**: PostgreSQL 5432, Redis 6380
✅ **Test config sudah fleksibel**: Bisa gunakan environment variables atau default values
✅ **Container names**: `karima_postgres` dan `karima_redis`

Untuk detail lebih lanjut, lihat:
- [`docs/PANDUAN_DEBUGGING_TEST_DB.md`](PANDUAN_DEBUGGING_TEST_DB.md)
- [`docs/TEST_DATABASE_TROUBLESHOOTING.md`](TEST_DATABASE_TROUBLESHOOTING.md)
