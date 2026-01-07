CREATE TABLE events (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  start_datetime TIMESTAMP NOT NULL,
  end_datetime TIMESTAMP NOT NULL,
  all_day BOOLEAN DEFAULT FALSE,
  location VARCHAR(255),
  category VARCHAR(50),
  course_id INTEGER REFERENCES courses(id) ON DELETE SET NULL,
  component_id INTEGER REFERENCES course_components(id) ON DELETE SET NULL,
  source VARCHAR(50) DEFAULT 'manual',
  recurrence_rule TEXT,
  reminder_time INTEGER,
  energy_level VARCHAR(20),
  is_completed BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_user ON events(user_id);
CREATE INDEX idx_events_start ON events(start_datetime);
CREATE INDEX idx_events_course ON events(course_id);
CREATE INDEX idx_events_component ON events(component_id);
CREATE INDEX idx_events_category ON events(category);
