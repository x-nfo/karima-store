# Karima Store

A modern, scalable e-commerce backend API built with Go, designed for high performance and reliability.

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Database Migrations](#database-migrations)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## 🎯 Overview

Karima Store is a comprehensive e-commerce backend solution that provides a robust foundation for building online stores. It features modular architecture, comprehensive authentication, payment processing, shipping integration, and more.

## ✨ Features

### Core Functionality
- **Product Management**: Create, read, update, and delete products with variants
- **Category Management**: Organize products into categories
- **Shopping Cart**: Add items to cart, update quantities, remove items
- **Checkout Process**: Complete order flow with payment integration
- **Order Management**: Track orders, update status, manage order history
- **User Authentication**: Secure authentication using Ory Kratos
- **Media Management**: Upload and manage product images (local or Cloudflare R2)
- **Coupon System**: Apply discount codes and manage promotions
- **Flash Sales**: Time-limited promotional campaigns
- **Reviews & Ratings**: Customer feedback system
- **Wishlist**: Save favorite products for later
- **Tax Management**: Configure tax rates and rules
- **Stock Management**: Track inventory levels and stock movements

### Integrations
- **Payment**: Midtrans payment gateway integration
- **Shipping**: RajaOngkir/Komerce shipping API
- **Notifications**: WhatsApp notifications via Fonnte
- **Storage**: Cloudflare R2 or local file storage
- **Email**: SMTP email configuration

### Technical Features
- **RESTful API**: Well-structured REST endpoints
- **API Documentation**: Auto-generated Swagger/OpenAPI documentation
- **Rate Limiting**: Protect against API abuse
- **CORS Support**: Cross-origin resource sharing configuration
- **Input Validation**: Request validation using struct tags
- **Database Migrations**: Version-controlled database schema
- **Caching**: Redis-based caching for improved performance
- **Docker Support**: Containerized deployment with Docker/Podman
- **Comprehensive Testing**: Unit and integration tests

## 🛠 Tech Stack

### Backend
- **Language**: Go 1.24.0
- **Web Framework**: Fiber v2 (high-performance HTTP framework)
- **ORM**: GORM (Go Object Relational Mapping)
- **Database**: PostgreSQL
- **Cache**: Redis

### Authentication
- **Identity Management**: Ory Kratos

### External Services
- **Payment Gateway**: Midtrans
- **Shipping**: RajaOngkir/Komerce
- **Notifications**: Fonnte (WhatsApp)
- **Storage**: Cloudflare R2 / Local
- **Email**: SMTP

### Development Tools
- **API Documentation**: Swagger/OpenAPI
- **Containerization**: Docker/Podman
- **Testing**: Go testing framework with testify
- **Database Migrations**: golang-migrate

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.24.0 or higher
- **PostgreSQL**: Version 12 or higher
- **Redis**: Version 6 or higher
- **Docker/Podman**: For containerized deployment (optional)
- **Make**: For running make commands

### Installation

#### macOS
```bash
brew install go postgresql redis
```

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install golang postgresql redis-server make
```

#### Windows
Download and install from:
- Go: https://golang.org/dl/
- PostgreSQL: https://www.postgresql.org/download/windows/
- Redis: https://redis.io/download

## 🚀 Installation

1. **Clone the repository**
```bash
git clone https://github.com/karima-store/karima-store.git
cd karima-store
```

2. **Install dependencies**
```bash
go mod download
```

3. **Set up environment variables**
```bash
cp .env.example .env.local
```

Edit `.env.local` with your configuration values (see [Configuration](#configuration) section).

4. **Set up the database**
```bash
# Start PostgreSQL and Redis (if using Docker)
podman-compose up -d postgres redis

# Run database migrations
make migrate
```

5. **Install Swagger CLI** (for generating API documentation)
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

## ⚙️ Configuration

The application uses environment variables for configuration. Copy [`.env.example`](.env.example) to `.env.local` (development) or `.env.production` (production) and configure the following:

### Server Configuration
```env
APP_PORT=8080
APP_ENV=development
API_VERSION=v1
```

### Database Configuration
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_secure_password
DB_NAME=karima_db
DB_SSL_MODE=disable
```

### Redis Configuration
```env
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

### Authentication (Ory Kratos)
```env
KRATOS_PUBLIC_URL=http://127.0.0.1:4433
KRATOS_ADMIN_URL=http://127.0.0.1:4434
KRATOS_UI_URL=http://127.0.0.1:4455
JWT_SECRET=your_super_secret_jwt_key
JWT_EXPIRATION=24h
```

### Payment Gateway (Midtrans)
```env
MIDTRANS_SERVER_KEY=YOUR_MIDTRANS_SERVER_KEY
MIDTRANS_CLIENT_KEY=YOUR_MIDTRANS_CLIENT_KEY
MIDTRANS_IS_PRODUCTION=false
MIDTRANS_API_BASE_URL=https://app.sandbox.midtrans.com/snap/v1
```

### Shipping (RajaOngkir/Komerce)
```env
RAJAONGKIR_API_KEY=YOUR_RAJAONGKIR_API_KEY
RAJAONGKIR_BASE_URL=https://api-sandbox.collaborator.komerce.id/tariff/api/v1/
```

### Storage (Cloudflare R2)
```env
FILE_STORAGE=local
FILE_UPLOAD_MAX_SIZE=10MB
R2_ACCOUNT_ID=your_r2_account_id
R2_ENDPOINT=https://your_account_id.r2.cloudflarestorage.com
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET_NAME=karima-media
R2_PUBLIC_URL=https://your-custom-domain.com
```

### Notifications (Fonnte)
```env
FONNTE_TOKEN=YOUR_FONNTE_TOKEN
FONNTE_URL=https://api.fonnte.com/send
```

### CORS Configuration
```env
# Development
CORS_ORIGIN=http://localhost:3000,http://localhost:8080

# Production (MUST be specific domains)
# CORS_ORIGIN=https://yourdomain.com,https://www.yourdomain.com
```

For a complete list of configuration options, see [`.env.example`](.env.example).

## 🏃 Running the Application

### Development Mode (Hybrid)
Run the application locally with database in Docker:

```bash
make dev-local
```

This starts the application on port 8080 with PostgreSQL and Redis running in Docker.

### Full Docker Mode
Run everything in containers:

```bash
make docker-up
```

To stop all services:

```bash
make docker-down
```

### With Ory Kratos Authentication
Start Kratos services along with the application:

```bash
make kratos-up
```

To stop Kratos services:

```bash
make kratos-down
```

### View Logs
View application logs:

```bash
make logs
```

### Database Shell
Access PostgreSQL shell:

```bash
make db-shell
```

## 📚 API Documentation

The API documentation is auto-generated using Swagger/OpenAPI.

### Generate Documentation
```bash
make swagger
```

### Access Documentation
Once the application is running, access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

The Swagger JSON is available at:
```
http://localhost:8080/swagger/doc.json
```

### API Standards
For detailed API standards and conventions, see [`docs/api_standards.md`](docs/api_standards.md).

## 📁 Project Structure

```
karima-store/
├── cmd/                    # Application entry points
│   ├── api/               # Main API application
│   ├── check_conn/        # Database connection checker
│   └── migrate/           # Database migration runner
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── database/         # Database connections and setup
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware (auth, CORS, rate limiting)
│   ├── models/           # Data models
│   ├── repository/       # Data access layer
│   ├── routes/           # Route definitions
│   ├── services/         # Business logic layer
│   ├── storage/          # Storage implementations (R2, local)
│   └── utils/            # Utility functions
├── migrations/           # Database migration files
├── docs/                 # Documentation
│   ├── swagger/          # Swagger documentation
│   ├── QA/               # QA reports and testing plans
│   └── Progress_flows/   # Workflow documentation
├── public/               # Public assets
├── scripts/              # Utility scripts
├── deploy/               # Deployment configurations
│   └── kratos/          # Ory Kratos configuration
├── .env.example          # Environment variables template
├── docker-compose.yml    # Main Docker Compose configuration
├── docker-compose.kratos.yml  # Kratos Docker Compose configuration
├── Dockerfile            # Application Docker image
├── Makefile              # Make commands
└── go.mod                # Go module dependencies
```

### Architecture Layers

1. **Handlers**: Handle HTTP requests and responses
2. **Services**: Contain business logic
3. **Repository**: Data access and database operations
4. **Models**: Data structures and database entities
5. **Middleware**: Cross-cutting concerns (auth, validation, logging)

For more details on architecture, see [`docs/architecture_path.md`](docs/architecture_path.md).

## 🗄 Database Migrations

The project uses `golang-migrate` for database version control.

### Running Migrations
```bash
# Run all pending migrations
make migrate

# Or run manually
go run cmd/migrate/main.go up
```

### Migration Files
Migration files are located in the [`migrations/`](migrations/) directory:
- `000001_init_schema.up.sql` - Initial schema
- `000002_add_pricing_enhancements.up.sql` - Pricing features
- `000003_add_media_table.up.sql` - Media management
- `000004_add_media_deleted_at.up.sql` - Soft delete for media
- `000005_fix_media_schema.up.sql` - Media schema fixes
- `000006_add_kratos_identity.up.sql` - Kratos integration

### Rolling Back Migrations
```bash
go run cmd/migrate/main.go down 1
```

For database schema documentation, see [`docs/system_erd.md`](docs/system_erd.md).

## 🧪 Testing

The project includes comprehensive unit and integration tests.

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for specific package
go test ./internal/services/

# Run tests with verbose output
go test -v ./...
```

### Test Structure
Test files follow the convention `*_test.go` and are located alongside the code they test:
- [`internal/services/product_service_test.go`](internal/services/product_service_test.go)
- [`internal/services/checkout_service_test.go`](internal/services/checkout_service_test.go)
- [`internal/services/media_service_test.go`](internal/services/media_service_test.go)
- [`internal/middleware/cors_test.go`](internal/middleware/cors_test.go)
- [`internal/middleware/kratos_test.go`](internal/middleware/kratos_test.go)
- [`internal/middleware/rate_limit_test.go`](internal/middleware/rate_limit_test.go)
- [`internal/middleware/validator_test.go`](internal/middleware/validator_test.go)

### Test Setup
The project uses a test setup file [`internal/test_setup.go`](internal/test_setup.go) for common test utilities.

For QA reports and testing plans, see the [`docs/QA/`](docs/QA/) directory.

## 🚢 Deployment

### Docker Deployment

Build and run with Docker Compose:

```bash
# Build and start all services
make docker-up

# View logs
make logs

# Stop services
make docker-down
```

### Environment-Specific Configuration

Use different environment files for different environments:
- `.env.local` - Local development
- `.env.production` - Production deployment

### Production Considerations

1. **Security**
   - Use strong, unique passwords
   - Set `APP_ENV=production`
   - Configure proper CORS origins
   - Use HTTPS in production
   - Enable SSL for database connections

2. **Performance**
   - Enable Redis caching
   - Configure appropriate rate limits
   - Use CDN for static assets
   - Enable database connection pooling

3. **Monitoring**
   - Set up log aggregation
   - Monitor application metrics
   - Set up alerts for errors

4. **Backups**
   - Regular database backups
   - Backup R2 storage if using cloud storage
   - Document disaster recovery procedures

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Code Style
- Follow Go coding standards
- Use meaningful variable and function names
- Add comments for complex logic
- Write tests for new features
- Update documentation as needed

### Naming Conventions
See [`docs/naming_convention.md`](docs/naming_convention.md) for project-specific naming guidelines.

## 📖 Additional Documentation

- [API Standards](docs/api_standards.md) - API design guidelines
- [Architecture Path](docs/architecture_path.md) - System architecture overview
- [System ERD](docs/system_erd.md) - Database schema documentation
- [Naming Convention](docs/naming_convention.md) - Naming guidelines
- [Master Playbook](docs/master_playbook.md) - Development playbook
- [QA Reports](docs/QA/) - Quality assurance reports and plans

## 🐛 Troubleshooting

### Database Connection Issues
```bash
# Check if PostgreSQL is running
podman ps | grep postgres

# Check database logs
podman logs karima_postgres

# Test connection
make check_conn
```

### Redis Connection Issues
```bash
# Check if Redis is running
podman ps | grep redis

# Test Redis connection
redis-cli ping
```

### Migration Issues
```bash
# Check migration status
go run cmd/migrate/main.go version

# Force specific version
go run cmd/migrate/main.go force 000001
```

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

## 📞 Support

For issues, questions, or contributions:
- Open an issue on GitHub
- Check existing documentation in the [`docs/`](docs/) directory
- Review QA reports in [`docs/QA/`](docs/QA/)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Built with [Fiber](https://gofiber.io/) web framework
- Authentication powered by [Ory Kratos](https://www.ory.sh/kratos/)
- Database ORM by [GORM](https://gorm.io/)
- API documentation with [Swagger](https://swagger.io/)

---

**Built with ❤️ for modern e-commerce**
