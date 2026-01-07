CREATE TABLE weekly_reviews (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  week_start_date DATE NOT NULL,
  week_end_date DATE NOT NULL,
  wins TEXT,
  challenges TEXT,
  lessons_learned TEXT,
  next_week_focus TEXT,
  overall_rating INTEGER CHECK (overall_rating >= 1 AND overall_rating <= 10),
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(user_id, week_start_date)
);

CREATE INDEX idx_weekly_reviews_user ON weekly_reviews(user_id);
CREATE INDEX idx_weekly_reviews_week ON weekly_reviews(week_start_date);
