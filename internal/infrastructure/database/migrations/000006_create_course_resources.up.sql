CREATE TABLE course_resources (
  id SERIAL PRIMARY KEY,
  course_id INTEGER REFERENCES courses(id) ON DELETE CASCADE,
  component_id INTEGER REFERENCES course_components(id) ON DELETE SET NULL,
  title VARCHAR(255) NOT NULL,
  type VARCHAR(50),
  url VARCHAR(500),
  file_path VARCHAR(500),
  description TEXT,
  tags TEXT[],
  is_primary BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_course_resources_course ON course_resources(course_id);
CREATE INDEX idx_course_resources_component ON course_resources(component_id);
CREATE INDEX idx_course_resources_type ON course_resources(type);
CREATE INDEX idx_course_resources_tags ON course_resources USING gin(tags);
