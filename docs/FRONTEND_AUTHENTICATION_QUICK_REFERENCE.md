# Frontend Authentication Quick Reference

## 🚀 Authentication Endpoints for Frontend

### 1. Register (Email/Password)
```
POST /api/v1/auth/register
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "full_name": "John Doe"
}
```
**Response:** User object + JWT token
**Cookie:** `jwt` (24h expiration)

---

### 2. Login (Email/Password)
```
POST /api/v1/auth/login
```
**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```
**Response:** User object + JWT token
**Cookie:** `jwt` (24h expiration)

---

### 3. Logout
```
GET /api/v1/auth/logout
```
**Response:** Success message
**Cookie:** `jwt` cleared

---

### 4. Get Current User (Protected)
```
GET /api/v1/auth/me
```
**Authentication Required:** Yes
**Response:** Current user object

---

### 5. OAuth Login (Google)
```
GET /api/v1/auth/google
```
**Behavior:** Redirects to Google OAuth login page

---

### 6. OAuth Callback
```
GET /api/v1/auth/google/callback
```
**Behavior:** Handles OAuth callback, creates/updates user, sets JWT cookie
**Response:** User object + JWT token

---

## 📋 Summary Table

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/auth/register` | POST | No | Register new user |
| `/api/v1/auth/login` | POST | No | Login with email/password |
| `/api/v1/auth/logout` | GET | No | Logout current user |
| `/api/v1/auth/me` | GET | Yes | Get current user |
| `/api/v1/auth/google` | GET | No | Initiate Google OAuth |
| `/api/v1/auth/google/callback` | GET | No | Handle Google OAuth callback |

---

## 🔐 Authentication Methods

### Cookie-Based (Web)
```javascript
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  credentials: 'include', // Important!
  body: JSON.stringify({ email, password })
});
```

### Bearer Token (Mobile/API)
```javascript
fetch('http://localhost:8080/api/v1/auth/me', {
  method: 'GET',
  headers: {
    'Authorization': `Bearer ${token}`
  }
});
```

---

## 📦 Response Structure

### Success Response
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
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Error Response
```json
{
  "error": "Error message here"
}
```

---

## 🍪 Cookies

### JWT Cookie
- **Name:** `jwt`
- **HTTPOnly:** Yes
- **Secure:** Yes (production only)
- **Expiration:** 24 hours

### CSRF Cookie
- **Name:** `csrf_token`
- **HTTPOnly:** No
- **Secure:** Yes (production only)
- **Expiration:** 24 hours

---

## 🔒 CSRF Protection

For state-changing requests (POST, PUT, DELETE, PATCH), include CSRF token:

```javascript
const csrfToken = document.cookie
  .split('; ')
  .find(row => row.startsWith('csrf_token='))
  ?.split('=')[1];

fetch('http://localhost:8080/api/v1/users/me', {
  method: 'PUT',
  headers: {
    'Content-Type': 'application/json',
    'X-CSRF-Token': csrfToken
  },
  credentials: 'include',
  body: JSON.stringify(data)
});
```

---

## 📱 Quick Implementation

### React Example
```jsx
import axios from 'axios';

axios.defaults.withCredentials = true;

// Login
const login = async (email, password) => {
  const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
    email,
    password
  });
  return response.data;
};

// Get Current User
const getCurrentUser = async () => {
  const response = await axios.get('http://localhost:8080/api/v1/auth/me');
  return response.data.user;
};

// Logout
const logout = async () => {
  await axios.get('http://localhost:8080/api/v1/auth/logout');
};

// Google OAuth
const loginWithGoogle = () => {
  window.location.href = 'http://localhost:8080/api/v1/auth/google';
};
```

### Vue Example
```javascript
import axios from 'axios';

axios.defaults.withCredentials = true;

export default {
  methods: {
    async login(email, password) {
      const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
        email,
        password
      });
      return response.data;
    },
    async getCurrentUser() {
      const response = await axios.get('http://localhost:8080/api/v1/auth/me');
      return response.data.user;
    },
    async logout() {
      await axios.get('http://localhost:8080/api/v1/auth/logout');
    },
    loginWithGoogle() {
      window.location.href = 'http://localhost:8080/api/v1/auth/google';
    }
  }
};
```

---

## 🧪 Testing

### Register
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"password123","full_name":"Test User"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -c cookies.txt \
  -d '{"email":"test@example.com","password":"password123"}'
```

### Get Current User
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -b cookies.txt
```

### Logout
```bash
curl -X GET http://localhost:8080/api/v1/auth/logout \
  -b cookies.txt
```

---

## 📚 Full Documentation

- **[Frontend Authentication Guide](./FRONTEND_AUTHENTICATION_GUIDE.md)** - Detailed implementation guide with React examples
- **[Authentication API Reference](./AUTHENTICATION_API_REFERENCE.md)** - Complete API reference
- **[Goth Fiber Integration Status](./GOTH_FIBER_INTEGRATION_STATUS.md)** - OAuth integration details
- **[API Standards](./api_standards.md)** - General API conventions
- **[Swagger Documentation](http://localhost:8080/swagger)** - Interactive API documentation

---

## ✅ Checklist for Frontend Integration

- [ ] Implement register endpoint
- [ ] Implement login endpoint
- [ ] Implement logout endpoint
- [ ] Implement get current user endpoint
- [ ] Implement Google OAuth login
- [ ] Handle JWT cookie (automatic with `credentials: 'include'`)
- [ ] Implement CSRF token handling for state-changing requests
- [ ] Handle authentication errors (401, 403, etc.)
- [ ] Implement protected routes
- [ ] Test all authentication flows

---

## 🎯 Key Points

1. **Cookie-Based Auth** (Recommended for Web): Set `credentials: 'include'` for all requests
2. **Bearer Token Auth** (Recommended for Mobile): Include `Authorization: Bearer <token>` header
3. **CSRF Protection**: Include `X-CSRF-Token` header for POST/PUT/DELETE/PATCH requests
4. **OAuth**: Redirect to `/api/v1/auth/google` to initiate Google login
5. **JWT Expiration**: Tokens expire after 24 hours
6. **Auto-Login**: JWT cookie is automatically sent with requests when using `credentials: 'include'`
