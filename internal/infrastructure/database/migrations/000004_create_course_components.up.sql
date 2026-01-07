CREATE TABLE course_components (
  id SERIAL PRIMARY KEY,
  course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
  type VARCHAR(50) NOT NULL,
  name VARCHAR(255) NOT NULL,
  weight DECIMAL(5,2),
  max_score DECIMAL(6,2),
  achieved_score DECIMAL(6,2),
  due_date TIMESTAMP,
  completion_date TIMESTAMP,
  is_completed BOOLEAN DEFAULT FALSE,
  notes TEXT,
  display_order INTEGER DEFAULT 0,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_course_components_course ON course_components(course_id);
CREATE INDEX idx_course_components_type ON course_components(type);
CREATE INDEX idx_course_components_due ON course_components(due_date);
