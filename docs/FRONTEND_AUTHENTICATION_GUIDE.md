# Frontend Authentication Integration Guide

## 📋 Overview

This document provides a comprehensive guide for frontend developers to integrate with the Karima Store API authentication system.

---

## 🔐 Authentication Methods

The API supports **two authentication methods**:

### 1. Cookie-Based Authentication (Recommended for Web)
- JWT token is automatically stored in an HTTP-only cookie named `jwt`
- No manual token management required
- More secure (XSS protection)
- Works seamlessly with browser requests

### 2. Bearer Token Authentication (Recommended for Mobile/API)
- JWT token returned in response body
- Must be included in `Authorization` header: `Bearer <token>`
- Manual token management required
- Suitable for mobile apps and API clients

---

## 📡 Authentication Endpoints

### 1. Register (Email/Password)

**Endpoint:** `POST /api/v1/auth/register`

**Description:** Register a new user account with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "full_name": "John Doe"
}
```

**Response (Success - 200):**
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

**Response (Error - 400):**
```json
{
  "error": "email already registered"
}
```

**Cookies Set:**
- `jwt` (HTTPOnly, Secure in production, 24h expiration)

**Frontend Implementation Example:**

```javascript
// Using fetch API
async function register(email, password, fullName) {
  const response = await fetch('http://localhost:8080/api/v1/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Important for cookies
    body: JSON.stringify({
      email,
      password,
      full_name: fullName
    })
  });

  const data = await response.json();
  
  if (response.ok) {
    // User registered successfully
    // JWT cookie is automatically set by browser
    console.log('Registration successful:', data.user);
    return data;
  } else {
    // Handle error
    console.error('Registration failed:', data.error);
    throw new Error(data.error);
  }
}

// Using Axios
import axios from 'axios';

async function register(email, password, fullName) {
  try {
    const response = await axios.post('http://localhost:8080/api/v1/auth/register', {
      email,
      password,
      full_name: fullName
    }, {
      withCredentials: true // Important for cookies
    });
    
    console.log('Registration successful:', response.data.user);
    return response.data;
  } catch (error) {
    console.error('Registration failed:', error.response?.data?.error);
    throw error;
  }
}
```

---

### 2. Login (Email/Password)

**Endpoint:** `POST /api/v1/auth/login`

**Description:** Authenticate user with email and password.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```

**Response (Success - 200):**
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

**Response (Error - 401):**
```json
{
  "error": "invalid credentials"
}
```

**Cookies Set:**
- `jwt` (HTTPOnly, Secure in production, 24h expiration)

**Frontend Implementation Example:**

```javascript
// Using fetch API
async function login(email, password) {
  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Important for cookies
    body: JSON.stringify({ email, password })
  });

  const data = await response.json();
  
  if (response.ok) {
    // Login successful
    // JWT cookie is automatically set by browser
    console.log('Login successful:', data.user);
    return data;
  } else {
    // Handle error
    console.error('Login failed:', data.error);
    throw new Error(data.error);
  }
}

// Using Axios
import axios from 'axios';

async function login(email, password) {
  try {
    const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
      email,
      password
    }, {
      withCredentials: true // Important for cookies
    });
    
    console.log('Login successful:', response.data.user);
    return response.data;
  } catch (error) {
    console.error('Login failed:', error.response?.data?.error);
    throw error;
  }
}
```

---

### 3. Logout

**Endpoint:** `GET /api/v1/auth/logout`

**Description:** Logout current user by clearing JWT cookie.

**Request:** No body required

**Response (Success - 200):**
```json
{
  "message": "Logged out successfully"
}
```

**Frontend Implementation Example:**

```javascript
// Using fetch API
async function logout() {
  const response = await fetch('http://localhost:8080/api/v1/auth/logout', {
    method: 'GET',
    credentials: 'include' // Important for cookies
  });

  if (response.ok) {
    console.log('Logout successful');
    // Redirect to login page
    window.location.href = '/login';
  }
}

// Using Axios
import axios from 'axios';

async function logout() {
  try {
    await axios.get('http://localhost:8080/api/v1/auth/logout', {
      withCredentials: true
    });
    
    console.log('Logout successful');
    // Redirect to login page
    window.location.href = '/login';
  } catch (error) {
    console.error('Logout failed:', error);
  }
}
```

---

### 4. Get Current User (Protected)

**Endpoint:** `GET /api/v1/auth/me`

**Description:** Get current authenticated user information.

**Authentication Required:** Yes (JWT cookie or Bearer token)

**Request:** No body required

**Response (Success - 200):**
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

**Response (Error - 401):**
```json
{
  "error": "Unauthorized"
}
```

**Frontend Implementation Example:**

```javascript
// Using fetch API (Cookie-based)
async function getCurrentUser() {
  const response = await fetch('http://localhost:8080/api/v1/auth/me', {
    method: 'GET',
    credentials: 'include' // Important for cookies
  });

  const data = await response.json();
  
  if (response.ok) {
    return data.user;
  } else {
    // User not authenticated
    console.error('Not authenticated:', data.error);
    throw new Error(data.error);
  }
}

// Using fetch API (Bearer token)
async function getCurrentUser(token) {
  const response = await fetch('http://localhost:8080/api/v1/auth/me', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  const data = await response.json();
  
  if (response.ok) {
    return data.user;
  } else {
    throw new Error(data.error);
  }
}

// Using Axios (Cookie-based)
import axios from 'axios';

async function getCurrentUser() {
  try {
    const response = await axios.get('http://localhost:8080/api/v1/auth/me', {
      withCredentials: true
    });
    
    return response.data.user;
  } catch (error) {
    console.error('Not authenticated:', error.response?.data?.error);
    throw error;
  }
}
```

---

### 5. OAuth Login (Google)

**Endpoint:** `GET /api/v1/auth/:provider`

**Description:** Initiate OAuth login flow with specified provider (e.g., Google).

**Parameters:**
- `provider` (path parameter): OAuth provider name (e.g., `google`)

**Request:** No body required

**Behavior:** Redirects to OAuth provider's login page

**Frontend Implementation Example:**

```javascript
// Redirect to Google OAuth login
function loginWithGoogle() {
  // Redirect user to OAuth login endpoint
  window.location.href = 'http://localhost:8080/api/v1/auth/google';
}

// Or open in new window
function loginWithGooglePopup() {
  const popup = window.open(
    'http://localhost:8080/api/v1/auth/google',
    'Google Login',
    'width=500,height=600'
  );
  
  // Listen for popup close or message
  const checkPopup = setInterval(() => {
    if (popup.closed) {
      clearInterval(checkPopup);
      // Check if user is authenticated
      getCurrentUser().then(user => {
        console.log('Logged in with Google:', user);
      });
    }
  }, 1000);
}
```

---

### 6. OAuth Callback

**Endpoint:** `GET /api/v1/auth/:provider/callback`

**Description:** Handle OAuth callback from provider.

**Parameters:**
- `provider` (path parameter): OAuth provider name (e.g., `google`)
- `code` (query parameter): Authorization code from OAuth provider
- `state` (query parameter): State parameter for CSRF protection

**Behavior:** 
- Completes OAuth authentication
- Creates or updates user in database
- Generates JWT token
- Sets JWT cookie
- Returns user data and token

**Response (Success - 200):**
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

**Frontend Implementation Example:**

```javascript
// The OAuth callback is handled automatically by the backend
// After OAuth login, the user will be redirected to this endpoint
// The backend will set the JWT cookie and return user data

// For popup-based OAuth, you can listen for messages:
window.addEventListener('message', (event) => {
  if (event.origin !== 'http://localhost:8080') return;
  
  if (event.data.type === 'oauth_success') {
    console.log('OAuth login successful:', event.data.user);
    // Redirect to dashboard
    window.location.href = '/dashboard';
  }
});

// For redirect-based OAuth, the callback page can:
// 1. Receive the token from URL (if backend redirects with token)
// 2. Or make a request to /api/v1/auth/me to get user data
// 3. Store token (if using Bearer token auth)
// 4. Redirect to dashboard

// Example callback page:
async function handleOAuthCallback() {
  try {
    // Get current user (JWT cookie is already set)
    const user = await getCurrentUser();
    console.log('OAuth login successful:', user);
    
    // Redirect to dashboard
    window.location.href = '/dashboard';
  } catch (error) {
    console.error('OAuth callback failed:', error);
    // Redirect to login page with error
    window.location.href = '/login?error=oauth_failed';
  }
}

// Call this function when callback page loads
handleOAuthCallback();
```

---

## 🛡️ Making Authenticated Requests

### Cookie-Based Authentication (Web)

```javascript
// All authenticated requests automatically include the JWT cookie
async function getOrders() {
  const response = await fetch('http://localhost:8080/api/v1/orders', {
    method: 'GET',
    credentials: 'include' // Important for cookies
  });

  if (response.ok) {
    const data = await response.json();
    return data.orders;
  } else {
    throw new Error('Failed to fetch orders');
  }
}

// Using Axios
import axios from 'axios';

// Configure axios to include cookies
axios.defaults.withCredentials = true;

async function getOrders() {
  const response = await axios.get('http://localhost:8080/api/v1/orders');
  return response.data.orders;
}
```

### Bearer Token Authentication (Mobile/API)

```javascript
// Store token after login
let authToken = null;

async function login(email, password) {
  const response = await fetch('http://localhost:8080/api/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email, password })
  });

  const data = await response.json();
  
  if (response.ok) {
    // Store token for future requests
    authToken = data.token;
    localStorage.setItem('auth_token', data.token);
    return data;
  } else {
    throw new Error(data.error);
  }
}

// Make authenticated request with Bearer token
async function getOrders() {
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch('http://localhost:8080/api/v1/orders', {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });

  if (response.ok) {
    const data = await response.json();
    return data.orders;
  } else if (response.status === 401) {
    // Token expired, redirect to login
    localStorage.removeItem('auth_token');
    window.location.href = '/login';
    throw new Error('Unauthorized');
  } else {
    throw new Error('Failed to fetch orders');
  }
}

// Using Axios with interceptor
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080'
});

// Add request interceptor to include token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Add response interceptor to handle 401 errors
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

async function getOrders() {
  const response = await api.get('/api/v1/orders');
  return response.data.orders;
}
```

---

## 🔒 CSRF Protection

The API uses CSRF protection for state-changing requests (POST, PUT, DELETE, PATCH).

### CSRF Token Handling

**CSRF Token Location:** 
- Cookie: `csrf_token`
- Header: `X-CSRF-Token`

**Excluded Paths:** Auth endpoints, OAuth callback, and public endpoints are excluded from CSRF protection.

**Frontend Implementation:**

```javascript
// For state-changing requests, include CSRF token
async function updateProfile(profileData) {
  // Get CSRF token from cookie
  const csrfToken = document.cookie
    .split('; ')
    .find(row => row.startsWith('csrf_token='))
    ?.split('=')[1];

  const response = await fetch('http://localhost:8080/api/v1/users/me', {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'X-CSRF-Token': csrfToken // Include CSRF token
    },
    credentials: 'include',
    body: JSON.stringify(profileData)
  });

  return response.json();
}

// Using Axios with CSRF interceptor
import axios from 'axios';

axios.interceptors.request.use((config) => {
  // Get CSRF token from cookie
  const csrfToken = document.cookie
    .split('; ')
    .find(row => row.startsWith('csrf_token='))
    ?.split('=')[1];

  // Add CSRF header for state-changing requests
  if (['post', 'put', 'delete', 'patch'].includes(config.method?.toLowerCase())) {
    config.headers['X-CSRF-Token'] = csrfToken;
  }

  return config;
});
```

---

## 📱 Complete React Example

```jsx
import { useState, useEffect } from 'react';
import axios from 'axios';

// Configure axios
axios.defaults.withCredentials = true;

function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  // Check authentication status on mount
  useEffect(() => {
    checkAuth();
  }, []);

  const checkAuth = async () => {
    try {
      const response = await axios.get('http://localhost:8080/api/v1/auth/me');
      setUser(response.data.user);
    } catch (error) {
      setUser(null);
    } finally {
      setLoading(false);
    }
  };

  const register = async (email, password, fullName) => {
    const response = await axios.post('http://localhost:8080/api/v1/auth/register', {
      email,
      password,
      full_name: fullName
    });
    setUser(response.data.user);
    return response.data;
  };

  const login = async (email, password) => {
    const response = await axios.post('http://localhost:8080/api/v1/auth/login', {
      email,
      password
    });
    setUser(response.data.user);
    return response.data;
  };

  const loginWithGoogle = () => {
    window.location.href = 'http://localhost:8080/api/v1/auth/google';
  };

  const logout = async () => {
    await axios.get('http://localhost:8080/api/v1/auth/logout');
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, loading, register, login, loginWithGoogle, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

// Login Component
function LoginForm() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const { login, loginWithGoogle } = useAuth();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await login(email, password);
    } catch (err) {
      setError(err.response?.data?.error || 'Login failed');
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="Email"
        required
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
        required
      />
      <button type="submit">Login</button>
      <button type="button" onClick={loginWithGoogle}>
        Login with Google
      </button>
      {error && <p className="error">{error}</p>}
    </form>
  );
}

// Protected Route Component
function ProtectedRoute({ children }) {
  const { user, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (!user) return <Navigate to="/login" />;

  return children;
}

// Usage
function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          <Route path="/login" element={<LoginForm />} />
          <Route
            path="/dashboard"
            element={
              <ProtectedRoute>
                <Dashboard />
              </ProtectedRoute>
            }
          />
        </Routes>
      </Router>
    </AuthProvider>
  );
}
```

---

## 🧪 Testing Checklist

### Authentication Flow
- [ ] Register new user with email/password
- [ ] Login with email/password
- [ ] Verify JWT cookie is set
- [ ] Access protected endpoint with cookie
- [ ] Logout and verify cookie is cleared

### OAuth Flow
- [ ] Initiate Google OAuth login
- [ ] Complete OAuth authentication
- [ ] Verify user is created/found in database
- [ ] Verify JWT cookie is set
- [ ] Access protected endpoint

### Error Handling
- [ ] Test invalid email/password
- [ ] Test duplicate email registration
- [ ] Test unauthorized access to protected endpoints
- [ ] Test expired token

---

## 📚 Additional Resources

- [API Standards](./api_standards.md)
- [Goth Fiber Integration Status](./GOTH_FIBER_INTEGRATION_STATUS.md)
- [Project Architecture](./architecture_path.md)
- [Swagger Documentation](http://localhost:8080/swagger)
