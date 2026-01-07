CREATE TABLE blocked_time_slots (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    source_type VARCHAR(50),
    source_id INTEGER,
    start_datetime TIMESTAMP NOT NULL,
    end_datetime TIMESTAMP NOT NULL,
    reason VARCHAR(255),
    is_flexible BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_blocked_slots_user_time ON blocked_time_slots(user_id, start_datetime, end_datetime);
CREATE INDEX idx_blocked_slots_source ON blocked_time_slots(source_type, source_id);
