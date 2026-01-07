# Habit API

Base URL: `/api/habits`

## Endpoints

### POST /habits
Create a new habit

### GET /habits
Get all habits (with completed_today status)

### GET /habits/active
Get only active habits

### GET /habits/{id}
Get habit by ID

### PUT /habits/{id}
Update habit

### DELETE /habits/{id}
Delete habit

### POST /habits/{id}/log
Log habit for today
- Body: `{ "count": 1, "notes": "optional" }`
- **Business Logic**: If count >= target_count, increments streak

**Streak Tracking:**
- `current_streak`: Consecutive days completed
- `longest_streak`: Best streak ever
- `completed_today`: Whether habit was logged today

For complete API documentation, see `/api/openapi.yaml`
