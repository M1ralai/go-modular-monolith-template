-- Add missing columns and fix mismatches between code models and database schema

-- Tasks Fixes
ALTER TABLE tasks ADD COLUMN IF NOT EXISTS progress_percentage FLOAT8 DEFAULT 0;

-- Habits Fixes
ALTER TABLE habits RENAME COLUMN title TO name;
ALTER TABLE habits RENAME COLUMN best_streak TO longest_streak;

-- People Fixes
ALTER TABLE people ADD COLUMN IF NOT EXISTS company VARCHAR(255);
