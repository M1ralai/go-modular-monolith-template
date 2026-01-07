# Calendar API

Base URL: `/api/calendar`

## Google Calendar Integration

### POST /calendar/google/connect
Get Google OAuth authorization URL

**Response:**
```json
{
  "auth_url": "https://accounts.google.com/o/oauth2/..."
}
```

### GET /calendar/google/callback?code=xxx
Complete OAuth flow with authorization code

### POST /calendar/google/disconnect
Disconnect Google Calendar

### POST /calendar/google/sync
Trigger manual sync with Google Calendar

### GET /calendar/status
Get sync status for all integrations

### GET /calendar/integrations
List all calendar integrations

**Features:**
- Google OAuth 2.0 flow
- Automatic token refresh
- Two-way sync (local ↔ Google)
- Sync queue with retry logic

**Environment Variables:**
- GOOGLE_CLIENT_ID
- GOOGLE_CLIENT_SECRET
- GOOGLE_REDIRECT_URI

For complete API documentation, see `/api/openapi.yaml`
