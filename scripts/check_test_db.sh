#!/bin/bash

# Test Database Connection Diagnostic Script
# This script helps diagnose why tests cannot connect to the database

echo "=========================================="
echo "Test Database Connection Diagnostic Tool"
echo "=========================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓${NC} $2"
    else
        echo -e "${RED}✗${NC} $2"
    fi
}

# Check if Docker is installed
echo "1. Checking Docker installation..."
if command -v docker &> /dev/null; then
    print_status 0 "Docker is installed"
    DOCKER_AVAILABLE=true
else
    print_status 1 "Docker is not installed"
    DOCKER_AVAILABLE=false
fi
echo ""

# Check if Docker Compose is installed
echo "2. Checking Docker Compose installation..."
if command -v docker-compose &> /dev/null; then
    print_status 0 "Docker Compose is installed"
    DOCKER_COMPOSE_AVAILABLE=true
else
    print_status 1 "Docker Compose is not installed"
    DOCKER_COMPOSE_AVAILABLE=false
fi
echo ""

# Check Docker containers
if [ "$DOCKER_AVAILABLE" = true ]; then
    echo "3. Checking Docker containers..."
    CONTAINERS=$(docker ps -a --filter "name=karima" --format "{{.Names}}\t{{.Status}}")

    if [ -z "$CONTAINERS" ]; then
        echo -e "${YELLOW}⚠${NC} No Karima Store containers found"
    else
        echo "$CONTAINERS" | while read -r line; do
            if echo "$line" | grep -q "Up"; then
                echo -e "${GREEN}✓${NC} $line"
            else
                echo -e "${RED}✗${NC} $line"
            fi
        done
    fi
    echo ""
fi

# Check PostgreSQL service
echo "4. Checking PostgreSQL service..."
if [ "$DOCKER_AVAILABLE" = true ]; then
    PG_CONTAINER=$(docker ps -q --filter "name=db")
    if [ -n "$PG_CONTAINER" ]; then
        print_status 0 "PostgreSQL container is running"

        # Get container details
        PG_PORT=$(docker port "$PG_CONTAINER" 5432 2>/dev/null | cut -d: -f2)
        if [ -n "$PG_PORT" ]; then
            echo "   PostgreSQL is mapped to port: $PG_PORT"
        fi

        # Test connection from inside container
        echo "   Testing connection from inside container..."
        if docker exec "$PG_CONTAINER" pg_isready -U karima_store &> /dev/null; then
            print_status 0 "PostgreSQL is ready inside container"
        else
            print_status 1 "PostgreSQL is not ready inside container"
        fi
    else
        print_status 1 "PostgreSQL container is not running"
    fi
else
    echo -e "${YELLOW}⚠${NC} Docker not available, skipping container checks"
fi
echo ""

# Check if PostgreSQL is running locally
echo "5. Checking local PostgreSQL..."
if command -v psql &> /dev/null; then
    print_status 0 "psql client is installed"

    # Test connection to localhost
    echo "   Testing connection to localhost:5432..."
    if PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db -c "SELECT 1;" &> /dev/null; then
        print_status 0 "Can connect to localhost:5432"
    else
        print_status 1 "Cannot connect to localhost:5432"
        echo "   Error details:"
        PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d karima_db -c "SELECT 1;" 2>&1 | grep -E "error|could not|failed" | head -3
    fi

    # Test connection to 127.0.0.1
    echo "   Testing connection to 127.0.0.1:5432..."
    if PGPASSWORD=lokal psql -h 127.0.0.1 -p 5432 -U karima_store -d karima_db -c "SELECT 1;" &> /dev/null; then
        print_status 0 "Can connect to 127.0.0.1:5432"
    else
        print_status 1 "Cannot connect to 127.0.0.1:5432"
    fi
else
    print_status 1 "psql client is not installed"
fi
echo ""

# Check environment variables
echo "6. Checking environment configuration..."
if [ -f ".env" ]; then
    print_status 0 ".env file exists"

    echo "   Database configuration from .env:"
    grep "^DB_" .env | while read -r line; do
        echo "     $line"
    done
else
    print_status 1 ".env file not found"
fi
echo ""

# Check test configuration
echo "7. Checking test configuration..."
if [ -f "internal/config/test_config.go" ]; then
    print_status 0 "test_config.go exists"

    echo "   Test database configuration:"
    grep -A 5 "DBHost:" internal/config/test_config.go | head -6 | sed 's/^/     /'
else
    print_status 1 "test_config.go not found"
fi
echo ""

# Check if test database exists
echo "8. Checking if test database exists..."
if command -v psql &> /dev/null; then
    if PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d postgres -c "\l" | grep -q "karima_db"; then
        print_status 0 "Database 'karima_db' exists"
    else
        print_status 1 "Database 'karima_db' does not exist"
        echo "   To create it, run:"
        echo "     docker exec -it <postgres_container> psql -U karima_store -c 'CREATE DATABASE karima_db;'"
    fi
fi
echo ""

# Summary and recommendations
echo "=========================================="
echo "Summary and Recommendations"
echo "=========================================="
echo ""

if [ "$DOCKER_AVAILABLE" = true ] && [ -n "$PG_CONTAINER" ]; then
    echo -e "${GREEN}Docker PostgreSQL is running${NC}"
    echo ""
    echo "Recommended solutions:"
    echo ""
    echo "1. Update test_config.go to use Docker service name:"
    echo "   Change: DBHost: \"localhost\""
    echo "   To:     DBHost: \"db\""
    echo ""
    echo "2. Or expose PostgreSQL port and use localhost:"
    echo "   Make sure docker-compose.yml has:"
    echo "   ports:"
    echo "     - \"5432:5432\""
    echo ""
    echo "3. Or run tests inside Docker network:"
    echo "   docker-compose -f docker-compose.test.yml up"
    echo ""
elif command -v psql &> /dev/null && PGPASSWORD=lokal psql -h localhost -p 5432 -U karima_store -d postgres -c "\l" | grep -q "karima_db"; then
    echo -e "${GREEN}Local PostgreSQL is available and database exists${NC}"
    echo ""
    echo "Test configuration looks correct. Try running:"
    echo "  go test -v ./internal/repository/user_repository_test.go"
    echo ""
else
    echo -e "${RED}Database is not accessible${NC}"
    echo ""
    echo "To fix this issue:"
    echo ""
    echo "Option 1: Start Docker containers"
    echo "  docker-compose up -d db"
    echo ""
    echo "Option 2: Install and start PostgreSQL locally"
    echo "  sudo apt install postgresql postgresql-contrib"
    echo "  sudo systemctl start postgresql"
    echo "  sudo -u postgres createuser karima_store"
    echo "  sudo -u postgres createdb karima_db -O karima_store"
    echo "  sudo -u postgres psql -c \"ALTER USER karima_store PASSWORD 'lokal';\""
    echo ""
    echo "Option 3: Use Docker Compose for tests"
    echo "  See docs/TEST_DATABASE_TROUBLESHOOTING.md for details"
    echo ""
fi

echo "For detailed troubleshooting guide, see:"
echo "  docs/TEST_DATABASE_TROUBLESHOOTING.md"
echo ""
