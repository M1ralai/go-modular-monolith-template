ALTER TABLE events
    DROP COLUMN IF EXISTS is_recurring,
    DROP COLUMN IF EXISTS recurrence_rule,
    DROP COLUMN IF EXISTS google_event_id,
    DROP COLUMN IF EXISTS apple_event_id,
    DROP COLUMN IF EXISTS last_synced_at;
