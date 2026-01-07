CREATE TABLE contact_logs (
  id SERIAL PRIMARY KEY,
  person_id INTEGER REFERENCES people(id) ON DELETE CASCADE,
  contact_datetime TIMESTAMP DEFAULT NOW(),
  contact_type VARCHAR(50),
  notes TEXT,
  mood VARCHAR(20),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_contact_logs_person ON contact_logs(person_id);
CREATE INDEX idx_contact_logs_datetime ON contact_logs(contact_datetime);
CREATE INDEX idx_contact_logs_type ON contact_logs(contact_type);
