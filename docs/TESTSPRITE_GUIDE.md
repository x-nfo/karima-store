# TestSprite Integration Guide

## Overview

This guide explains how to use TestSprite to test the Karima Store e-commerce backend. TestSprite provides a comprehensive testing framework that integrates with Go's native testing tools to provide enhanced test reporting, coverage analysis, and CI/CD integration.

## Prerequisites

Before using TestSprite with this project, ensure you have:

- Go 1.24.0 or higher installed
- PostgreSQL and Redis running (either locally or via Docker/Podman)
- Make utility installed
- TestSprite CLI installed (if using TestSprite directly)

## Installation

### Install TestSprite CLI (Optional)

If you want to use TestSprite's advanced features:

```bash
# Install TestSprite CLI
npm install -g testsprite

# Or using Go
go install github.com/testsprite/testsprite@latest
```

### Install Additional Tools

Some TestSprite features require additional tools:

```bash
# Install gotestsum for JUnit XML output
go install gotest.tools/gotestsum@latest

# Install reflex for watch mode
go install github.com/cespare/reflex@latest
```

## Configuration

TestSprite configuration is defined in [`.testsprite.json`](../.testsprite.json). This file includes:

- **Test Patterns**: Glob patterns for test files
- **Test Suites**: Organized test groups (unit, handler, middleware, integration)
- **Coverage Settings**: Coverage thresholds and report formats
- **Pre/Post Test Hooks**: Database setup and cleanup
- **CI Integration**: GitHub Actions, GitLab CI, Jenkins support
- **Flaky Test Detection**: Automatic retry and detection of flaky tests

## Quick Start

### Run All Tests with TestSprite

```bash
# Run all tests with TestSprite configuration
make test-sprite

# Run tests with coverage report
make test-sprite-coverage

# Run full test suite with all reports
make test-sprite-full
```

### Run Specific Test Suites

```bash
# Run only unit tests
make test-sprite-unit

# Run only handler tests
make test-sprite-handler

# Run only middleware tests
make test-sprite-middleware
```

### Watch Mode

```bash
# Run tests automatically when files change
make test-sprite-watch
```

## Test Suites

### 1. Unit Tests

Tests for individual components in isolation:

- **Services**: Business logic tests
  - [`internal/services/product_service_test.go`](../internal/services/product_service_test.go)
  - [`internal/services/checkout_service_test.go`](../internal/services/checkout_service_test.go)
  - [`internal/services/media_service_test.go`](../internal/services/media_service_test.go)
  - [`internal/services/pricing_service_test.go`](../internal/services/pricing_service_test.go)

- **Models**: Data model tests
  - [`internal/models/product_test.go`](../internal/models/product_test.go)

- **Repository**: Data access layer tests
  - [`internal/repository/product_repository_test.go`](../internal/repository/product_repository_test.go)

### 2. Handler Tests

HTTP endpoint and request handler tests:

- [`internal/handlers/product_handler_test.go`](../internal/handlers/product_handler_test.go)
- [`internal/handlers/category_handler_test.go`](../internal/handlers/category_handler_test.go)
- [`internal/handlers/checkout_handler_test.go`](../internal/handlers/checkout_handler_test.go)
- [`internal/handlers/media_handler_test.go`](../internal/handlers/media_handler_test.go)
- [`internal/handlers/order_handler_test.go`](../internal/handlers/order_handler_test.go)
- [`internal/handlers/pricing_handler_test.go`](../internal/handlers/pricing_handler_test.go)
- [`internal/handlers/variant_handler_test.go`](../internal/handlers/variant_handler_test.go)
- [`internal/handlers/whatsapp_handler_test.go`](../internal/handlers/whatsapp_handler_test.go)
- [`internal/handlers/komerce_handler_test.go`](../internal/handlers/komerce_handler_test.go)

### 3. Middleware Tests

Middleware component tests:

- [`internal/middleware/api_key_test.go`](../internal/middleware/api_key_test.go)
- [`internal/middleware/cors_test.go`](../internal/middleware/cors_test.go)
- [`internal/middleware/csrf_test.go`](../internal/middleware/csrf_test.go)
- [`internal/middleware/error_handler_test.go`](../internal/middleware/error_handler_test.go)
- [`internal/middleware/kratos_test.go`](../internal/middleware/kratos_test.go)
- [`internal/middleware/rate_limit_test.go`](../internal/middleware/rate_limit_test.go)
- [`internal/middleware/security_test.go`](../internal/middleware/security_test.go)
- [`internal/middleware/validator_test.go`](../internal/middleware/validator_test.go)

### 4. Integration Tests

Tests that verify interactions between multiple components:

- Full request/response cycles
- Database operations with real connections
- Redis caching integration
- External API mocking

## Test Database Setup

TestSprite automatically handles test database setup and cleanup:

```bash
# Setup test database
make db-setup

# Run migrations
make migrate-test

# Cleanup test database
make db-cleanup
```

These commands are automatically run as part of `make test-sprite-full`.

## Coverage Reports

TestSprite generates comprehensive coverage reports:

### Generate Coverage Report

```bash
# Generate HTML coverage report
make test-sprite-coverage

# View coverage report in browser
make view-coverage
```

### Coverage Thresholds

The project enforces a minimum coverage threshold of 60% as defined in [`.testsprite.json`](../.testsprite.json). If coverage falls below this threshold, tests will fail.

### Coverage Reports

- **HTML Report**: `coverage.html` - Interactive coverage visualization
- **JSON Report**: `coverage.out` - Machine-readable coverage data
- **Terminal Output**: Summary displayed in console

## Test Results

### JUnit XML Output

Generate JUnit XML for CI/CD integration:

```bash
make test-sprite-junit
```

Output: `test-results/junit.xml`

### Test Report Format

TestSprite generates reports in multiple formats:

- **JUnit XML**: For CI/CD systems
- **JSON**: For programmatic analysis
- **HTML**: For human-readable reports

## CI/CD Integration

### GitHub Actions

TestSprite includes GitHub Actions integration:

```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24.0'
      - name: Run TestSprite tests
        run: make test-sprite-full
      - name: Upload coverage
        uses: actions/upload-artifact@v3
        with:
          name: coverage
          path: coverage.html
```

### GitLab CI

```yaml
# .gitlab-ci.yml
test:
  image: golang:1.24
  script:
    - make test-sprite-full
  artifacts:
    paths:
      - coverage.html
      - test-results/
```

### Jenkins

```groovy
// Jenkinsfile
pipeline {
    agent any
    stages {
        stage('Test') {
            steps {
                sh 'make test-sprite-full'
            }
        }
    }
    post {
        always {
            publishHTML(target: [
                reportDir: '.',
                reportFiles: 'coverage.html',
                reportName: 'Coverage Report'
            ])
        }
    }
}
```

## Flaky Test Detection

TestSprite automatically detects flaky tests by:

1. Running tests up to 3 times if they fail
2. Tracking test results across multiple runs
3. Flagging tests that fail intermittently

Configuration in [`.testsprite.json`](../.testsprite.json):

```json
{
  "flakyTestDetection": {
    "enabled": true,
    "maxRetries": 3,
    "threshold": 2
  }
}
```

## Performance Tracking

TestSprite tracks test execution time:

- **Slow Test Threshold**: 5 seconds
- **Execution Time Tracking**: Enabled by default
- **Alerts**: Notifies on slow tests

View performance metrics in the TestSprite dashboard.

## TestSprite Dashboard

Enable the TestSprite dashboard for real-time test monitoring:

```json
{
  "dashboard": {
    "enabled": true,
    "refreshInterval": "30s",
    "showTrends": true,
    "showFlakyTests": true
  }
}
```

Start the dashboard:

```bash
testsprite dashboard
```

## Best Practices

### 1. Test Organization

- Keep tests alongside the code they test
- Use descriptive test names
- Group related tests in test suites
- Follow Go testing conventions

### 2. Test Isolation

- Each test should be independent
- Clean up test data after each test
- Use test fixtures for common setup
- Avoid shared state between tests

### 3. Coverage Goals

- Aim for >80% coverage on critical paths
- Maintain >60% overall coverage
- Focus on business logic and handlers
- Don't test generated code or simple getters/setters

### 4. Performance

- Keep tests fast (<5s per test)
- Use table-driven tests for similar cases
- Mock external dependencies
- Run integration tests separately

### 5. CI/CD

- Run tests on every commit
- Block merges on test failures
- Upload coverage reports
- Track flaky tests

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
podman ps | grep postgres

# Start PostgreSQL
podman-compose up -d postgres

# Check database logs
podman logs karima_postgres
```

### Redis Connection Issues

```bash
# Check if Redis is running
podman ps | grep redis

# Start Redis
podman-compose up -d redis

# Test Redis connection
redis-cli ping
```

### Test Failures

```bash
# Run tests with verbose output
make test-sprite

# Run specific test
go test -v ./internal/services/ -run TestProductService

# Run tests with race detection
go test -race ./...
```

### Coverage Issues

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage details
go tool cover -func=coverage.out

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html
```

## Advanced Usage

### Custom Test Configurations

Create custom TestSprite configurations for different scenarios:

```bash
# Create custom config
cp .testsprite.json .testsprite.ci.json

# Edit for CI environment
# Then run with custom config
testsprite run --config .testsprite.ci.json
```

### Parallel Test Execution

TestSprite supports parallel test execution:

```json
{
  "parallel": {
    "enabled": true,
    "maxWorkers": 4
  }
}
```

### Test Filtering

Run specific tests using patterns:

```bash
# Run tests matching pattern
go test -v ./... -run TestProduct

# Run tests excluding pattern
go test -v ./... -skip TestIntegration

# Run tests in specific package
go test -v ./internal/services/...
```

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [TestSprite Documentation](https://docs.testsprite.io)
- [Testify Documentation](https://github.com/stretchr/testify)
- [Project README](../README.md)
- [API Standards](./api_standards.md)
- [Architecture Path](./architecture_path.md)

## Support

For issues or questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review existing test files for examples
3. Consult the [QA Reports](./QA/) directory
4. Open an issue on GitHub

## Summary

TestSprite provides a comprehensive testing solution for the Karima Store backend with:

- ✅ Automated test execution and reporting
- ✅ Coverage analysis with HTML reports
- ✅ CI/CD integration for multiple platforms
- ✅ Flaky test detection and retry logic
- ✅ Performance tracking and slow test alerts
- ✅ Real-time test monitoring dashboard
- ✅ JUnit XML output for test result aggregation

Use the provided Make commands to quickly run tests and generate reports. Configure [`.testsprite.json`](../.testsprite.json) to customize test behavior to your needs.
