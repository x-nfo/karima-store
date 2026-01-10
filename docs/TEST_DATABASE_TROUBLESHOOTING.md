# Test Database Connection Troubleshooting Guide

## Problem
Tests cannot connect to the database when running `go test`.

## Root Cause
The test configuration in [`internal/config/test_config.go`](../internal/config/test_config.go) uses hardcoded values that may not match your actual database setup:

```go
DBHost:     "localhost",  // ← May be wrong if using Docker
DBPort:     "5432",
DBUser:     "karima_store",
DBPassword: "lokal",
DBName:     "karima_db",
```

## Diagnostic Steps

### Step 1: Check if Database is Running

**Option A: Using Docker Compose**
```bash
docker-compose ps
```
Look for the `db` service. It should show "Up" status.

**Option B: Using Docker directly**
```bash
docker ps | grep postgres
```

**Option C: Using psql (if installed locally)**
```bash
psql -h localhost -p 5432 -U karima_store -d karima_db
```

### Step 2: Test Database Connection

Create a test connection script:

```bash
# Create a test file
cat > test_db_connection.sh << 'EOF'
#!/bin/bash

echo "Testing database connection..."

# Try localhost
echo "Testing localhost:5432..."
PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db -c "SELECT version();" 2>&1

# Try Docker service name
echo ""
echo "Testing db:5432 (Docker)..."
PGPASSWORD=lokal psql -h db -p 5432 -U karima_store -d karima_db -c "SELECT version();" 2>&1

# Try 127.0.0.1
echo ""
echo "Testing 127.0.0.1:5432..."
PGPASSWORD=lokal psql -h 127.0.0.1 -p 5432 -U karima_store -d karima_db -c "SELECT version();" 2>&1
EOF

chmod +x test_db_connection.sh
./test_db_connection.sh
```

### Step 3: Check Docker Network

```bash
# List Docker networks
docker network ls

# Check which network your containers use
docker inspect <container_id> | grep NetworkMode
```

## Solutions

### Solution 1: Update Test Config to Use Environment Variables (Recommended)

Update [`internal/config/test_config.go`](../internal/config/test_config.go) to read from environment variables:

```go
package config

import (
	"os"
)

// TestConfig returns a test configuration for unit tests
func TestConfig() *Config {
	return &Config{
		AppEnv:            "test",
		AppPort:           "8080",
		DBHost:            getEnv("TEST_DB_HOST", "localhost"),
		DBPort:            getEnv("TEST_DB_PORT", "5432"),
		DBUser:            getEnv("TEST_DB_USER", "karima_store"),
		DBPassword:        getEnv("TEST_DB_PASSWORD", "lokal"),
		DBName:            getEnv("TEST_DB_NAME", "karima_db_test"),
		DBSSLMode:         "disable",
		RedisHost:         getEnv("TEST_REDIS_HOST", "localhost"),
		RedisPort:         getEnv("TEST_REDIS_PORT", "6380"),
		RedisPassword:     "",
		JWTSecret:         "test-secret-key-for-testing-only",
		GoogleKey:         "test-google-key",
		GoogleSecret:      "test-google-secret",
		CallbackURL:       "http://localhost:8080/api/v1/auth/google/callback",
		FileStorage:       "local",
		RateLimitLimit:    "100",
		RateLimitWindow:   "1m",
		R2AccountID:       "",
		R2AccessKeyID:     "",
		R2SecretAccessKey: "",
		R2BucketName:      "",
		R2PublicURL:       "",
		R2Region:          "",
	}
}
```

Then create a `.env.test` file:

```bash
# Test Database Configuration
TEST_DB_HOST=db
TEST_DB_PORT=5432
TEST_DB_USER=karima_store
TEST_DB_PASSWORD=lokal
TEST_DB_NAME=karima_db_test

# Test Redis Configuration
TEST_REDIS_HOST=redis
TEST_REDIS_PORT=6379
```

### Solution 2: Use Docker Compose for Tests

Update [`docker-compose.yml`](../docker-compose.yml) to include test database service:

```yaml
version: '3.8'

services:
  # Main database
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: karima_store
      POSTGRES_PASSWORD: lokal
      POSTGRES_DB: karima_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U karima_store"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Test database (separate for testing)
  db_test:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: karima_store
      POSTGRES_PASSWORD: lokal
      POSTGRES_DB: karima_db_test
    ports:
      - "5433:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U karima_store"]
      interval: 10s
      timeout: 5s
      retries: 5

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"

volumes:
  postgres_data:
```

Then update test config to use test database:

```go
DBHost:     "localhost",
DBPort:     "5433",  // Different port for test database
DBName:     "karima_db_test",
```

### Solution 3: Run Tests Inside Docker

Create a `docker-compose.test.yml`:

```yaml
version: '3.8'

services:
  db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: karima_store
      POSTGRES_PASSWORD: lokal
      POSTGRES_DB: karima_db_test
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
      - TEST_DB_NAME=karima_db_test
    depends_on:
      db:
        condition: service_healthy
```

Run tests with:
```bash
docker-compose -f docker-compose.test.yml up --abort-on-container-exit
```

### Solution 4: Create Test Database

If running locally, create the test database:

```bash
# Connect to PostgreSQL
psql -U karima_store -d postgres

# Create test database
CREATE DATABASE karima_db_test;

# Exit
\q
```

### Solution 5: Quick Fix for Local Development

If you're running tests locally and have PostgreSQL on port 5432:

1. Make sure PostgreSQL is running
2. Create the database:
   ```bash
   createdb -U karima_store karima_db
   ```
3. Run tests:
   ```bash
   go test ./...
   ```

## Running Tests with Database

### Option 1: Run All Tests
```bash
go test -v ./...
```

### Option 2: Run Specific Test File
```bash
go test -v ./internal/repository/user_repository_test.go
```

### Option 3: Run Tests with Coverage
```bash
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Option 4: Run Tests with Race Detection
```bash
go test -race ./...
```

## Common Error Messages and Solutions

### Error: `connection refused`
**Cause**: Database not running or wrong host/port
**Solution**:
- Check if PostgreSQL is running: `docker ps` or `systemctl status postgresql`
- Verify host and port in test config

### Error: `authentication failed`
**Cause**: Wrong username or password
**Solution**:
- Check credentials in test config match database
- Verify user exists: `psql -U postgres -c "\du"`

### Error: `database "karima_db" does not exist`
**Cause**: Database not created
**Solution**:
- Create database: `createdb -U karima_store karima_db`
- Or update test config to use existing database

### Error: `timeout`
**Cause**: Connection timeout
**Solution**:
- Check network connectivity
- Verify firewall rules
- Increase timeout in test setup

## Best Practices

1. **Use Separate Test Database**: Don't use production database for tests
2. **Environment Variables**: Use environment variables for test configuration
3. **Docker Compose**: Use Docker Compose for consistent test environment
4. **Cleanup**: Always clean up test data after tests
5. **Isolation**: Each test should be independent

## Makefile Commands

Add these commands to your [`Makefile`](../Makefile):

```makefile
# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Run tests in Docker
test-docker:
	docker-compose -f docker-compose.test.yml up --abort-on-container-exit

# Create test database
create-test-db:
	docker exec -it karima_store-db psql -U karima_store -c "CREATE DATABASE IF NOT EXISTS karima_db_test;"

# Check database connection
check-db:
	@echo "Checking database connection..."
	@docker exec -it karima_store-db pg_isready -U karima_store
```

## Verification

After implementing a solution, verify tests work:

```bash
# Run a single test to verify
go test -v -run TestUserRepository_Create ./internal/repository

# If successful, run all tests
go test -v ./...
```

## Additional Resources

- [GORM Database Connection](https://gorm.io/docs/connecting_to_the_database.html)
- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Docker Compose Healthchecks](https://docs.docker.com/compose/compose-file/compose-file-v3/#healthcheck)
