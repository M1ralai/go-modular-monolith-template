-- Add missing columns to habits table
ALTER TABLE habits ADD COLUMN icon VARCHAR(50);
ALTER TABLE habits ADD COLUMN time_of_day VARCHAR(20);
ALTER TABLE habits ADD COLUMN reminder_time VARCHAR(10);
