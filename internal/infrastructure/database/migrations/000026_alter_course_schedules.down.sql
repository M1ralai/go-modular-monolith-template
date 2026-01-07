ALTER TABLE course_schedules
    DROP COLUMN IF EXISTS semester_start_date,
    DROP COLUMN IF EXISTS semester_end_date,
    DROP COLUMN IF EXISTS exclude_dates;
