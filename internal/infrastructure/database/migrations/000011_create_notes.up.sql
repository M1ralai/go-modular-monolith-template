CREATE TABLE notes (
  id SERIAL PRIMARY KEY,
  user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
  parent_note_id INTEGER REFERENCES notes(id) ON DELETE SET NULL,
  title VARCHAR(255) NOT NULL,
  content TEXT,
  life_area_id INTEGER REFERENCES life_areas(id) ON DELETE SET NULL,
  course_id INTEGER REFERENCES courses(id) ON DELETE SET NULL,
  component_id INTEGER REFERENCES course_components(id) ON DELETE SET NULL,
  linked_event_id INTEGER REFERENCES events(id) ON DELETE SET NULL,
  linked_task_id INTEGER REFERENCES tasks(id) ON DELETE SET NULL,
  tags TEXT[],
  is_archived BOOLEAN DEFAULT FALSE,
  is_favorite BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_notes_user ON notes(user_id);
CREATE INDEX idx_notes_course ON notes(course_id);
CREATE INDEX idx_notes_component ON notes(component_id);
CREATE INDEX idx_notes_parent ON notes(parent_note_id);
CREATE INDEX idx_notes_tags ON notes USING gin(tags);
CREATE INDEX idx_notes_search ON notes USING gin(to_tsvector('english', title || ' ' || COALESCE(content, '')));
