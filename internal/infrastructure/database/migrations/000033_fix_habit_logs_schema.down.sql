-- Revert habit_logs schema changes
ALTER TABLE habit_logs RENAME COLUMN log_date TO completed_at;
ALTER TABLE habit_logs ALTER COLUMN completed_at TYPE TIMESTAMP;
ALTER TABLE habit_logs RENAME COLUMN notes TO note;
ALTER TABLE habit_logs DROP COLUMN count;
ALTER TABLE habit_logs DROP COLUMN is_completed;
ALTER TABLE habit_logs DROP COLUMN created_at;

DROP INDEX IF EXISTS idx_habit_logs_unique;
CREATE UNIQUE INDEX idx_habit_logs_unique ON habit_logs(habit_id, DATE(completed_at));
