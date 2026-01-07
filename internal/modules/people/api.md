# People API

Base URL: `/api/people`

## Endpoints

### POST /people - Create person
### GET /people - Get all people
### GET /people/search?q={query} - Search by name/email/company
### GET /people/tag/{tag} - Search by tag (PostgreSQL array)
### GET /people/{id} - Get person by ID
### PUT /people/{id} - Update person
### DELETE /people/{id} - Delete person

**Features:** Tag-based filtering (TEXT[] column), relationship tracking

For complete API documentation, see `/api/openapi.yaml`
