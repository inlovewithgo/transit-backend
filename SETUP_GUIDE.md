# Authentication System Setup Guide

## Quick Setup

1. **Copy environment variables:**
   ```bash
   copy .env.example .env
   ```

2. **Update the `.env` file with your values:**
   ```bash
   # JWT Configuration - Use a strong secret key
   JWT_SECRET=your_super_secret_jwt_key_here_make_it_long_and_random
   
   # Resend Email Configuration - Get from https://resend.com
   RESEND_API_KEY=re_your_resend_api_key_here
   ```

3. **Update email domain in mail service:**
   - Edit `main/service/mail-service.go`
   - Replace `noreply@yourdomain.com` with your verified Resend domain

4. **Start PostgreSQL database:**
   ```bash
   docker compose -f docker\pg\docker-compose.yaml up -d
   ```

5. **Run the application:**
   ```bash
   go run ./cmd
   ```

6. **Test the API:**
   ```bash
   # PowerShell
   .\test_auth_api.ps1
   
   # Bash (if available)
   ./test_auth_api.sh
   ```

## Resend Setup

1. Sign up at https://resend.com
2. Verify your domain or use the test domain
3. Get your API key from the dashboard
4. Add it to your `.env` file

## Database Schema

The application will automatically create the `users` table on startup with proper indexes and constraints.

## API Endpoints

- **POST** `/api/v1/auth/register` - Register new user
- **POST** `/api/v1/auth/login` - Login user
- **GET** `/api/v1/profile` - Get user profile (protected)
- **POST** `/api/v1/logout` - Logout user (protected)

## Frontend Integration

Use the `access_token` from login/register responses in the Authorization header:

```javascript
// Example fetch request
fetch('/api/v1/profile', {
  headers: {
    'Authorization': 'Bearer ' + accessToken,
    'Content-Type': 'application/json'
  }
})
```

## Email Templates

The system sends two types of emails:
1. **Welcome Email** - After successful registration
2. **Login Notification** - After successful login

Both emails are sent asynchronously and won't block the API response.
