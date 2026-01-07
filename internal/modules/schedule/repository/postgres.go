package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/domain"
	"github.com/jmoiron/sqlx"
)

type postgresBlockedSlotRepo struct{ db *sqlx.DB }

func NewBlockedTimeSlotRepository(db *sqlx.DB) BlockedTimeSlotRepository {
	return &postgresBlockedSlotRepo{db: db}
}

func (r *postgresBlockedSlotRepo) Create(ctx context.Context, s *domain.BlockedTimeSlot) (*domain.BlockedTimeSlot, error) {
	query := `INSERT INTO blocked_time_slots (user_id, source_type, source_id, start_datetime, end_datetime, reason, is_flexible, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`
	now := time.Now()
	model := BlockedTimeSlotFromDomain(s)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.SourceType, model.SourceID, model.StartDatetime, model.EndDatetime, model.Reason, model.IsFlexible, now).Scan(&model.ID, &model.CreatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresBlockedSlotRepo) GetByID(ctx context.Context, id int) (*domain.BlockedTimeSlot, error) {
	var model BlockedTimeSlotModel
	if err := r.db.GetContext(ctx, &model, `SELECT * FROM blocked_time_slots WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresBlockedSlotRepo) GetByUserAndTimeRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.BlockedTimeSlot, error) {
	var models []BlockedTimeSlotModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM blocked_time_slots WHERE user_id = $1 AND start_datetime < $3 AND end_datetime > $2 ORDER BY start_datetime`, userID, start, end); err != nil {
		return nil, err
	}
	result := make([]*domain.BlockedTimeSlot, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresBlockedSlotRepo) GetByUserAndDate(ctx context.Context, userID int, date time.Time) ([]*domain.BlockedTimeSlot, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	dayEnd := dayStart.Add(24 * time.Hour)
	return r.GetByUserAndTimeRange(ctx, userID, dayStart, dayEnd)
}

func (r *postgresBlockedSlotRepo) GetBySource(ctx context.Context, sourceType string, sourceID int) ([]*domain.BlockedTimeSlot, error) {
	var models []BlockedTimeSlotModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM blocked_time_slots WHERE source_type = $1 AND source_id = $2`, sourceType, sourceID); err != nil {
		return nil, err
	}
	result := make([]*domain.BlockedTimeSlot, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresBlockedSlotRepo) Update(ctx context.Context, s *domain.BlockedTimeSlot) error {
	model := BlockedTimeSlotFromDomain(s)
	_, err := r.db.ExecContext(ctx, `UPDATE blocked_time_slots SET start_datetime = $1, end_datetime = $2, reason = $3, is_flexible = $4 WHERE id = $5`, model.StartDatetime, model.EndDatetime, model.Reason, model.IsFlexible, model.ID)
	return err
}

func (r *postgresBlockedSlotRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blocked_time_slots WHERE id = $1`, id)
	return err
}

func (r *postgresBlockedSlotRepo) DeleteBySource(ctx context.Context, sourceType string, sourceID int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM blocked_time_slots WHERE source_type = $1 AND source_id = $2`, sourceType, sourceID)
	return err
}

func (r *postgresBlockedSlotRepo) DeleteOld(ctx context.Context, before time.Time) (int, error) {
	result, err := r.db.ExecContext(ctx, `DELETE FROM blocked_time_slots WHERE end_datetime < $1`, before)
	if err != nil {
		return 0, err
	}
	rows, _ := result.RowsAffected()
	return int(rows), nil
}
