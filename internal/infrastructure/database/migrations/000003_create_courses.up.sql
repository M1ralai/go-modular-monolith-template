CREATE TABLE courses (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  code VARCHAR(50),
  instructor VARCHAR(255),
  credits DECIMAL(3,1),
  semester VARCHAR(50),
  type VARCHAR(50),
  color VARCHAR(7),
  syllabus_url VARCHAR(500),
  final_grade VARCHAR(10),
  is_active BOOLEAN DEFAULT TRUE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_courses_user ON courses(user_id);
CREATE INDEX idx_courses_semester ON courses(semester);
CREATE INDEX idx_courses_active ON courses(is_active);
