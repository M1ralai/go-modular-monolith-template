# Schedule API

Base URL: `/api/schedule`

## Conflict Detection

### POST /schedule/check-conflict
Check if a time range conflicts with existing blocked slots

**Request:**
```json
{
  "start": "2025-02-10T09:30:00Z",
  "end": "2025-02-10T11:00:00Z"
}
```

**Response (Conflict):**
```json
{
  "has_conflict": true,
  "reason": "Matematik Dersi (A-204)",
  "suggestions": [
    {"start": "2025-02-10T11:00:00Z", "end": "2025-02-10T12:30:00Z", "duration_minutes": 90}
  ]
}
```

## Free Slots

### GET /schedule/free-slots?date=2025-02-10&duration=120
Find available time slots for a given date and duration

## Blocked Slots

### GET /schedule/blocked-slots?date=2025-02-10
List blocked time slots for a date

### POST /schedule/blocked-slots
Create a blocked time slot

### DELETE /schedule/blocked-slots/{id}
Delete a blocked time slot

## Event Generation

### POST /schedule/generate-events
Generate semester events from course schedule

**Request:**
```json
{
  "title": "Matematik",
  "day_of_week": "Monday",
  "start_time": "09:00",
  "end_time": "10:30",
  "location": "A-204",
  "semester_start_date": "2025-02-10",
  "semester_end_date": "2025-06-15",
  "exclude_dates": ["2025-04-21", "2025-04-22"]
}
```

**Response:**
```json
{
  "events_generated": 14,
  "blocked_slots_created": 14
}
```

**Features:**
- Automatic weekly event generation
- Holiday/break exclusion
- Conflict-free blocked slot creation
- Time overlap detection

For complete API documentation, see `/api/openapi.yaml`
