# Security Testing Plan - Karima Store

## Executive Summary

This comprehensive security testing plan outlines the unit testing strategy for the Karima Store backend application. The plan focuses on ensuring full security coverage across all critical components, including authentication, authorization, input validation, rate limiting, CORS, and error handling.

## 1. Testing Scope and Objectives

### 1.1 Scope
- Authentication and Authorization middleware (Ory Kratos)
- Input validation and sanitization
- Rate limiting implementation
- CORS configuration
- Error handling and security headers
- File upload security
- Database security and transactions
- Session management

### 1.2 Objectives
- Identify and mitigate security vulnerabilities
- Ensure proper authentication and authorization controls
- Validate input validation effectiveness
- Test rate limiting and abuse prevention
- Verify CORS security configuration
- Ensure proper error handling and information disclosure prevention
- Validate file upload security measures

## 2. Testing Strategy

### 2.1 Approach
- **Unit Testing**: Test individual components in isolation
- **Integration Testing**: Test component interactions
- **Security Testing**: Focus on attack vectors and vulnerabilities
- **Boundary Testing**: Test edge cases and invalid inputs
- **Performance Testing**: Verify rate limiting effectiveness

### 2.2 Test Environment
- Development environment with test database
- Mock external services (Kratos, Komerce, Midtrans)
- Test data with realistic scenarios
- Security testing tools and frameworks

## 3. Test Cases

### 3.1 Authentication & Authorization Testing

#### 3.1.1 Ory Kratos Authentication Middleware
```go
// Test cases for Ory Kratos authentication
func TestKratosAuthentication(t *testing.T) {
    // Test valid Kratos session cookie
    // Test expired Kratos session
    // Test invalid session cookie
    // Test missing session cookie
    // Test session validation with Kratos API
    // Test Bearer token authentication
    // Test token validation with invalid token
    // Test session validation failures
}

// Test cases for Kratos role-based access
func TestKratosRoleAccess(t *testing.T) {
    // Test admin role access via Kratos
    // Test user role access restrictions
    // Test role validation with different user roles
    // Test session validation failures
    // Test token-based authentication for API clients
}
```

### 3.2 Input Validation & Sanitization Testing

#### 3.2.1 Validation Middleware
```go
// Test cases for input validation
func TestInputValidation(t *testing.T) {
    // Test valid product creation input
    // Test invalid product data (missing required fields)
    // Test invalid data types (string instead of number)
    // Test SQL injection attempts
    // Test XSS attack vectors
    // Test command injection attempts
    // Test path traversal attempts
    // Test email validation
    // Test numeric range validation
    // Test string length validation
}

// Test cases for request body parsing
func TestRequestBodyParsing(t *testing.T) {
    // Test valid JSON request body
    // Test malformed JSON
    // Test empty request body
    // Test oversized request body
    // Test invalid content type
}
```

#### 3.2.2 File Upload Security
```go
// Test cases for file upload security
func TestFileUploadSecurity(t *testing.T) {
    // Test valid image file upload
    // Test file type validation (only images allowed)
    // Test file size limits
    // Test malicious file upload (script files)
    // Test file extension spoofing
    // Test magic bytes validation
    // Test directory traversal in filenames
    // Test concurrent file uploads
}

// Test cases for media validation
func TestMediaValidation(t *testing.T) {
    // Test image dimension validation
    // Test file content validation
    // Test duplicate file uploads
    // Test file storage security
}
```

### 3.3 Rate Limiting & CORS Testing

#### 3.3.1 Rate Limiting Middleware
```go
// Test cases for rate limiting
func TestRateLimiting(t *testing.T) {
    // Test normal request rate (within limits)
    // Test exceeding request rate (should be blocked)
    // Test rate limit reset functionality
    // Test different IP address rate limiting
    // Test production vs development rate limits
    // Test Redis-backed rate limiting
    // Test rate limit bypass attempts
    // Test concurrent request rate limiting
}

// Test cases for rate limit configuration
func TestRateLimitConfiguration(t *testing.T) {
    // Test environment-specific rate limits
    // Test custom rate limit configuration
    // Test rate limit window configuration
    // Test rate limit override functionality
}
```

#### 3.3.2 CORS Configuration
```go
// Test cases for CORS middleware
func TestCORSConfiguration(t *testing.T) {
    // Test allowed origin headers
    // Test allowed methods
    // Test allowed headers
    // Test credentials support
    // Test preflight request handling
    // Test CORS with authentication
    // Test CORS with different origins
    // Test CORS security headers
}

// Test cases for CORS security
func TestCORSSecurity(t *testing.T) {
    // Test cross-origin request blocking
    // Test sensitive data exposure via CORS
    // Test CORS misconfiguration attacks
    // Test CORS with authentication tokens
}
```

### 3.4 Error Handling & Security Headers

#### 3.4.1 Error Handling
```go
// Test cases for error handling
func TestErrorHandling(t *testing.T) {
    // Test generic error responses
    // Test validation error responses
    // Test authentication error responses
    // Test authorization error responses
    // Test server error responses
    // Test error message sanitization
    // Test error logging
    // Test error code consistency
    // Test error detail levels (dev vs production)
}

// Test cases for security error messages
func TestSecurityErrorMessages(t *testing.T) {
    // Test no sensitive information in error messages
    // Test consistent error formatting
    // Test error message localization
    // Test error message obfuscation
}
```

#### 3.4.2 Security Headers
```go
// Test cases for security headers
func TestSecurityHeaders(t *testing.T) {
    // Test HSTS header presence and configuration
    // Test CSP header configuration
    // Test X-Frame-Options header
    // Test X-Content-Type-Options header
    // Test X-XSS-Protection header
    // Test Referrer-Policy header
    // Test Content-Security-Policy configuration
    // Test security header inheritance
}

// Test cases for helmet middleware
func TestHelmetMiddleware(t *testing.T) {
    // Test all security headers are present
    // Test header values are correct
    // Test header configuration options
    // Test security header effectiveness
}
```

### 3.5 Database & Session Security

#### 3.5.1 Database Security
```go
// Test cases for database security
func TestDatabaseSecurity(t *testing.T) {
    // Test SQL injection prevention
    // Test ORM security features
    // Test database connection pooling
    // Test database transaction security
    // Test input sanitization for database queries
    // Test prepared statement usage
    // Test database error handling
}

// Test cases for transaction management
func TestTransactionSecurity(t *testing.T) {
    // Test atomic transaction operations
    // Test transaction rollback on failure
    // Test concurrent transaction handling
    // Test transaction isolation levels
    // Test transaction timeout handling
}
```

#### 3.5.2 Session Management
```go
// Test cases for session security
func TestSessionSecurity(t *testing.T) {
    // Test session token generation
    // Test session token expiration
    // Test session token validation
    // Test session fixation prevention
    // Test session hijacking prevention
    // Test session timeout handling
    // Test session invalidation
    // Test secure session cookies
    // Test token-based session validation
}

// Test cases for session data security
func TestSessionDataSecurity(t *testing.T) {
    // Test session data encryption
    // Test session data integrity
    // Test session data size limits
    // Test session data validation
    // Test Kratos session integration
}
```

## 4. Testing Implementation

### 4.1 Test Structure
```
internal/
├── middleware/
│   ├── kratos_test.go          # Kratos authentication tests
│   ├── validator_test.go       # Input validation tests
│   ├── rate_limit_test.go      # Rate limiting tests
│   └── cors_test.go           # CORS configuration tests
├── services/
│   ├── product_service_test.go # Product service tests
│   ├── checkout_service_test.go # Checkout service tests
│   └── media_service_test.go   # Media service tests
└── utils/
    └── response_test.go        # Error handling tests
```

### 4.2 Testing Framework
- **Testing Library**: Go's built-in testing package
- **Mocking**: testify/mock for service mocking
- **HTTP Testing**: httptest for API endpoint testing
- **Database Testing**: sqlmock for database mocking
- **Security Testing**: gosec for static analysis

### 4.3 Test Coverage Targets
- Authentication middleware: 100% coverage
- Input validation: 95%+ coverage
- Rate limiting: 90%+ coverage
- CORS: 100% coverage
- Error handling: 90%+ coverage
- File security: 85%+ coverage

## 5. Security Testing Tools

### 5.1 Static Analysis
- **gosec**: Go security checker
- **govulncheck**: Vulnerability scanner
- **staticcheck**: Static analysis tool

### 5.2 Dynamic Analysis
- **OWASP ZAP**: Web application security scanner
- **Burp Suite**: Web security testing
- **Postman**: API testing and security testing

### 5.3 Code Quality
- **golangci-lint**: Linter for Go
- **misspell**: Spelling checker
- **ineffassign**: Ineffective assignment checker

## 6. Test Data Management

### 6.1 Test Data Strategy
- **Realistic Test Data**: Use production-like data for testing
- **Sensitive Data Masking**: Mask sensitive information in test data
- **Data Isolation**: Separate test data from production data
- **Data Cleanup**: Automated test data cleanup

### 6.2 Test Data Categories
- **Valid Data**: Normal, expected input values
- **Invalid Data**: Malformed, out-of-range, or malicious input
- **Boundary Data**: Edge cases and boundary values
- **Performance Data**: Large datasets for load testing

## 7. Test Execution and Reporting

### 7.1 Test Execution
- **Automated Testing**: CI/CD pipeline integration
- **Scheduled Testing**: Regular security test execution
- **Ad-hoc Testing**: On-demand security testing
- **Regression Testing**: After code changes

### 7.2 Test Reporting
- **Test Results Dashboard**: Visual test results
- **Security Vulnerability Reports**: Detailed vulnerability findings
- **Code Coverage Reports**: Test coverage metrics
- **Remediation Tracking**: Issue tracking and resolution

## 8. Risk Assessment and Mitigation

### 8.1 High Priority Risks
- **Authentication Bypass**: Test token validation and session management
- **Authorization Bypass**: Test role-based access controls
- **SQL Injection**: Test input validation and ORM security
- **Cross-Site Scripting (XSS)**: Test input sanitization
- **Cross-Site Request Forgery (CSRF)**: Test CSRF protection

### 8.2 Medium Priority Risks
- **Rate Limiting Bypass**: Test rate limiting effectiveness
- **CORS Misconfiguration**: Test CORS security headers
- **Information Disclosure**: Test error message sanitization
- **File Upload Vulnerabilities**: Test file security measures
- **Session Fixation**: Test session management security

### 8.3 Low Priority Risks
- **Denial of Service**: Test rate limiting and resource limits
- **Insecure Direct Object References**: Test access controls
- **Security Header Issues**: Test security header configuration
- **Logging Insufficiencies**: Test logging and monitoring

## 9. Test Maintenance and Evolution

### 9.1 Test Maintenance
- **Regular Updates**: Update tests with code changes
- **New Feature Testing**: Add tests for new features
- **Vulnerability Updates**: Update tests for new vulnerabilities
- **Framework Updates**: Update testing tools and libraries

### 9.2 Test Evolution
- **Enhanced Coverage**: Increase test coverage over time
- **New Testing Types**: Add new testing methodologies
- **Performance Testing**: Expand performance and load testing
- **Compliance Testing**: Add compliance-specific tests

## 10. Conclusion

This security testing plan provides a comprehensive framework for ensuring the security of the Karima Store application. By implementing these tests, we can identify and mitigate security vulnerabilities, ensure proper authentication and authorization controls, validate input validation effectiveness, and maintain a secure application environment. The plan covers all critical security aspects and provides a structured approach to security testing that can be continuously improved and expanded as the application evolves.

**Note**: This plan has been updated to reflect the transition to Full Ory Kratos authentication, removing JWT-related tests and focusing exclusively on Kratos-based authentication and session management.