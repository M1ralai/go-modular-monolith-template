package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) TaskRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, task *domain.Task) (*domain.Task, error) {
	query := `
		INSERT INTO tasks (user_id, parent_task_id, title, description, due_date,
						   estimated_start, estimated_end, priority, is_completed,
						   progress_percentage, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	model := FromDomain(task)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.UserID,
		model.ParentTaskID,
		model.Title,
		model.Description,
		model.DueDate,
		model.EstimatedStart,
		model.EstimatedEnd,
		model.Priority,
		model.IsCompleted,
		model.ProgressPercentage,
		now,
		now,
	).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Task, error) {
	query := `
		SELECT id, user_id, parent_task_id, title, description, due_date,
			   estimated_start, estimated_end, actual_start, actual_end,
			   priority, is_completed, completed_at, progress_percentage,
			   created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	var model TaskModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Task, error) {
	query := `
		SELECT id, user_id, parent_task_id, title, description, due_date,
			   estimated_start, estimated_end, actual_start, actual_end,
			   priority, is_completed, completed_at, progress_percentage,
			   created_at, updated_at
		FROM tasks
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	var models []TaskModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(models))
	for i, m := range models {
		tasks[i] = m.ToDomain()
	}

	return tasks, nil
}

func (r *postgresRepository) GetSubtasks(ctx context.Context, parentID int) ([]*domain.Task, error) {
	query := `
		SELECT id, user_id, parent_task_id, title, description, due_date,
			   estimated_start, estimated_end, actual_start, actual_end,
			   priority, is_completed, completed_at, progress_percentage,
			   created_at, updated_at
		FROM tasks
		WHERE parent_task_id = $1
		ORDER BY created_at ASC
	`

	var models []TaskModel
	err := r.db.SelectContext(ctx, &models, query, parentID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(models))
	for i, m := range models {
		tasks[i] = m.ToDomain()
	}

	return tasks, nil
}

func (r *postgresRepository) GetParentTasks(ctx context.Context, userID int) ([]*domain.Task, error) {
	query := `
		SELECT id, user_id, parent_task_id, title, description, due_date,
			   estimated_start, estimated_end, actual_start, actual_end,
			   priority, is_completed, completed_at, progress_percentage,
			   created_at, updated_at
		FROM tasks
		WHERE user_id = $1 AND parent_task_id IS NULL
		ORDER BY created_at DESC
	`

	var models []TaskModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]*domain.Task, len(models))
	for i, m := range models {
		tasks[i] = m.ToDomain()
	}

	return tasks, nil
}

func (r *postgresRepository) Update(ctx context.Context, task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, due_date = $3, estimated_start = $4,
			estimated_end = $5, actual_start = $6, actual_end = $7, priority = $8,
			is_completed = $9, completed_at = $10, progress_percentage = $11, updated_at = $12
		WHERE id = $13
	`

	model := FromDomain(task)
	_, err := r.db.ExecContext(
		ctx, query,
		model.Title,
		model.Description,
		model.DueDate,
		model.EstimatedStart,
		model.EstimatedEnd,
		model.ActualStart,
		model.ActualEnd,
		model.Priority,
		model.IsCompleted,
		model.CompletedAt,
		model.ProgressPercentage,
		time.Now(),
		model.ID,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) CountSubtasks(ctx context.Context, parentID int) (total int, completed int, err error) {
	query := `
		SELECT
			COUNT(*) as total,
			COUNT(*) FILTER (WHERE is_completed = true) as completed
		FROM tasks
		WHERE parent_task_id = $1
	`

	err = r.db.QueryRowxContext(ctx, query, parentID).Scan(&total, &completed)
	return total, completed, err
}

func (r *postgresRepository) GetStats(ctx context.Context, userID int) (completedToday, dueToday, dueTomorrow, overdue int, err error) {
	// Completed today
	completedTodayQuery := `
		SELECT COUNT(*) 
		FROM tasks 
		WHERE user_id = $1 
		AND is_completed = true 
		AND completed_at IS NOT NULL
		AND DATE(completed_at) = CURRENT_DATE
	`
	err = r.db.QueryRowxContext(ctx, completedTodayQuery, userID).Scan(&completedToday)
	if err != nil {
		return
	}

	// Due today (not completed)
	dueTodayQuery := `
		SELECT COUNT(*) 
		FROM tasks 
		WHERE user_id = $1 
		AND is_completed = false 
		AND due_date IS NOT NULL
		AND DATE(due_date) = CURRENT_DATE
	`
	err = r.db.QueryRowxContext(ctx, dueTodayQuery, userID).Scan(&dueToday)
	if err != nil {
		return
	}

	// Due tomorrow (not completed)
	dueTomorrowQuery := `
		SELECT COUNT(*) 
		FROM tasks 
		WHERE user_id = $1 
		AND is_completed = false 
		AND due_date IS NOT NULL
		AND DATE(due_date) = CURRENT_DATE + INTERVAL '1 day'
	`
	err = r.db.QueryRowxContext(ctx, dueTomorrowQuery, userID).Scan(&dueTomorrow)
	if err != nil {
		return
	}

	// Overdue (not completed, due date before today)
	overdueQuery := `
		SELECT COUNT(*) 
		FROM tasks 
		WHERE user_id = $1 
		AND is_completed = false 
		AND due_date IS NOT NULL
		AND DATE(due_date) < CURRENT_DATE
	`
	err = r.db.QueryRowxContext(ctx, overdueQuery, userID).Scan(&overdue)
	return
}
