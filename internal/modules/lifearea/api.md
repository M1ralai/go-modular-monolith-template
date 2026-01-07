# LifeArea API

Base URL: `/api/life-areas`

## Endpoints

### POST /life-areas
Create a new life area
- Auth: Required
- Body: `CreateLifeAreaRequest`

### GET /life-areas
Get all life areas for current user
- Auth: Required

### GET /life-areas/{id}
Get life area by ID
- Auth: Required

### PUT /life-areas/{id}
Update life area
- Auth: Required

### DELETE /life-areas/{id}
Delete life area
- Auth: Required

For complete API documentation, see `/api/openapi.yaml`
