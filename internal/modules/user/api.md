# User API

Base URL: `/api/users`

## Endpoints

### POST /users
Create a new user
- Auth: Not required (registration)
- Body: `CreateUserRequest`

### GET /users
Get all users
- Auth: Required

### GET /users/{id}
Get user by ID
- Auth: Required

### PUT /users/{id}
Update user
- Auth: Required

### DELETE /users/{id}
Delete user
- Auth: Required

For complete API documentation, see `/api/openapi.yaml`
