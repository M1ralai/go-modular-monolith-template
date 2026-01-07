-- Revert goals schema fixes

ALTER TABLE goals DROP COLUMN IF EXISTS progress_percentage;
ALTER TABLE goals DROP COLUMN IF EXISTS priority;
ALTER TABLE goals RENAME COLUMN completed_at TO achieved_at;
ALTER TABLE goals RENAME COLUMN is_completed TO is_achieved;
