package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/domain"
	"github.com/jmoiron/sqlx"
)

type postgresCalendarIntegrationRepo struct{ db *sqlx.DB }

func NewCalendarIntegrationRepository(db *sqlx.DB) CalendarIntegrationRepository {
	return &postgresCalendarIntegrationRepo{db: db}
}

func (r *postgresCalendarIntegrationRepo) Create(ctx context.Context, c *domain.CalendarIntegration) (*domain.CalendarIntegration, error) {
	query := `INSERT INTO calendar_integrations (user_id, provider, access_token, refresh_token, expires_at, calendar_id, is_active, last_sync_at, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := CalendarIntegrationFromDomain(c)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.Provider, model.AccessToken, model.RefreshToken, model.ExpiresAt, model.CalendarID, model.IsActive, model.LastSyncAt, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresCalendarIntegrationRepo) GetByID(ctx context.Context, id int) (*domain.CalendarIntegration, error) {
	var model CalendarIntegrationModel
	if err := r.db.GetContext(ctx, &model, `SELECT * FROM calendar_integrations WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresCalendarIntegrationRepo) GetByUserAndProvider(ctx context.Context, userID int, provider string) (*domain.CalendarIntegration, error) {
	var model CalendarIntegrationModel
	if err := r.db.GetContext(ctx, &model, `SELECT * FROM calendar_integrations WHERE user_id = $1 AND provider = $2`, userID, provider); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresCalendarIntegrationRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.CalendarIntegration, error) {
	var models []CalendarIntegrationModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM calendar_integrations WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}
	result := make([]*domain.CalendarIntegration, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresCalendarIntegrationRepo) GetActiveByProvider(ctx context.Context, provider string) ([]*domain.CalendarIntegration, error) {
	var models []CalendarIntegrationModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM calendar_integrations WHERE provider = $1 AND is_active = true`, provider); err != nil {
		return nil, err
	}
	result := make([]*domain.CalendarIntegration, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresCalendarIntegrationRepo) Update(ctx context.Context, c *domain.CalendarIntegration) error {
	model := CalendarIntegrationFromDomain(c)
	_, err := r.db.ExecContext(ctx, `UPDATE calendar_integrations SET access_token = $1, refresh_token = $2, expires_at = $3, calendar_id = $4, is_active = $5, last_sync_at = $6, updated_at = $7 WHERE id = $8`, model.AccessToken, model.RefreshToken, model.ExpiresAt, model.CalendarID, model.IsActive, model.LastSyncAt, time.Now(), model.ID)
	return err
}

func (r *postgresCalendarIntegrationRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM calendar_integrations WHERE id = $1`, id)
	return err
}

type postgresRecurringEventRepo struct{ db *sqlx.DB }

func NewRecurringEventRepository(db *sqlx.DB) RecurringEventRepository {
	return &postgresRecurringEventRepo{db: db}
}

func (r *postgresRecurringEventRepo) Create(ctx context.Context, e *domain.RecurringEvent) (*domain.RecurringEvent, error) {
	query := `INSERT INTO recurring_events (user_id, course_schedule_id, title, day_of_week, start_time, end_time, location, start_date, end_date, google_event_id, apple_event_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id, created_at`
	now := time.Now()
	var loc, gID, aID *string
	if e.Location != "" {
		loc = &e.Location
	}
	if e.GoogleEventID != "" {
		gID = &e.GoogleEventID
	}
	if e.AppleEventID != "" {
		aID = &e.AppleEventID
	}
	var id int
	var createdAt time.Time
	err := r.db.QueryRowxContext(ctx, query, e.UserID, e.CourseScheduleID, e.Title, e.DayOfWeek, e.StartTime, e.EndTime, loc, e.StartDate, e.EndDate, gID, aID, now).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	e.ID = id
	e.CreatedAt = createdAt
	return e, nil
}

func (r *postgresRecurringEventRepo) GetByID(ctx context.Context, id int) (*domain.RecurringEvent, error) {
	var model RecurringEventModel
	if err := r.db.GetContext(ctx, &model, `SELECT id, user_id, course_schedule_id, title, day_of_week, start_time::text, end_time::text, location, start_date, end_date, google_event_id, apple_event_id, last_synced_at, created_at FROM recurring_events WHERE id = $1`, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRecurringEventRepo) GetByUserID(ctx context.Context, userID int) ([]*domain.RecurringEvent, error) {
	var models []RecurringEventModel
	if err := r.db.SelectContext(ctx, &models, `SELECT id, user_id, course_schedule_id, title, day_of_week, start_time::text, end_time::text, location, start_date, end_date, google_event_id, apple_event_id, last_synced_at, created_at FROM recurring_events WHERE user_id = $1`, userID); err != nil {
		return nil, err
	}
	result := make([]*domain.RecurringEvent, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresRecurringEventRepo) GetByCourseScheduleID(ctx context.Context, scheduleID int) ([]*domain.RecurringEvent, error) {
	var models []RecurringEventModel
	if err := r.db.SelectContext(ctx, &models, `SELECT id, user_id, course_schedule_id, title, day_of_week, start_time::text, end_time::text, location, start_date, end_date, google_event_id, apple_event_id, last_synced_at, created_at FROM recurring_events WHERE course_schedule_id = $1`, scheduleID); err != nil {
		return nil, err
	}
	result := make([]*domain.RecurringEvent, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresRecurringEventRepo) Update(ctx context.Context, e *domain.RecurringEvent) error {
	var loc, gID, aID *string
	if e.Location != "" {
		loc = &e.Location
	}
	if e.GoogleEventID != "" {
		gID = &e.GoogleEventID
	}
	if e.AppleEventID != "" {
		aID = &e.AppleEventID
	}
	_, err := r.db.ExecContext(ctx, `UPDATE recurring_events SET title = $1, day_of_week = $2, start_time = $3, end_time = $4, location = $5, start_date = $6, end_date = $7, google_event_id = $8, apple_event_id = $9, last_synced_at = $10 WHERE id = $11`, e.Title, e.DayOfWeek, e.StartTime, e.EndTime, loc, e.StartDate, e.EndDate, gID, aID, e.LastSyncedAt, e.ID)
	return err
}

func (r *postgresRecurringEventRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM recurring_events WHERE id = $1`, id)
	return err
}

type postgresSyncQueueRepo struct{ db *sqlx.DB }

func NewSyncQueueRepository(db *sqlx.DB) SyncQueueRepository { return &postgresSyncQueueRepo{db: db} }

func (r *postgresSyncQueueRepo) Create(ctx context.Context, s *domain.SyncQueue) (*domain.SyncQueue, error) {
	query := `INSERT INTO calendar_sync_queue (user_id, event_id, provider, action, status, retry_count, error_message, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at`
	now := time.Now()
	var errMsg *string
	if s.ErrorMessage != "" {
		errMsg = &s.ErrorMessage
	}
	var id int
	var createdAt time.Time
	err := r.db.QueryRowxContext(ctx, query, s.UserID, s.EventID, s.Provider, s.Action, "pending", 0, errMsg, now).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	s.ID = id
	s.CreatedAt = createdAt
	return s, nil
}

func (r *postgresSyncQueueRepo) GetPending(ctx context.Context, limit int) ([]*domain.SyncQueue, error) {
	var models []SyncQueueModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM calendar_sync_queue WHERE status = 'pending' ORDER BY created_at LIMIT $1`, limit); err != nil {
		return nil, err
	}
	result := make([]*domain.SyncQueue, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresSyncQueueRepo) GetByEventID(ctx context.Context, eventID int) ([]*domain.SyncQueue, error) {
	var models []SyncQueueModel
	if err := r.db.SelectContext(ctx, &models, `SELECT * FROM calendar_sync_queue WHERE event_id = $1`, eventID); err != nil {
		return nil, err
	}
	result := make([]*domain.SyncQueue, len(models))
	for i, m := range models {
		result[i] = m.ToDomain()
	}
	return result, nil
}

func (r *postgresSyncQueueRepo) UpdateStatus(ctx context.Context, id int, status string, errorMsg string) error {
	var errPtr *string
	if errorMsg != "" {
		errPtr = &errorMsg
	}
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `UPDATE calendar_sync_queue SET status = $1, error_message = $2, synced_at = $3 WHERE id = $4`, status, errPtr, now, id)
	return err
}

func (r *postgresSyncQueueRepo) IncrementRetry(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `UPDATE calendar_sync_queue SET retry_count = retry_count + 1 WHERE id = $1`, id)
	return err
}

func (r *postgresSyncQueueRepo) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM calendar_sync_queue WHERE id = $1`, id)
	return err
}
