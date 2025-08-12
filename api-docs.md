# Transit Backend API Documentation

Base URL: `http://localhost:3030/api/v1`

---

## Authentication Endpoints

### Register User
**POST** `/auth/register`

#### Request
```json
{
  "email": "user@example.com",
  "password": "yourpassword",
  "firstName": "John",
  "lastName": "Doe"
}
```

#### Response (Success)
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  },
  "accessToken": "JWT_TOKEN_HERE",
  "refreshToken": "REFRESH_TOKEN_HERE"
}
```

#### Response (Error)
```json
{
  "error": "Registration failed",
  "message": "Email already exists"
}
```

---

### Login User
**POST** `/auth/login`

#### Request
```json
{
  "email": "user@example.com",
  "password": "yourpassword"
}
```

#### Response (Success)
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  },
  "accessToken": "JWT_TOKEN_HERE",
  "refreshToken": "REFRESH_TOKEN_HERE"
}
```

#### Response (Error)
```json
{
  "error": "Login failed",
  "message": "Invalid credentials"
}
```

---

### Get User Profile
**GET** `/profile`
**Headers:** `Authorization: Bearer <accessToken>`

#### Response (Success)
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe"
  }
}
```

#### Response (Error)
```json
{
  "error": "Unauthorized",
  "message": "Invalid or missing token"
}
```

---

### Logout
**POST** `/logout`
**Headers:** `Authorization: Bearer <accessToken>`

#### Response
```json
{
  "message": "Logout successful"
}
```

---

## Waitlist Endpoints

### Join Waitlist
**POST** `/w aitlist`

#### Request
```json
{
  "email": "waitlistuser@example.com",
  "firstName": "Jane",
  "lastName": "Smith"
}
```

#### Response (Success)
```json
{
  "message": "Successfully joined the waitlist",
  "user": {
    "id": 2,
    "email": "waitlistuser@example.com",
    "firstName": "Jane",
    "lastName": "Smith"
  }
}
```

#### Response (Error)
```json
{
  "error": "Already on waitlist",
  "message": "This email is already registered"
}
```

---

### Get Waitlist Stats
**GET** `/waitlist/stats`

#### Response
```json
{
  "total": 123,
  "recent": [
    {
      "email": "waitlistuser@example.com",
      "firstName": "Jane",
      "lastName": "Smith",
      "joinedAt": "2025-08-12T10:00:00Z"
    }
    // ...more users
  ]
}
```

---

## Health Endpoints

### Basic Health Check
**GET** `/health`

#### Response
```json
{
  "status": "healthy",
  "timestamp": "2025-08-12T10:00:00Z",
  "database": "connected",
  "version": "1.0.0",
  "uptime": "1h23m"
}
```

---

### Detailed Health Check
**GET** `/health/detailed`

#### Response
```json
{
  "status": "healthy",
  "timestamp": "2025-08-12T10:00:00Z",
  "services": {
    "database": {
      "status": "healthy",
      "response_time": "10ms",
      "message": "Database connection successful"
    }
    // ...other services
  },
  "system": {
    "version": "1.0.0",
    "uptime": "1h23m"
  }
}
```

---

### Readiness Check
**GET** `/health/ready`

#### Response
```json
{
  "status": "ready"
}
```

---

### Liveness Check
**GET** `/health/live`

#### Response
```json
{
  "status": "alive",
  "timestamp": "2025-08-12T10:00:00Z"
}
```

---

## Usage Notes
- All endpoints expect and return JSON.
- Protected endpoints require a valid JWT in the `Authorization` header.
- For registration/login, use the returned `accessToken` for subsequent requests.
- Waitlist endpoints may be rate-limited.
