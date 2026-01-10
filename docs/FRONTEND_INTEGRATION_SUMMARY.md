# Frontend Authentication Integration Summary

## 📋 Overview

This document provides a complete summary of authentication requirements for the Karima Store frontend integration.

---

## 🔐 Authentication Methods Available

The Karima Store API supports **two authentication methods**:

### 1. Email/Password Authentication
- Register new users
- Login with email and password
- Logout functionality
- Session management via JWT cookies

### 2. OAuth Authentication (Google)
- Google OAuth 2.0 integration
- Automatic user creation/updates
- JWT token generation
- Session management via JWT cookies

---

## 📡 Required Authentication Endpoints

### Public Endpoints (No Authentication Required)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/auth/register` | POST | Register new user with email/password |
| `/api/v1/auth/login` | POST | Login with email/password |
| `/api/v1/auth/logout` | GET | Logout current user |
| `/api/v1/auth/google` | GET | Initiate Google OAuth login |
| `/api/v1/auth/google/callback` | GET | Handle Google OAuth callback |

### Protected Endpoints (Authentication Required)

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/api/v1/auth/me` | GET | Get current authenticated user |
| `/api/v1/users/me` | GET | Get current user profile |
| `/api/v1/orders` | GET | Get user's orders |
| `/api/v1/orders/:id` | GET | Get specific order |
| `/api/v1/checkout` | POST | Create checkout |

---

## 🍪 Cookie-Based Authentication (Recommended for Web)

### How It Works

1. **Login/Register**: Backend sets `jwt` cookie (HTTPOnly, 24h expiration)
2. **Authenticated Requests**: Browser automatically sends cookie with requests
3. **Logout**: Backend clears `jwt` cookie

### Frontend Implementation

```javascript
// Configure to include cookies
axios.defaults.withCredentials = true;

// Login
const login = async (email, password) => {
  const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
    email,
    password
  });
  // JWT cookie is automatically set by browser
  return response.data;
};

// Authenticated request
const getOrders = async () => {
  const response = await axios.get('http://localhost:8080/api/v1/orders');
  // JWT cookie is automatically sent
  return response.data.orders;
};
```

### Benefits

✅ Automatic token management
✅ More secure (HTTPOnly cookies prevent XSS)
✅ No manual token handling required
✅ Works seamlessly with browser requests

---

## 🎫 Bearer Token Authentication (Recommended for Mobile/API)

### How It Works

1. **Login/Register**: Backend returns JWT token in response body
2. **Storage**: Frontend stores token (localStorage, AsyncStorage, etc.)
3. **Authenticated Requests**: Frontend includes token in `Authorization` header

### Frontend Implementation

```javascript
// Login and store token
const login = async (email, password) => {
  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, password })
  });
  const data = await response.json();
  localStorage.setItem('auth_token', data.token);
  return data;
};

// Authenticated request
const getOrders = async () => {
  const token = localStorage.getItem('auth_token');
  const response = await fetch('http://localhost:8080/api/v1/orders', {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};
```

### Benefits

✅ Works with mobile apps
✅ No cookies required
✅ More control over token storage
✅ Suitable for API-only clients

---

## 🔒 CSRF Protection

### What Is CSRF?

Cross-Site Request Forgery (CSRF) protection prevents unauthorized requests from malicious websites.

### How It Works

1. Backend sets `csrf_token` cookie
2. Frontend reads cookie and includes it in `X-CSRF-Token` header
3. Backend validates token for state-changing requests (POST, PUT, DELETE, PATCH)

### Frontend Implementation

```javascript
// Get CSRF token from cookie
const getCsrfToken = () => {
  return document.cookie
    .split('; ')
    .find(row => row.startsWith('csrf_token='))
    ?.split('=')[1];
};

// Include CSRF token in state-changing requests
const updateProfile = async (data) => {
  const csrfToken = getCsrfToken();
  const response = await fetch('http://localhost:8080/api/v1/users/me', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken
    },
    credentials: 'include',
    body: JSON.stringify(data)
  });
  return response.json();
};
```

### Excluded Paths

The following endpoints are **excluded** from CSRF protection:
- `/api/v1/auth/*` (All auth endpoints)
- `/api/v1/auth/google/callback` (OAuth callback)
- `/api/v1/health` (Health check)
- `/swagger/*` (Documentation)
- All GET requests

---

## 📱 OAuth Integration (Google)

### OAuth Flow

1. **Initiate**: Frontend redirects user to `/api/v1/auth/google`
2. **Google Login**: User authenticates with Google
3. **Callback**: Google redirects to `/api/v1/auth/google/callback`
4. **Backend Processing**: Backend creates/updates user, generates JWT token
5. **Session**: Backend sets JWT cookie

### Frontend Implementation

```javascript
// Redirect to Google OAuth
const loginWithGoogle = () => {
  window.location.href = 'http://localhost:8080/api/v1/auth/google';
};

// Or open in popup
const loginWithGooglePopup = () => {
  const popup = window.open(
    'http://localhost:8080/api/v1/auth/google',
    'Google Login',
    'width=500,height=600'
  );

  // Monitor popup for completion
  const checkPopup = setInterval(() => {
    if (popup.closed) {
      clearInterval(checkPopup);
      // Check authentication status
      getCurrentUser();
    }
  }, 1000);
};
```

### User Data from OAuth

When user logs in via Google, the following data is available:
- Email (from Google account)
- Full name (from Google account)
- Avatar URL (from Google profile)
- Provider: `google`
- User is automatically marked as verified (`is_verified: true`)

---

## 📦 Response Structures

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
    "last_login_at": "2024-01-01T00:00:00Z",
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

### Common Error Codes

- `400 Bad Request`: Invalid request body or duplicate email
- `401 Unauthorized`: Invalid credentials or missing token
- `404 Not Found`: User not found
- `500 Internal Server Error`: Server error

---

## 🎯 Implementation Checklist

### Phase 1: Basic Authentication
- [ ] Implement register form
- [ ] Implement login form
- [ ] Implement logout functionality
- [ ] Store JWT token (if using Bearer token auth)
- [ ] Handle authentication errors

### Phase 2: Protected Routes
- [ ] Implement protected route wrapper
- [ ] Redirect unauthenticated users to login
- [ ] Check authentication status on app load
- [ ] Display current user info

### Phase 3: OAuth Integration
- [ ] Add "Login with Google" button
- [ ] Implement OAuth redirect flow
- [ ] Handle OAuth callback
- [ ] Update user state after OAuth login

### Phase 4: CSRF Protection
- [ ] Implement CSRF token retrieval
- [ ] Include CSRF token in state-changing requests
- [ ] Handle CSRF validation errors

### Phase 5: Error Handling
- [ ] Display user-friendly error messages
- [ ] Handle network errors
- [ ] Handle token expiration
- [ ] Implement retry logic for failed requests

---

## 🧪 Testing Checklist

### Authentication Flow
- [ ] Register new user
- [ ] Login with email/password
- [ ] Verify JWT cookie is set
- [ ] Access protected endpoint
- [ ] Logout and verify cookie is cleared

### OAuth Flow
- [ ] Initiate Google OAuth
- [ ] Complete Google authentication
- [ ] Verify user is created/found
- [ ] Verify JWT cookie is set
- [ ] Access protected endpoint

### Error Handling
- [ ] Test invalid email/password
- [ ] Test duplicate email registration
- [ ] Test unauthorized access
- [ ] Test expired token

### CSRF Protection
- [ ] Verify CSRF token is set
- [ ] Include CSRF token in requests
- [ ] Test CSRF validation

---

## 📚 Documentation Links

### Quick Start
- **[Frontend Authentication Quick Reference](./FRONTEND_AUTHENTICATION_QUICK_REFERENCE.md)** - Quick reference guide with code examples

### Detailed Guides
- **[Frontend Authentication Guide](./FRONTEND_AUTHENTICATION_GUIDE.md)** - Comprehensive implementation guide with React examples
- **[Authentication API Reference](./AUTHENTICATION_API_REFERENCE.md)** - Complete API reference with all endpoints

### Backend Integration
- **[Goth Fiber Integration Status](./GOTH_FIBER_INTEGRATION_STATUS.md)** - OAuth integration details and status
- **[API Standards](./api_standards.md)** - General API conventions and best practices

### Interactive Documentation
- **[Swagger Documentation](http://localhost:8080/swagger)** - Interactive API documentation

---

## 💡 Best Practices

### Security
1. **Never store tokens in localStorage** for cookie-based auth (use HTTPOnly cookies)
2. **Always use HTTPS** in production
3. **Validate user input** before sending to API
4. **Handle token expiration** gracefully
5. **Implement logout** on all pages

### User Experience
1. **Show loading states** during authentication
2. **Display clear error messages**
3. **Remember user session** (cookies handle this automatically)
4. **Provide multiple login options** (email/password + OAuth)
5. **Implement password strength validation**

### Performance
1. **Minimize API calls** (cache user data)
2. **Use optimistic updates** where appropriate
3. **Implement request debouncing** for search inputs
4. **Lazy load protected routes**

---

## 🆘 Troubleshooting

### Common Issues

**Issue:** JWT cookie not being set
- **Solution:** Ensure `credentials: 'include'` is set in fetch/axios config

**Issue:** 401 Unauthorized errors
- **Solution:** Check if JWT cookie is present and valid, or if Bearer token is included

**Issue:** CSRF validation errors
- **Solution:** Ensure `X-CSRF-Token` header is included in state-changing requests

**Issue:** OAuth callback fails
- **Solution:** Verify Google OAuth credentials and callback URL configuration

**Issue:** User not authenticated after OAuth
- **Solution:** Check if JWT cookie is set after OAuth callback

---

## 📞 Support

For questions or issues related to frontend authentication integration:

1. Check the documentation links above
2. Review the Swagger documentation at `http://localhost:8080/swagger`
3. Contact the backend team for API-specific issues

---

## ✅ Summary

The Karima Store API provides a complete authentication system with:

✅ Email/Password authentication
✅ Google OAuth integration
✅ Cookie-based authentication (web)
✅ Bearer token authentication (mobile/API)
✅ CSRF protection
✅ JWT token management
✅ User session management
✅ Protected routes
✅ Comprehensive documentation

Frontend developers can choose between cookie-based authentication (recommended for web) or Bearer token authentication (recommended for mobile/API), depending on their use case.

