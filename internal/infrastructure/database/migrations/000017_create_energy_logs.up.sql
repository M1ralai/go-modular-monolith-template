CREATE TABLE energy_logs (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  log_datetime TIMESTAMP DEFAULT NOW(),
  energy_level INTEGER CHECK (energy_level >= 1 AND energy_level <= 10),
  mood INTEGER CHECK (mood >= 1 AND mood <= 5),
  context TEXT,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_energy_logs_user ON energy_logs(user_id);
CREATE INDEX idx_energy_logs_datetime ON energy_logs(log_datetime);
