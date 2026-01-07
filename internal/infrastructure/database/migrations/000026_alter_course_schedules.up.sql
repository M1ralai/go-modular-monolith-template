ALTER TABLE course_schedules
    ADD COLUMN IF NOT EXISTS semester_start_date DATE,
    ADD COLUMN IF NOT EXISTS semester_end_date DATE,
    ADD COLUMN IF NOT EXISTS exclude_dates DATE[];
