-- Fix habit_logs table to match repository and add missing columns
ALTER TABLE habit_logs RENAME COLUMN completed_at TO log_date;
ALTER TABLE habit_logs ALTER COLUMN log_date TYPE DATE;
ALTER TABLE habit_logs RENAME COLUMN note TO notes;
ALTER TABLE habit_logs ADD COLUMN count INTEGER DEFAULT 0;
ALTER TABLE habit_logs ADD COLUMN is_completed BOOLEAN DEFAULT FALSE;
ALTER TABLE habit_logs ADD COLUMN created_at TIMESTAMP DEFAULT NOW();

-- Update unique index for the new column name
DROP INDEX IF EXISTS idx_habit_logs_unique;
CREATE UNIQUE INDEX idx_habit_logs_unique ON habit_logs(habit_id, log_date);
