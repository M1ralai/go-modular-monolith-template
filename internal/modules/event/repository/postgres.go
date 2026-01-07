package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct{ db *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) EventRepository { return &postgresRepository{db: db} }

func (r *postgresRepository) Create(ctx context.Context, event *domain.Event) (*domain.Event, error) {
	query := `INSERT INTO events (user_id, life_area_id, title, description, start_time, end_time, location, is_all_day, is_recurring, recurrence, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := FromDomain(event)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.LifeAreaID, model.Title, model.Description, model.StartTime, model.EndTime, model.Location, model.IsAllDay, model.IsRecurring, model.Recurrence, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Event, error) {
	query := `SELECT id, user_id, life_area_id, title, description, start_time, end_time, location, is_all_day, is_recurring, recurrence, created_at, updated_at FROM events WHERE id = $1`
	var model EventModel
	if err := r.db.GetContext(ctx, &model, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Event, error) {
	query := `SELECT id, user_id, life_area_id, title, description, start_time, end_time, location, is_all_day, is_recurring, recurrence, created_at, updated_at FROM events WHERE user_id = $1 ORDER BY start_time`
	var models []EventModel
	if err := r.db.SelectContext(ctx, &models, query, userID); err != nil {
		return nil, err
	}
	events := make([]*domain.Event, len(models))
	for i, m := range models {
		events[i] = m.ToDomain()
	}
	return events, nil
}

func (r *postgresRepository) GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.Event, error) {
	query := `SELECT id, user_id, life_area_id, title, description, start_time, end_time, location, is_all_day, is_recurring, recurrence, created_at, updated_at FROM events WHERE user_id = $1 AND start_time BETWEEN $2 AND $3 ORDER BY start_time`
	var models []EventModel
	if err := r.db.SelectContext(ctx, &models, query, userID, start, end); err != nil {
		return nil, err
	}
	events := make([]*domain.Event, len(models))
	for i, m := range models {
		events[i] = m.ToDomain()
	}
	return events, nil
}

func (r *postgresRepository) Update(ctx context.Context, event *domain.Event) error {
	query := `UPDATE events SET title = $1, description = $2, start_time = $3, end_time = $4, location = $5, is_all_day = $6, is_recurring = $7, recurrence = $8, life_area_id = $9, updated_at = $10 WHERE id = $11`
	model := FromDomain(event)
	_, err := r.db.ExecContext(ctx, query, model.Title, model.Description, model.StartTime, model.EndTime, model.Location, model.IsAllDay, model.IsRecurring, model.Recurrence, model.LifeAreaID, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	return err
}
