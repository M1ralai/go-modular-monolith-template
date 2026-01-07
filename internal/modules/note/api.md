# Note API

Base URL: `/api/notes`

## Endpoints

### POST /notes
Create a new note
- Auth: Required
- Can optionally link to course, component, or life area

### GET /notes
Get all notes for current user
- Auth: Required

### GET /notes/favorites
Get favorite notes
- Auth: Required

### GET /notes/search?q={query}
Search notes by title and content
- Auth: Required

### GET /notes/{id}
Get note by ID with backlinks
- Auth: Required
- Returns: outgoing_links, backlinks, backlink_count

### PUT /notes/{id}
Update note
- Auth: Required

### DELETE /notes/{id}
Delete note (also deletes associated links)
- Auth: Required

## Backlinks (Obsidian-like)

### GET /notes/{id}/backlinks
Get all notes linking to this note
- Auth: Required

### POST /notes/{id}/links
Create a link from this note to another
- Auth: Required
- Body: `{ "target_note_id": 123, "link_text": "optional" }`

### DELETE /notes/links/{linkId}
Delete a note link
- Auth: Required

**Note about backlinks:**
- When you create a link from Note A to Note B, Note B's `backlinks` will include Note A
- This enables Obsidian-like bidirectional linking
- When a note is deleted, all its links (both outgoing and incoming references) are cleaned up

For complete API documentation, see `/api/openapi.yaml`
