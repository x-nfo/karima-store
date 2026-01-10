# Authentication API Reference for Frontend

## 📡 Authentication Endpoints Summary

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login with email/password | No |
| GET | `/api/v1/auth/logout` | Logout current user | No |
| GET | `/api/v1/auth/me` | Get current user | Yes |
| GET | `/api/v1/auth/:provider` | Initiate OAuth login | No |
| GET | `/api/v1/auth/:provider/callback` | OAuth callback handler | No |

---

## 🔐 1. Register (Email/Password)

**Endpoint:** `POST /api/v1/auth/register`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "full_name": "John Doe"
}
```

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "customer",
    "is_active": true,
    "is_verified": false,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cookies Set:**
- `jwt` (HTTPOnly, Secure in production, 24h expiration)

**Frontend Implementation:**
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ email, password, full_name })
});
```

---

## 🔐 2. Login (Email/Password)

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "role": "customer",
    "is_active": true,
    "is_verified": false,
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Cookies Set:**
- `jwt` (HTTPOnly, Secure in production, 24h expiration)

**Frontend Implementation:**
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include',
  body: JSON.stringify({ email, password })
});
```

---

## 🔐 3. Logout

**Endpoint:** `GET /api/v1/auth/logout`

**Request:** No body required

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

**Frontend Implementation:**
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/logout', {
  method: 'GET',
  credentials: 'include'
});
```

---

## 🔐 4. Get Current User (Protected)

**Endpoint:** `GET /api/v1/auth/me`

**Authentication Required:** Yes (JWT cookie or Bearer token)

**Request:** No body required

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "full_name": "John Doe",
    "avatar": "https://example.com/avatar.jpg",
    "role": "customer",
    "is_active": true,
    "is_verified": false,
    "last_login_at": "2024-01-01T00:00:00Z",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**Frontend Implementation (Cookie-based):**
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/me', {
  method: 'GET',
  credentials: 'include'
});
```

**Frontend Implementation (Bearer token):**
```javascript
const response = await fetch('http://localhost:8080/api/v1/auth/me', {
  method: 'GET',
  headers: { 'Authorization': `Bearer ${token}` }
});
```

---

## 🔐 5. OAuth Login (Google)

**Endpoint:** `GET /api/v1/auth/:provider`

**Parameters:**
- `provider` (path): OAuth provider name (e.g., `google`)

**Request:** No body required

**Behavior:** Redirects to OAuth provider's login page

**Frontend Implementation:**
```javascript
// Redirect to Google OAuth
window.location.href = 'http://localhost:8080/api/v1/auth/google';

// Or open in popup
const popup = window.open('http://localhost:8080/api/v1/auth/google', 'Google Login', 'width=500,height=600');
```

---

## 🔐 6. OAuth Callback

**Endpoint:** `GET /api/v1/auth/:provider/callback`

**Parameters:**
- `provider` (path): OAuth provider name (e.g., `google`)
- `code` (query): Authorization code from OAuth provider
- `state` (query): State parameter for CSRF protection

**Behavior:**
- Completes OAuth authentication
- Creates/updates user in database
- Generates JWT token
- Sets JWT cookie
- Returns user data and token

**Response (200):**
```json
{
  "user": {
    "id": 1,
    "email": "user@gmail.com",
    "full_name": "John Doe",
    "avatar": "https://lh3.googleusercontent.com/...",
    "role": "customer",
    "is_active": true,
    "is_verified": true,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Frontend Implementation:**
```javascript
// The callback is handled automatically by the backend
// After OAuth login, the user will be redirected to this endpoint
// The backend will set the JWT cookie and return user data

// Check if user is authenticated after OAuth callback
async function checkAuthAfterOAuth() {
  const response = await fetch('http://localhost:8080/api/v1/auth/me', {
    method: 'GET',
    credentials: 'include'
  });

  if (response.ok) {
    const data = await response.json();
    return data.user;
  }
  return null;
}
```

---

## 🛡️ Authentication Methods

### Method 1: Cookie-Based (Recommended for Web)

**Pros:**
- Automatic token management
- More secure (HTTPOnly cookies)
- No manual token handling required

**Implementation:**
```javascript
// Set credentials to 'include' for all requests
axios.defaults.withCredentials = true;

// Or use fetch with credentials
fetch('http://localhost:8080/api/v1/orders', {
  credentials: 'include'
});
```

### Method 2: Bearer Token (Recommended for Mobile/API)

**Pros:**
- Works with mobile apps
- No cookies required
- More control over token storage

**Implementation:**
```javascript
// Store token after login
const token = localStorage.getItem('auth_token');

// Include token in Authorization header
fetch('http://localhost:8080/api/v1/orders', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

---

## 🔒 CSRF Protection

**CSRF Token Location:**
- Cookie: `csrf_token`
- Header: `X-CSRF-Token`

**Excluded Paths:** Auth endpoints, OAuth callback, and public endpoints

**Frontend Implementation:**
```javascript
// Get CSRF token from cookie
const csrfToken = document.cookie
  .split('; ')
  .find(row => row.startsWith('csrf_token='))
  ?.split('=')[1];

// Include CSRF header for state-changing requests
fetch('http://localhost:8080/api/v1/users/me', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': csrfToken
  },
  credentials: 'include',
  body: JSON.stringify(profileData)
});
```

---

## 📋 User Object Structure

```typescript
interface User {
  id: number;
  email: string;
  full_name: string;
  avatar?: string;
  role: 'customer' | 'admin';
  is_active: boolean;
  is_verified: boolean;
  last_login_at?: string;
  created_at: string;
  updated_at: string;
}
```

---

## ❌ Error Responses

### 400 Bad Request
```json
{
  "error": "email already registered"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized"
}
```

### 404 Not Found
```json
{
  "error": "User not found"
}
```

---

## 🧪 Quick Test Commands

```bash
# Register new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"password123","full_name":"Test User"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"password123"}'

# Get current user
curl -X GET http://localhost:8080/api/v1/auth/me \
  -b cookies.txt

# Logout
curl -X GET http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

---

## 📚 Related Documentation

- [Frontend Authentication Guide](./FRONTEND_AUTHENTICATION_GUIDE.md) - Detailed implementation guide
- [Goth Fiber Integration Status](./GOTH_FIBER_INTEGRATION_STATUS.md) - OAuth integration details
- [API Standards](./api_standards.md) - General API conventions
- [Swagger Documentation](http://localhost:8080/swagger) - Interactive API documentation

