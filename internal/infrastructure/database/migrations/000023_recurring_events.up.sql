CREATE TABLE recurring_events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_schedule_id INTEGER REFERENCES course_schedules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    day_of_week VARCHAR(10) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    location VARCHAR(255),
    start_date DATE NOT NULL,
    end_date DATE,
    exclude_dates DATE[],
    google_event_id VARCHAR(255),
    apple_event_id VARCHAR(255),
    last_synced_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_recurring_events_user ON recurring_events(user_id);
CREATE INDEX idx_recurring_events_schedule ON recurring_events(course_schedule_id);
