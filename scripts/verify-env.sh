#!/bin/bash

# ============================================
# KARIMA STORE - Environment Verification Script
# ============================================
# This script verifies that all required environment
# variables are set and valid for production deployment
#
# Usage: ./scripts/verify-env.sh [env-file]
# Example: ./scripts/verify-env.sh .env.production

# set -e  <-- Removed to allow all checks to run

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Counters
ERRORS=0
WARNINGS=0
PASSED=0

# Environment file
ENV_FILE="${1:-.env.production}"

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}KARIMA STORE - Environment Verification${NC}"
echo -e "${BLUE}============================================${NC}"
echo -e "Checking: ${YELLOW}$ENV_FILE${NC}"
echo ""

# Check if file exists
if [ ! -f "$ENV_FILE" ]; then
    echo -e "${RED}✗ Error: File $ENV_FILE not found!${NC}"
    exit 1
fi

# Load environment file
set -a
source "$ENV_FILE"
set +a

# Function to check if variable is set and not placeholder
check_required() {
    local var_name=$1
    local var_value="${!var_name}"
    local placeholder_pattern=$2
    
    if [ -z "$var_value" ]; then
        echo -e "${RED}✗ $var_name is not set${NC}"
        ((ERRORS++))
        return 1
    elif [[ "$var_value" =~ $placeholder_pattern ]]; then
        echo -e "${RED}✗ $var_name contains placeholder value${NC}"
        ((ERRORS++))
        return 1
    else
        echo -e "${GREEN}✓ $var_name is set${NC}"
        ((PASSED++))
        return 0
    fi
}

# Function to check optional variable
check_optional() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo -e "${YELLOW}⚠ $var_name is not set (optional)${NC}"
        ((WARNINGS++))
    else
        echo -e "${GREEN}✓ $var_name is set${NC}"
        ((PASSED++))
    fi
}

# Function to validate password strength
check_password_strength() {
    local var_name=$1
    local var_value="${!var_name}"
    local min_length=${2:-32}
    
    if [ ${#var_value} -lt $min_length ]; then
        echo -e "${RED}✗ $var_name is too short (minimum $min_length characters)${NC}"
        ((ERRORS++))
        return 1
    else
        echo -e "${GREEN}✓ $var_name has sufficient length${NC}"
        ((PASSED++))
        return 0
    fi
}

# Function to test connection
test_connection() {
    local service=$1
    local host=$2
    local port=$3
    
    if timeout 5 bash -c "cat < /dev/null > /dev/tcp/$host/$port" 2>/dev/null; then
        echo -e "${GREEN}✓ $service connection successful ($host:$port)${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${YELLOW}⚠ $service connection failed ($host:$port)${NC}"
        echo -e "  ${YELLOW}Note: This might be expected if running locally${NC}"
        ((WARNINGS++))
        return 1
    fi
}

echo -e "${BLUE}--- Server Configuration ---${NC}"
check_required "APP_PORT" "^$"
check_required "APP_ENV" "^$"
check_required "API_VERSION" "^$"
echo ""

echo -e "${BLUE}--- Database Configuration ---${NC}"
check_required "DB_HOST" "your_|placeholder"
check_required "DB_PORT" "^$"
check_required "DB_USER" "your_|placeholder"
check_required "DB_PASSWORD" "REPLACE_|your_|placeholder|lokal"
check_required "DB_NAME" "^$"
check_required "DB_SSL_MODE" "^$"

if [ "$APP_ENV" = "production" ] && [ "$DB_SSL_MODE" != "require" ] && [ "$DB_SSL_MODE" != "verify-full" ]; then
    echo -e "${RED}✗ DB_SSL_MODE should be 'require' or 'verify-full' in production${NC}"
    ((ERRORS++))
fi

# Test database connection
if [ -n "$DB_HOST" ] && [ "$DB_HOST" != "your_production_db_host" ]; then
    test_connection "PostgreSQL" "$DB_HOST" "$DB_PORT"
fi
echo ""

echo -e "${BLUE}--- Redis Configuration ---${NC}"
check_required "REDIS_HOST" "your_|placeholder"
check_required "REDIS_PORT" "^$"
check_required "REDIS_PASSWORD" "REPLACE_|your_|placeholder|^$"

# Test Redis connection
if [ -n "$REDIS_HOST" ] && [ "$REDIS_HOST" != "your_production_redis_host" ]; then
    test_connection "Redis" "$REDIS_HOST" "$REDIS_PORT"
fi
echo ""

echo -e "${BLUE}--- Authentication (Kratos) ---${NC}"
check_required "KRATOS_PUBLIC_URL" "your|placeholder"
check_required "KRATOS_ADMIN_URL" "your|placeholder"
check_required "KRATOS_UI_URL" "your|placeholder"

if [[ "$KRATOS_PUBLIC_URL" != https://* ]] && [ "$APP_ENV" = "production" ]; then
    echo -e "${YELLOW}⚠ KRATOS_PUBLIC_URL should use HTTPS in production${NC}"
    ((WARNINGS++))
fi
echo ""

echo -e "${BLUE}--- JWT Configuration (Optional - Using Ory Kratos) ---${NC}"
# JWT is optional since we use Ory Kratos for authentication
if [ -n "$JWT_SECRET" ] && [ "$JWT_SECRET" != "REPLACE_WITH_GENERATED_HEX_64_CHARS" ]; then
    check_password_strength "JWT_SECRET" 64
else
    echo -e "${YELLOW}⚠ JWT_SECRET not set (optional - using Ory Kratos)${NC}"
    ((WARNINGS++))
fi
check_optional "JWT_EXPIRATION"
echo ""

echo -e "${BLUE}--- Payment Gateway (Midtrans) ---${NC}"
check_required "MIDTRANS_SERVER_KEY" "YOUR_|your_|placeholder"
check_required "MIDTRANS_CLIENT_KEY" "YOUR_|your_|placeholder"
check_required "MIDTRANS_IS_PRODUCTION" "^$"

if [ "$APP_ENV" = "production" ] && [ "$MIDTRANS_IS_PRODUCTION" != "true" ]; then
    echo -e "${RED}✗ MIDTRANS_IS_PRODUCTION should be 'true' in production${NC}"
    ((ERRORS++))
fi

if [[ "$MIDTRANS_SERVER_KEY" == *"sandbox"* ]] && [ "$APP_ENV" = "production" ]; then
    echo -e "${RED}✗ Using sandbox Midtrans key in production!${NC}"
    ((ERRORS++))
fi
echo ""

echo -e "${BLUE}--- Shipping (RajaOngkir) ---${NC}"
check_required "RAJAONGKIR_API_KEY" "YOUR_|your_|placeholder"
check_required "RAJAONGKIR_BASE_URL" "^$"

if [[ "$RAJAONGKIR_BASE_URL" == *"sandbox"* ]] && [ "$APP_ENV" = "production" ]; then
    echo -e "${YELLOW}⚠ Using sandbox RajaOngkir URL in production${NC}"
    ((WARNINGS++))
fi
echo ""

echo -e "${BLUE}--- Storage (Cloudflare R2) ---${NC}"
check_required "FILE_STORAGE" "^$"
check_required "R2_ACCOUNT_ID" "your_|placeholder"
check_required "R2_ACCESS_KEY_ID" "your_|placeholder"
check_required "R2_SECRET_ACCESS_KEY" "your_|placeholder"
check_required "R2_BUCKET_NAME" "^$"

if [ "$APP_ENV" = "production" ] && [ "$FILE_STORAGE" = "local" ]; then
    echo -e "${YELLOW}⚠ FILE_STORAGE is 'local' in production (should be 'r2')${NC}"
    ((WARNINGS++))
fi
echo ""

echo -e "${BLUE}--- Email Configuration ---${NC}"
check_required "EMAIL_HOST" "^$"
check_required "EMAIL_PORT" "^$"
check_required "EMAIL_USER" "your_|placeholder"
check_required "EMAIL_PASSWORD" "your_|placeholder|minyakkayuputih"
echo ""

echo -e "${BLUE}--- WhatsApp (Fonnte) ---${NC}"
check_required "FONNTE_TOKEN" "YOUR_|your_|placeholder"
check_required "FONNTE_URL" "^$"
echo ""

echo -e "${BLUE}--- Logging Configuration ---${NC}"
check_required "LOG_LEVEL" "^$"
check_required "LOG_FILE" "^$"

if [ "$APP_ENV" = "production" ] && [ "$LOG_LEVEL" = "debug" ]; then
    echo -e "${YELLOW}⚠ LOG_LEVEL is 'debug' in production (should be 'warn' or 'error')${NC}"
    ((WARNINGS++))
fi
echo ""

echo -e "${BLUE}--- CORS Configuration ---${NC}"
check_required "CORS_ORIGIN" "^$"

if [[ "$CORS_ORIGIN" == *"localhost"* ]] && [ "$APP_ENV" = "production" ]; then
    echo -e "${YELLOW}⚠ CORS_ORIGIN contains 'localhost' in production${NC}"
    ((WARNINGS++))
fi

if [ "$CORS_ORIGIN" = "*" ] && [ "$APP_ENV" = "production" ]; then
    echo -e "${RED}✗ CORS_ORIGIN is wildcard (*) in production - security risk!${NC}"
    ((ERRORS++))
fi
echo ""

echo -e "${BLUE}--- Rate Limiting ---${NC}"
check_required "RATE_LIMIT_WINDOW" "^$"
check_required "RATE_LIMIT_LIMIT" "^$"
echo ""

# Summary
echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Verification Summary${NC}"
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${YELLOW}Warnings: $WARNINGS${NC}"
echo -e "${RED}Errors: $ERRORS${NC}"
echo ""

if [ $ERRORS -gt 0 ]; then
    echo -e "${RED}✗ Verification FAILED - Please fix errors before deploying${NC}"
    exit 1
elif [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠ Verification completed with warnings${NC}"
    echo -e "${YELLOW}  Review warnings before deploying to production${NC}"
    exit 0
else
    echo -e "${GREEN}✓ Verification PASSED - Environment is ready for deployment${NC}"
    exit 0
fi
