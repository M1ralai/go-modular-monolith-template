CREATE TABLE habit_logs (
  id SERIAL PRIMARY KEY,
  habit_id INTEGER REFERENCES habits(id) ON DELETE CASCADE,
  completed_at TIMESTAMP DEFAULT NOW(),
  note TEXT,
  mood INTEGER CHECK (mood >= 1 AND mood <= 5),
  skipped BOOLEAN DEFAULT FALSE,
  skip_reason TEXT
);

CREATE UNIQUE INDEX idx_habit_logs_unique ON habit_logs(habit_id, DATE(completed_at));
CREATE INDEX idx_habit_logs_habit ON habit_logs(habit_id);
CREATE INDEX idx_habit_logs_date ON habit_logs(completed_at);
