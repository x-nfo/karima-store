# Variabel
BINARY_NAME=main
MAIN_PATH=./cmd/api/main.go

# --- MODE LOKAL (Hybrid) ---
# Gunakan ini saat develop kode agar cepat (DB di Podman, App di Terminal)
dev-local:
	@echo "Menjalankan Karima Store di port 8080..."
	APP_PORT=8080 DB_HOST=localhost REDIS_HOST=localhost REDIS_PORT=6380 go run $(MAIN_PATH)

# --- MODE PODMAN (Full Container) ---
# Bangun image dan jalankan semua layanan di kontainer
docker-up:
	@echo "Membangun dan menjalankan semua layanan di Podman..."
	podman-compose up -d --build

# Matikan semua layanan
docker-down:
	@echo "Menghentikan semua layanan..."
	podman-compose down

# --- KRATOS AUTH ---
# Jalankan Kratos Services
kratos-up:
	@echo "Menjalankan layanan Ory Kratos..."
	podman-compose -f docker-compose.yml -f docker-compose.kratos.yml up -d kratos-migrate kratos kratos-selfservice-ui-node mailslurper

# Matikan Kratos Services
kratos-down:
	@echo "Menghentikan layanan Ory Kratos..."
	podman-compose -f docker-compose.yml -f docker-compose.kratos.yml stop kratos-migrate kratos kratos-selfservice-ui-node mailslurper

# Lihat log aplikasi backend saja
logs:
	podman logs -f karima_store_backend

# --- UTILITY ---
# Merapikan library Go
tidy:
	go mod tidy
	go mod verify

# Masuk ke terminal database postgres
db-shell:
	podman exec -it karima_postgres psql -U karima_store -d karima_db

# Bersihkan image sampah (<none>)
clean:
	podman image prune -f

# --- SWAGGER ---
# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	$(HOME)/go/bin/swag init -g cmd/api/main.go

# --- ENVIRONMENT SETUP ---
# Generate secure secrets for production
generate-secrets:
	@echo "Generating secure secrets..."
	@chmod +x scripts/generate-env-secrets.sh
	@./scripts/generate-env-secrets.sh

# Verify environment configuration
verify-env:
	@echo "Verifying environment configuration..."
	@chmod +x scripts/verify-env.sh
	@./scripts/verify-env.sh .env.production

# Verify local environment
verify-env-local:
	@echo "Verifying local environment configuration..."
	@chmod +x scripts/verify-env.sh
	@./scripts/verify-env.sh .env.local

# Create production env from template
setup-prod-env:
	@echo "Setting up production environment..."
	@if [ -f .env.production ]; then \
		echo "⚠ .env.production already exists. Backup created as .env.production.backup"; \
		cp .env.production .env.production.backup; \
	fi
	@cp .env.production.template .env.production
	@chmod 600 .env.production
	@echo "✓ Created .env.production from template"
	@echo "⚠ Please edit .env.production and fill in all required values"
	@echo "  Then run: make verify-env"

# --- PRODUCTION DEPLOYMENT ---
# Deploy to production server
deploy-prod:
	@echo "Deploying to production..."
	@make verify-env
	@echo "✓ Environment verified"
	@echo "Copying files to production server..."
	@# Add your deployment commands here
	@echo "⚠ Configure your deployment commands in Makefile"

# Check production health
prod-health:
	@echo "Checking production health..."
	@curl -f http://localhost:8080/health || echo "⚠ Health check failed"

# View production logs
prod-logs:
	@echo "Viewing production logs..."
	@tail -f logs/app.log

# --- TESTING ---
# Run all tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report generated: coverage.html"

# --- HELP ---
# Show available commands
help:
	@echo "Karima Store - Available Commands:"
	@echo ""
	@echo "Development:"
	@echo "  make dev-local       - Run app locally (DB in Podman)"
	@echo "  make docker-up       - Run all services in containers"
	@echo "  make docker-down     - Stop all services"
	@echo ""
	@echo "Authentication (Kratos):"
	@echo "  make kratos-up       - Start Ory Kratos services"
	@echo "  make kratos-down     - Stop Ory Kratos services"
	@echo ""
	@echo "Environment Setup:"
	@echo "  make generate-secrets    - Generate secure passwords"
	@echo "  make setup-prod-env      - Create .env.production from template"
	@echo "  make verify-env          - Verify production environment"
	@echo "  make verify-env-local    - Verify local environment"
	@echo ""
	@echo "Production:"
	@echo "  make deploy-prod     - Deploy to production"
	@echo "  make prod-health     - Check production health"
	@echo "  make prod-logs       - View production logs"
	@echo ""
	@echo "Testing:"
	@echo "  make test            - Run all tests"
	@echo "  make test-coverage   - Run tests with coverage report"
	@echo ""
	@echo "Utilities:"
	@echo "  make tidy            - Clean up Go modules"
	@echo "  make db-shell        - Access database shell"
	@echo "  make logs            - View backend logs"
	@echo "  make clean           - Clean up Docker images"
	@echo "  make swagger         - Generate API documentation"
	@echo ""

.PHONY: dev-local docker-up docker-down kratos-up kratos-down logs tidy db-shell clean swagger \
        generate-secrets verify-env verify-env-local setup-prod-env deploy-prod prod-health \
        prod-logs test test-coverage help