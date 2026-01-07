# Auth API

Base URL: `/auth`

## Endpoints

### POST /auth/login
User login
- Auth: Not required
- Body: `LoginRequest`
- Returns: JWT token + user info

### POST /auth/register
User registration
- Auth: Not required
- Body: `RegisterRequest`
- Returns: JWT token + user info

For complete API documentation, see `/api/openapi.yaml`
