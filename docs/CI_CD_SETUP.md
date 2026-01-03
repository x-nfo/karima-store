# CI/CD Pipeline Documentation

## Overview
Automated CI/CD pipeline using GitHub Actions for testing, building, and deploying Karima Store backend.

## Workflows

### 1. CI - Continuous Integration (`.github/workflows/ci.yml`)
**Triggers:** Push to main/develop, Pull Requests

**Jobs:**
- **Lint & Format Check**
  - gofmt validation
  - go vet analysis
  - golangci-lint with multiple linters
  
- **Unit Tests**
  - Runs all tests with race detection
  - Generates coverage reports
  - Uses PostgreSQL and Redis services
  - Uploads coverage to Codecov
  
- **Build**
  - Builds Docker image
  - Validates Dockerfile
  - Uses build cache for speed
  
- **Security**
  - Gosec security scanner
  - Trivy vulnerability scanner
  - Uploads results to GitHub Security tab

### 2. CD - Continuous Deployment (`.github/workflows/cd.yml`)
**Triggers:** Push to main branch

**Jobs:**
- **Deploy to Production**
  - Builds and pushes Docker image to Docker Hub
  - Deploys to VPS via SSH
  - Performs health check
  - Notifies on success/failure

### 3. Test Coverage (`.github/workflows/test-coverage.yml`)
**Triggers:** Pull Requests

**Jobs:**
- **Coverage Report**
  - Generates detailed coverage report
  - Comments on PR with coverage stats
  - Enforces minimum 20% coverage threshold

## Required GitHub Secrets

Configure these secrets in your GitHub repository settings:

### Docker Hub
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub password/token

### VPS Deployment
- `VPS_HOST` - VPS server hostname/IP
- `VPS_USERNAME` - SSH username
- `VPS_SSH_KEY` - Private SSH key for authentication
- `VPS_PORT` - SSH port (default: 22)

## Setup Instructions

### 1. Configure Secrets
```bash
# Go to GitHub repository
# Settings → Secrets and variables → Actions → New repository secret
# Add all required secrets listed above
```

### 2. Enable GitHub Actions
- GitHub Actions should be enabled by default
- Workflows will run automatically on push/PR

### 3. VPS Setup
On your VPS server:

```bash
# Create app directory
sudo mkdir -p /opt/karima-store
cd /opt/karima-store

# Clone repository
git clone https://github.com/x-nfo/karima-store.git .

# Create docker-compose.yml for production
# (Copy from repository)

# Make deploy script executable
chmod +x scripts/deploy.sh

# Initial deployment
./scripts/deploy.sh
```

### 4. Docker Hub Setup
```bash
# Login to Docker Hub
docker login

# Tag and push initial image
docker build -t yourusername/karima-store:latest .
docker push yourusername/karima-store:latest
```

## Linting Configuration

`.golangci.yml` configures the following linters:
- errcheck - Check for unchecked errors
- gosimple - Simplify code
- govet - Vet examines Go source code
- ineffassign - Detect ineffectual assignments
- staticcheck - Static analysis
- unused - Check for unused code
- gofmt - Format checking
- goimports - Import checking
- misspell - Spell checking
- goconst - Find repeated strings
- gocyclo - Cyclomatic complexity
- revive - Fast, configurable linter
- gosec - Security audit

## Deployment Script

`scripts/deploy.sh` performs:
1. Database backup before deployment
2. Pull latest Docker images
3. Stop current containers
4. Start new containers
5. Health check
6. Rollback on failure
7. Cleanup old images and backups

## Monitoring

### CI Status
- Check workflow runs in GitHub Actions tab
- View test results and coverage
- Review security scan results

### Deployment Status
- Monitor deployment logs in Actions
- Check VPS health endpoint: `http://your-vps:8080/api/v1/health`
- Review application logs: `docker-compose logs -f`

## Troubleshooting

### CI Failures
```bash
# Run tests locally
go test ./...

# Run linter locally
golangci-lint run

# Check formatting
gofmt -s -l .
```

### Deployment Failures
```bash
# SSH to VPS
ssh user@your-vps

# Check container status
docker-compose ps

# View logs
docker-compose logs

# Manual rollback
cd /opt/karima-store
docker-compose down
# Restore from backup if needed
```

### Coverage Issues
```bash
# Generate coverage locally
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Best Practices

1. **Always create PRs** - Don't push directly to main
2. **Wait for CI** - Ensure all checks pass before merging
3. **Review coverage** - Maintain minimum 20% coverage
4. **Monitor deployments** - Check health after each deployment
5. **Keep secrets secure** - Never commit secrets to repository

## Maintenance

### Update Dependencies
```bash
go get -u ./...
go mod tidy
```

### Update Workflows
- Edit workflow files in `.github/workflows/`
- Test changes on a feature branch first
- Monitor first run after changes

### Backup Management
- Backups stored in `/opt/backups/karima-store/`
- Automatically keeps last 7 backups
- Manual cleanup if needed: `rm /opt/backups/karima-store/old_backup.sql`
