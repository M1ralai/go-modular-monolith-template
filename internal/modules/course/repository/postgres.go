package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) CourseRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, course *domain.Course) (*domain.Course, error) {
	query := `
		INSERT INTO courses (user_id, name, code, instructor, credits, semester, type, color, syllabus_url, final_grade, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	model := FromDomain(course)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.UserID,
		model.Name,
		model.Code,
		model.Instructor,
		model.Credits,
		model.Semester,
		model.Type,
		model.Color,
		model.SyllabusURL,
		model.FinalGrade,
		model.IsActive,
		now,
		now,
	).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Course, error) {
	query := `
		SELECT id, user_id, name, code, instructor, credits, semester, type, color, syllabus_url, final_grade, is_active, created_at, updated_at
		FROM courses
		WHERE id = $1
	`

	var model CourseModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Course, error) {
	query := `
		SELECT id, user_id, name, code, instructor, credits, semester, type, color, syllabus_url, final_grade, is_active, created_at, updated_at
		FROM courses
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var models []CourseModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	courses := make([]*domain.Course, len(models))
	for i, m := range models {
		courses[i] = m.ToDomain()
	}

	return courses, nil
}

func (r *postgresRepository) GetActiveCourses(ctx context.Context, userID int) ([]*domain.Course, error) {
	query := `
		SELECT id, user_id, name, code, instructor, credits, semester, type, color, syllabus_url, final_grade, is_active, created_at, updated_at
		FROM courses
		WHERE user_id = $1 AND is_active = true
		ORDER BY name ASC
	`

	var models []CourseModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	courses := make([]*domain.Course, len(models))
	for i, m := range models {
		courses[i] = m.ToDomain()
	}

	return courses, nil
}

func (r *postgresRepository) Update(ctx context.Context, course *domain.Course) error {
	query := `
		UPDATE courses
		SET name = $1, code = $2, instructor = $3, credits = $4, semester = $5, type = $6, color = $7, syllabus_url = $8, final_grade = $9, is_active = $10, updated_at = $11
		WHERE id = $12
	`

	model := FromDomain(course)
	_, err := r.db.ExecContext(
		ctx, query,
		model.Name,
		model.Code,
		model.Instructor,
		model.Credits,
		model.Semester,
		model.Type,
		model.Color,
		model.SyllabusURL,
		model.FinalGrade,
		model.IsActive,
		time.Now(),
		model.ID,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM courses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) CreateComponent(ctx context.Context, comp *domain.Component) (*domain.Component, error) {
	query := `
		INSERT INTO course_components (course_id, type, name, weight, max_score, achieved_score, due_date, is_completed, notes, display_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	model := FromDomainComponent(comp)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.CourseID,
		model.Type,
		model.Name,
		model.Weight,
		model.MaxScore,
		model.AchievedScore,
		model.DueDate,
		model.IsCompleted,
		model.Notes,
		model.DisplayOrder,
		now,
		now,
	).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) UpdateComponent(ctx context.Context, comp *domain.Component) error {
	query := `
		UPDATE course_components
		SET type = $1, name = $2, weight = $3, max_score = $4, achieved_score = $5, due_date = $6, completion_date = $7, is_completed = $8, notes = $9, display_order = $10, updated_at = $11
		WHERE id = $12
	`

	model := FromDomainComponent(comp)
	_, err := r.db.ExecContext(
		ctx, query,
		model.Type,
		model.Name,
		model.Weight,
		model.MaxScore,
		model.AchievedScore,
		model.DueDate,
		model.CompletionDate,
		model.IsCompleted,
		model.Notes,
		model.DisplayOrder,
		time.Now(),
		model.ID,
	)

	return err
}

func (r *postgresRepository) DeleteComponent(ctx context.Context, id int) error {
	query := `DELETE FROM course_components WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) GetComponents(ctx context.Context, courseID int) ([]*domain.Component, error) {
	query := `
		SELECT id, course_id, type, name, weight, max_score, achieved_score, due_date, completion_date, is_completed, notes, display_order, created_at, updated_at
		FROM course_components
		WHERE course_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	var models []ComponentModel
	err := r.db.SelectContext(ctx, &models, query, courseID)
	if err != nil {
		return nil, err
	}

	components := make([]*domain.Component, len(models))
	for i, m := range models {
		components[i] = m.ToDomain()
	}

	return components, nil
}

func (r *postgresRepository) GetComponentByID(ctx context.Context, id int) (*domain.Component, error) {
	query := `
		SELECT id, course_id, type, name, weight, max_score, achieved_score, due_date, completion_date, is_completed, notes, display_order, created_at, updated_at
		FROM course_components
		WHERE id = $1
	`

	var model ComponentModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) CreateSchedule(ctx context.Context, sched *domain.Schedule) (*domain.Schedule, error) {
	query := `
		INSERT INTO course_schedules (course_id, day_of_week, start_time, end_time, location, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	now := time.Now()
	model := FromDomainSchedule(sched)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.CourseID,
		model.DayOfWeek,
		model.StartTime,
		model.EndTime,
		model.Location,
		now,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) UpdateSchedule(ctx context.Context, sched *domain.Schedule) error {
	query := `
		UPDATE course_schedules
		SET day_of_week = $1, start_time = $2, end_time = $3, location = $4
		WHERE id = $5
	`

	model := FromDomainSchedule(sched)
	_, err := r.db.ExecContext(
		ctx, query,
		model.DayOfWeek,
		model.StartTime,
		model.EndTime,
		model.Location,
		model.ID,
	)

	return err
}

func (r *postgresRepository) DeleteSchedule(ctx context.Context, id int) error {
	query := `DELETE FROM course_schedules WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) GetSchedules(ctx context.Context, courseID int) ([]*domain.Schedule, error) {
	query := `
		SELECT 
			id, 
			course_id, 
			day_of_week, 
			start_time::text as start_time, 
			end_time::text as end_time, 
			location, 
			created_at
		FROM course_schedules
		WHERE course_id = $1
		ORDER BY 
			CASE day_of_week
				WHEN 'Monday' THEN 1
				WHEN 'Tuesday' THEN 2
				WHEN 'Wednesday' THEN 3
				WHEN 'Thursday' THEN 4
				WHEN 'Friday' THEN 5
				WHEN 'Saturday' THEN 6
				WHEN 'Sunday' THEN 7
				ELSE 8
			END,
			start_time ASC
	`

	var models []ScheduleModel
	err := r.db.SelectContext(ctx, &models, query, courseID)
	if err != nil {
		return nil, err
	}

	schedules := make([]*domain.Schedule, len(models))
	for i, m := range models {
		schedules[i] = m.ToDomain()
	}

	return schedules, nil
}

func (r *postgresRepository) GetScheduleByID(ctx context.Context, id int) (*domain.Schedule, error) {
	query := `
		SELECT 
			id, 
			course_id, 
			day_of_week, 
			start_time::text as start_time, 
			end_time::text as end_time, 
			location, 
			created_at
		FROM course_schedules
		WHERE id = $1
	`

	var model ScheduleModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}
