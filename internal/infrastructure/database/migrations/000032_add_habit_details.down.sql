-- Remove missing columns from habits table
ALTER TABLE habits DROP COLUMN icon;
ALTER TABLE habits DROP COLUMN time_of_day;
ALTER TABLE habits DROP COLUMN reminder_time;
