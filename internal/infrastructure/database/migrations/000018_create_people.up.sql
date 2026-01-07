CREATE TABLE people (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  relationship VARCHAR(50),
  tags TEXT[],
  importance VARCHAR(20),
  last_contact_date DATE,
  contact_frequency VARCHAR(20),
  birthday DATE,
  phone VARCHAR(50),
  email VARCHAR(255),
  notes TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_people_user ON people(user_id);
CREATE INDEX idx_people_tags ON people USING gin(tags);
CREATE INDEX idx_people_relationship ON people(relationship);
CREATE INDEX idx_people_importance ON people(importance);
