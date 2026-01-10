# Panduan Debugging Koneksi Database untuk Testing

## Masalah
Test tidak bisa terhubung ke database saat menjalankan `go test`.

## Penyebab Utama
Konfigurasi test di [`internal/config/test_config.go`](../internal/config/test_config.go) menggunakan nilai hardcoded yang mungkin tidak sesuai dengan setup database Anda:

```go
DBHost:     "localhost",  // ← Mungkin salah jika menggunakan Docker
DBPort:     "5432",
DBUser:     "karima_store",
DBPassword: "lokal",
DBName:     "karima_db",
```

## Cara Cek Masalah

### Langkah 1: Jalankan Script Diagnostik

Script ini akan mengecek semua kemungkinan masalah:

```bash
./scripts/check_test_db.sh
```

Script ini akan mengecek:
- Apakah Docker terinstall dan berjalan
- Apakah container PostgreSQL berjalan
- Apakah bisa connect ke database dari localhost
- Konfigurasi environment
- Apakah database test sudah dibuat

### Langkah 2: Cek Container Docker

```bash
# Lihat semua container yang berjalan
docker ps

# Lihat container karima_store
docker ps | grep karima

# Cek status container database
docker logs karima_store-db
```

### Langkah 3: Test Koneksi Manual

```bash
# Coba connect ke database
docker exec -it karima_store-db psql -U karima_store -d karima_db

# Atau dari host
PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db
```

## Solusi

### Solusi 1: Update Test Config untuk Menggunakan Docker (Rekomendasi)

Edit [`internal/config/test_config.go`](../internal/config/test_config.go):

```go
func TestConfig() *Config {
	return &Config{
		AppEnv:            "test",
		AppPort:           "8080",
		DBHost:            "db",           // ← Ubah dari "localhost" ke "db"
		DBPort:            "5432",
		DBUser:            "karima_store",
		DBPassword:        "lokal",
		DBName:            "karima_db",
		// ... config lainnya
	}
}
```

**Catatan**: Ini hanya bekerja jika test dijalankan di dalam Docker network.

### Solusi 2: Expose Port PostgreSQL di Docker Compose

Pastikan [`docker-compose.yml`](../docker-compose.yml) memiliki konfigurasi ini:

```yaml
services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: karima_store
      POSTGRES_PASSWORD: lokal
      POSTGRES_DB: karima_db
    ports:
      - "5432:5432"  # ← Pastikan port ini ada
    volumes:
      - postgres_data:/var/lib/postgresql/data
```

Kemudian restart container:

```bash
docker-compose down
docker-compose up -d db
```

### Solusi 3: Jalankan Test di Dalam Docker

Buat file `docker-compose.test.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: karima_store
      POSTGRES_PASSWORD: lokal
      POSTGRES_DB: karima_db
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U karima_store"]
      interval: 10s
      timeout: 5s
      retries: 5

  test:
    build: .
    command: go test -v ./...
    environment:
      - TEST_DB_HOST=db
      - TEST_DB_PORT=5432
      - TEST_DB_USER=karima_store
      - TEST_DB_PASSWORD=lokal
      - TEST_DB_NAME=karima_db
    depends_on:
      db:
        condition: service_healthy
```

Jalankan test:

```bash
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Solusi 4: Install PostgreSQL Lokal

Jika tidak ingin menggunakan Docker:

```bash
# Install PostgreSQL (Ubuntu/Debian)
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start service
sudo systemctl start postgresql

# Buat user dan database
sudo -u postgres createuser karima_store
sudo -u postgres createdb karima_db -O karima_store
sudo -u postgres psql -c "ALTER USER karima_store PASSWORD 'lokal';"

# Test koneksi
PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db

# Jalankan test
go test -v ./...
```

### Solusi 5: Buat Database Test di Docker

```bash
# Masuk ke container PostgreSQL
docker exec -it karima_store-db psql -U karima_store -d postgres

# Buat database
CREATE DATABASE karima_db;

# Exit
\q
```

## Menjalankan Test

### Jalankan Semua Test
```bash
go test -v ./...
```

### Jalankan Test Spesifik
```bash
go test -v ./internal/repository/user_repository_test.go
```

### Jalankan Test dengan Coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Jalankan Test dengan Race Detection
```bash
go test -race ./...
```

## Error Umum dan Solusi

### Error: `connection refused`
**Penyebab**: Database tidak berjalan atau host/port salah
**Solusi**:
- Cek apakah PostgreSQL berjalan: `docker ps` atau `systemctl status postgresql`
- Verifikasi host dan port di test config

### Error: `authentication failed`
**Penyebab**: Username atau password salah
**Solusi**:
- Cek credentials di test config
- Verifikasi user exists: `psql -U postgres -c "\du"`

### Error: `database "karima_db" does not exist`
**Penyebab**: Database belum dibuat
**Solusi**:
- Buat database: `docker exec -it karima_store-db psql -U karima_store -c "CREATE DATABASE karima_db;"`
- Atau update test config untuk menggunakan database yang sudah ada

### Error: `timeout`
**Penyebab**: Connection timeout
**Solusi**:
- Cek network connectivity
- Verifikasi firewall rules
- Increase timeout di test setup

## Tips Tambahan

### Cek Log Container
```bash
# Lihat log database
docker logs karima_store-db

# Lihat log real-time
docker logs -f karima_store-db
```

### Cek Network Docker
```bash
# Lihat network yang tersedia
docker network ls

# Cek network container
docker inspect karima_store-db | grep NetworkMode
```

### Restart Container
```bash
# Restart database
docker-compose restart db

# Atau restart semua
docker-compose restart
```

## Best Practices

1. **Gunakan Database Test Terpisah**: Jangan gunakan database production untuk testing
2. **Environment Variables**: Gunakan environment variables untuk konfigurasi test
3. **Docker Compose**: Gunakan Docker Compose untuk environment test yang konsisten
4. **Cleanup**: Selalu cleanup data test setelah test selesai
5. **Isolasi**: Setiap test harus independen

## Resources Tambahan

- [GORM Database Connection](https://gorm.io/docs/connecting_to_the_database.html)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Docker Compose Healthchecks](https://docs.docker.com/compose/compose-file/compose-file-v3/#healthcheck)

## Quick Reference

### Cek Status Database
```bash
./scripts/check_test_db.sh
```

### Start Database (Docker)
```bash
docker-compose up -d db
```

### Start Database (Local)
```bash
sudo systemctl start postgresql
```

### Jalankan Test
```bash
go test -v ./...
```

### Cek Log Test
```bash
go test -v ./internal/repository/user_repository_test.go 2>&1 | tee test.log
```

## Troubleshooting Flowchart

```
Start
  │
  ├─→ Jalankan check_test_db.sh
  │   │
  │   ├─→ Docker berjalan?
  │   │   ├─ Ya → Cek container PostgreSQL
  │   │   └─ Tidak → Install Docker atau gunakan PostgreSQL lokal
  │   │
  │   ├─→ PostgreSQL berjalan?
  │   │   ├─ Ya → Cek koneksi
  │   │   └─ Tidak → Start PostgreSQL
  │   │
  │   ├─→ Bisa connect?
  │   │   ├─ Ya → Cek database exists
  │   │   └─ Tidak → Cek host/port/credentials
  │   │
  │   └─→ Database exists?
  │       ├─ Ya → Update test config
  │       └─ Tidak → Buat database
  │
  └─→ Jalankan test
      │
      └─→ Berhasil?
          ├─ Ya → Selesai!
          └─ Tidak → Cek error message dan cari solusi di panduan ini
```

