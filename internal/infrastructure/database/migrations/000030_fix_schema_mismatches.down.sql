-- Down migration to revert schema fixes

-- People Revert
ALTER TABLE people DROP COLUMN IF EXISTS company;

-- Habits Revert
ALTER TABLE habits RENAME COLUMN longest_streak TO best_streak;
ALTER TABLE habits RENAME COLUMN name TO title;

-- Tasks Revert
ALTER TABLE tasks DROP COLUMN IF EXISTS progress_percentage;
