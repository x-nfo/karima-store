#!/bin/bash

# ============================================
# KARIMA STORE - Environment Secrets Generator
# ============================================
# This script generates secure random passwords and secrets
# for production environment configuration
# 
# Usage: ./scripts/generate-env-secrets.sh

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}KARIMA STORE - Secrets Generator${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# Function to generate random password
generate_password() {
    local length=$1
    openssl rand -base64 $length | tr -d "=+/" | cut -c1-$length
}

# Function to generate hex string
generate_hex() {
    local length=$1
    openssl rand -hex $length
}

# Function to generate alphanumeric
generate_alphanum() {
    local length=$1
    LC_ALL=C tr -dc 'A-Za-z0-9' < /dev/urandom | head -c $length
}

echo -e "${GREEN}Generating secure secrets...${NC}"
echo ""

# Database Password
DB_PASSWORD=$(generate_password 32)
echo -e "${YELLOW}Database Configuration:${NC}"
echo "DB_PASSWORD=$DB_PASSWORD"
echo ""

# Redis Password
REDIS_PASSWORD=$(generate_password 32)
echo -e "${YELLOW}Redis Configuration:${NC}"
echo "REDIS_PASSWORD=$REDIS_PASSWORD"
echo ""

# JWT Secret (64 characters for extra security)
JWT_SECRET=$(generate_hex 32)
echo -e "${YELLOW}JWT Configuration:${NC}"
echo "JWT_SECRET=$JWT_SECRET"
echo ""

# API Key for internal services
API_KEY=$(generate_hex 24)
echo -e "${YELLOW}API Key:${NC}"
echo "API_KEY=$API_KEY"
echo ""

# Session Secret
SESSION_SECRET=$(generate_hex 32)
echo -e "${YELLOW}Session Secret:${NC}"
echo "SESSION_SECRET=$SESSION_SECRET"
echo ""

# Encryption Key (for sensitive data)
ENCRYPTION_KEY=$(generate_hex 32)
echo -e "${YELLOW}Encryption Key:${NC}"
echo "ENCRYPTION_KEY=$ENCRYPTION_KEY"
echo ""

echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}✓ Secrets generated successfully!${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "${YELLOW}IMPORTANT SECURITY NOTES:${NC}"
echo "1. Store these secrets in a secure password manager"
echo "2. Never commit these secrets to version control"
echo "3. Use different secrets for each environment"
echo "4. Rotate secrets regularly (every 90 days recommended)"
echo "5. Use environment variables or secret management tools"
echo ""
echo -e "${YELLOW}Recommended Secret Management Tools:${NC}"
echo "- HashiCorp Vault"
echo "- AWS Secrets Manager"
echo "- Google Cloud Secret Manager"
echo "- Azure Key Vault"
echo "- Doppler"
echo "- 1Password for Teams"
echo ""

# Optional: Save to a temporary file
read -p "Save secrets to a temporary file? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]
then
    TEMP_FILE="/tmp/karima-secrets-$(date +%s).txt"
    cat > "$TEMP_FILE" << EOF
# ============================================
# KARIMA STORE - Generated Secrets
# Generated: $(date)
# ============================================
# WARNING: Delete this file after copying secrets to secure storage!

# Database
DB_PASSWORD=$DB_PASSWORD

# Redis
REDIS_PASSWORD=$REDIS_PASSWORD

# JWT
JWT_SECRET=$JWT_SECRET

# API Key
API_KEY=$API_KEY

# Session
SESSION_SECRET=$SESSION_SECRET

# Encryption
ENCRYPTION_KEY=$ENCRYPTION_KEY

# ============================================
# NEXT STEPS:
# ============================================
# 1. Copy these secrets to your password manager
# 2. Update your .env.production file
# 3. DELETE THIS FILE: rm $TEMP_FILE
# 4. Verify secrets are not in git: git status
# ============================================
EOF
    echo -e "${GREEN}✓ Secrets saved to: $TEMP_FILE${NC}"
    echo -e "${RED}⚠ REMEMBER TO DELETE THIS FILE AFTER USE!${NC}"
    echo -e "   Run: ${YELLOW}rm $TEMP_FILE${NC}"
fi

echo ""
