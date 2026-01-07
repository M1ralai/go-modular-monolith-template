-- Fix goals table schema mismatches

ALTER TABLE goals RENAME COLUMN is_achieved TO is_completed;
ALTER TABLE goals RENAME COLUMN achieved_at TO completed_at;
ALTER TABLE goals ADD COLUMN IF NOT EXISTS priority VARCHAR(20) DEFAULT 'medium';
ALTER TABLE goals ADD COLUMN IF NOT EXISTS progress_percentage FLOAT8 DEFAULT 0;
