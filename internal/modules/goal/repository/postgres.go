package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct{ db *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) GoalRepository { return &postgresRepository{db: db} }

func (r *postgresRepository) Create(ctx context.Context, goal *domain.Goal) (*domain.Goal, error) {
	query := `INSERT INTO goals (user_id, life_area_id, title, description, target_date, is_completed, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := FromDomain(goal)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.LifeAreaID, model.Title, model.Description, model.TargetDate, false, model.Priority, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Goal, error) {
	query := `SELECT id, user_id, life_area_id, title, description, target_date, is_completed, completed_at, priority, created_at, updated_at FROM goals WHERE id = $1`
	var model GoalModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Goal, error) {
	query := `SELECT id, user_id, life_area_id, title, description, target_date, is_completed, completed_at, priority, created_at, updated_at FROM goals WHERE user_id = $1 ORDER BY created_at DESC`
	var models []GoalModel
	if err := r.db.SelectContext(ctx, &models, query, userID); err != nil {
		return nil, err
	}
	goals := make([]*domain.Goal, len(models))
	for i, m := range models {
		goals[i] = m.ToDomain()
	}
	return goals, nil
}

func (r *postgresRepository) Update(ctx context.Context, goal *domain.Goal) error {
	query := `UPDATE goals SET title = $1, description = $2, target_date = $3, is_completed = $4, completed_at = $5, priority = $6, life_area_id = $7, updated_at = $8 WHERE id = $9`
	model := FromDomain(goal)
	_, err := r.db.ExecContext(ctx, query, model.Title, model.Description, model.TargetDate, model.IsCompleted, model.CompletedAt, model.Priority, model.LifeAreaID, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM goals WHERE id = $1`, id)
	return err
}

func (r *postgresRepository) CountMilestones(ctx context.Context, goalID int) (total int, completed int, err error) {
	query := `SELECT COUNT(*) as total, COUNT(*) FILTER (WHERE is_completed = true) as completed FROM milestones WHERE goal_id = $1`
	err = r.db.QueryRowxContext(ctx, query, goalID).Scan(&total, &completed)
	return
}
