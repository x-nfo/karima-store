#!/bin/bash

# ============================================
# Karima Store - VPS Initial Setup Script
# ============================================
# This script prepares a fresh VPS for production deployment
# 
# Usage: 
#   1. Copy to VPS: scp scripts/setup-vps.sh user@vps:~/
#   2. SSH to VPS: ssh user@vps
#   3. Run: chmod +x setup-vps.sh && sudo ./setup-vps.sh
#
# Requirements:
#   - Ubuntu 22.04 LTS or newer
#   - Root or sudo access
# ============================================

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}============================================${NC}"
echo -e "${BLUE}Karima Store - VPS Setup${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo -e "${RED}This script must be run as root${NC}" 
   exit 1
fi

# Variables
APP_USER="karima"
APP_DIR="/opt/karima-store"
DB_NAME="karima_db_prod"
DB_USER="karima_prod_user"

echo -e "${YELLOW}This script will:${NC}"
echo "1. Update system packages"
echo "2. Install required dependencies"
echo "3. Setup PostgreSQL database"
echo "4. Setup Redis"
echo "5. Configure firewall"
echo "6. Create application user and directories"
echo "7. Install and configure Nginx"
echo "8. Setup SSL with Let's Encrypt"
echo ""
read -p "Continue? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    exit 1
fi

# ============================================
# 1. System Update
# ============================================
echo -e "${BLUE}[1/8] Updating system packages...${NC}"
apt update
apt upgrade -y
apt install -y curl wget git unzip software-properties-common

# ============================================
# 2. Install Dependencies
# ============================================
echo -e "${BLUE}[2/8] Installing dependencies...${NC}"

# PostgreSQL
echo -e "${YELLOW}Installing PostgreSQL...${NC}"
apt install -y postgresql postgresql-contrib

# Redis
echo -e "${YELLOW}Installing Redis...${NC}"
apt install -y redis-server

# Nginx
echo -e "${YELLOW}Installing Nginx...${NC}"
apt install -y nginx

# Certbot for SSL
echo -e "${YELLOW}Installing Certbot...${NC}"
apt install -y certbot python3-certbot-nginx

# Monitoring tools
apt install -y htop iotop nethogs

echo -e "${GREEN}✓ Dependencies installed${NC}"

# ============================================
# 3. Setup PostgreSQL
# ============================================
echo -e "${BLUE}[3/8] Configuring PostgreSQL...${NC}"

# Generate random password
DB_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)

# Create database and user
sudo -u postgres psql <<EOF
CREATE DATABASE $DB_NAME;
CREATE USER $DB_USER WITH ENCRYPTED PASSWORD '$DB_PASSWORD';
GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;
ALTER DATABASE $DB_NAME OWNER TO $DB_USER;
\q
EOF

# Configure PostgreSQL for remote connections (if needed)
PG_VERSION=$(psql --version | awk '{print $3}' | cut -d. -f1)
PG_CONF="/etc/postgresql/$PG_VERSION/main/postgresql.conf"
PG_HBA="/etc/postgresql/$PG_VERSION/main/pg_hba.conf"

# Enable SSL
sed -i "s/#ssl = off/ssl = on/" $PG_CONF

# Restart PostgreSQL
systemctl restart postgresql
systemctl enable postgresql

echo -e "${GREEN}✓ PostgreSQL configured${NC}"
echo -e "${YELLOW}Database credentials:${NC}"
echo "  DB_NAME: $DB_NAME"
echo "  DB_USER: $DB_USER"
echo "  DB_PASSWORD: $DB_PASSWORD"
echo -e "${RED}⚠ Save these credentials securely!${NC}"
echo ""

# ============================================
# 4. Setup Redis
# ============================================
echo -e "${BLUE}[4/8] Configuring Redis...${NC}"

# Generate Redis password
REDIS_PASSWORD=$(openssl rand -base64 32 | tr -d "=+/" | cut -c1-32)

# Configure Redis
REDIS_CONF="/etc/redis/redis.conf"
cp $REDIS_CONF ${REDIS_CONF}.backup

# Set password
sed -i "s/# requirepass foobared/requirepass $REDIS_PASSWORD/" $REDIS_CONF

# Bind to localhost only
sed -i "s/bind 127.0.0.1 ::1/bind 127.0.0.1/" $REDIS_CONF

# Enable persistence
sed -i "s/appendonly no/appendonly yes/" $REDIS_CONF

# Restart Redis
systemctl restart redis-server
systemctl enable redis-server

echo -e "${GREEN}✓ Redis configured${NC}"
echo -e "${YELLOW}Redis password: $REDIS_PASSWORD${NC}"
echo -e "${RED}⚠ Save this password securely!${NC}"
echo ""

# ============================================
# 5. Configure Firewall
# ============================================
echo -e "${BLUE}[5/8] Configuring firewall...${NC}"

# Install UFW if not present
apt install -y ufw

# Default policies
ufw default deny incoming
ufw default allow outgoing

# Allow SSH (important!)
ufw allow 22/tcp

# Allow HTTP and HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# Enable firewall
echo "y" | ufw enable

echo -e "${GREEN}✓ Firewall configured${NC}"

# ============================================
# 6. Create Application User and Directories
# ============================================
echo -e "${BLUE}[6/8] Creating application user and directories...${NC}"

# Create user
if ! id "$APP_USER" &>/dev/null; then
    useradd -r -s /bin/bash -d $APP_DIR -m $APP_USER
    echo -e "${GREEN}✓ User $APP_USER created${NC}"
else
    echo -e "${YELLOW}⚠ User $APP_USER already exists${NC}"
fi

# Create directories
mkdir -p $APP_DIR/{logs,migrations,static}
chown -R $APP_USER:$APP_USER $APP_DIR
chmod 755 $APP_DIR

echo -e "${GREEN}✓ Directories created${NC}"

# ============================================
# 7. Configure Nginx
# ============================================
echo -e "${BLUE}[7/8] Configuring Nginx...${NC}"

# Remove default site
rm -f /etc/nginx/sites-enabled/default

# Create basic configuration (will be replaced with proper config later)
cat > /etc/nginx/sites-available/karima-store <<'EOF'
server {
    listen 80;
    server_name _;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

# Enable site
ln -sf /etc/nginx/sites-available/karima-store /etc/nginx/sites-enabled/

# Test configuration
nginx -t

# Restart Nginx
systemctl restart nginx
systemctl enable nginx

echo -e "${GREEN}✓ Nginx configured${NC}"

# ============================================
# 8. Setup SSL (Optional - requires domain)
# ============================================
echo -e "${BLUE}[8/8] SSL Setup${NC}"
echo -e "${YELLOW}To setup SSL, you need a domain name pointing to this server${NC}"
echo ""
read -p "Do you have a domain configured? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    read -p "Enter your domain (e.g., api.yourdomain.com): " DOMAIN
    read -p "Enter your email for Let's Encrypt: " EMAIL
    
    # Get SSL certificate
    certbot --nginx -d $DOMAIN --non-interactive --agree-tos -m $EMAIL
    
    echo -e "${GREEN}✓ SSL certificate installed${NC}"
else
    echo -e "${YELLOW}⚠ Skipping SSL setup. You can run this later:${NC}"
    echo "  certbot --nginx -d yourdomain.com -m your@email.com"
fi

# ============================================
# Summary
# ============================================
echo ""
echo -e "${BLUE}============================================${NC}"
echo -e "${GREEN}✓ VPS Setup Complete!${NC}"
echo -e "${BLUE}============================================${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo "1. Save these credentials to your password manager:"
echo "   - PostgreSQL: $DB_USER / $DB_PASSWORD"
echo "   - Redis: $REDIS_PASSWORD"
echo ""
echo "2. Update your .env.production file with these values"
echo ""
echo "3. Deploy your application:"
echo "   - Copy binary: scp karima-store $APP_USER@server:$APP_DIR/"
echo "   - Copy .env: scp .env.production $APP_USER@server:$APP_DIR/"
echo "   - Copy migrations: scp -r migrations $APP_USER@server:$APP_DIR/"
echo ""
echo "4. Setup systemd service:"
echo "   - Copy service file: sudo cp deploy/systemd/karima-store.service /etc/systemd/system/"
echo "   - Enable: sudo systemctl enable karima-store"
echo "   - Start: sudo systemctl start karima-store"
echo ""
echo "5. Update Nginx configuration:"
echo "   - Copy config: sudo cp deploy/nginx/karima-store.conf /etc/nginx/sites-available/karima-store"
echo "   - Test: sudo nginx -t"
echo "   - Reload: sudo systemctl reload nginx"
echo ""
echo -e "${YELLOW}Application Directory: $APP_DIR${NC}"
echo -e "${YELLOW}Application User: $APP_USER${NC}"
echo ""
echo -e "${GREEN}Happy Deploying! 🚀${NC}"
echo ""

# Save credentials to file
CREDS_FILE="/root/karima-credentials.txt"
cat > $CREDS_FILE <<EOF
Karima Store - Server Credentials
Generated: $(date)

PostgreSQL:
  Database: $DB_NAME
  User: $DB_USER
  Password: $DB_PASSWORD
  Host: localhost
  Port: 5432

Redis:
  Host: localhost
  Port: 6379
  Password: $REDIS_PASSWORD

Application:
  User: $APP_USER
  Directory: $APP_DIR

IMPORTANT: Delete this file after saving credentials securely!
Command: rm $CREDS_FILE
EOF

chmod 600 $CREDS_FILE
echo -e "${YELLOW}Credentials saved to: $CREDS_FILE${NC}"
echo -e "${RED}⚠ Remember to delete this file after copying credentials!${NC}"
