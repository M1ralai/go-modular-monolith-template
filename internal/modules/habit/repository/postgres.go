package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) HabitRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, habit *domain.Habit) (*domain.Habit, error) {
	query := `
		INSERT INTO habits (user_id, life_area_id, name, icon, description, frequency, frequency_config, target_count, time_of_day, reminder_time, current_streak, longest_streak, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, created_at, updated_at
	`
	now := time.Now()
	model := FromDomain(habit)
	err := r.db.QueryRowxContext(ctx, query,
		model.UserID, model.LifeAreaID, model.Name, model.Icon, model.Description,
		model.Frequency, model.FrequencyConfig, model.TargetCount, model.TimeOfDay, model.ReminderTime,
		0, 0, true, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Habit, error) {
	query := `SELECT id, user_id, life_area_id, name, icon, description, frequency, frequency_config, target_count, time_of_day, reminder_time, current_streak, longest_streak, is_active, created_at, updated_at FROM habits WHERE id = $1`
	var model HabitModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Habit, error) {
	query := `SELECT id, user_id, life_area_id, name, icon, description, frequency, frequency_config, target_count, time_of_day, reminder_time, current_streak, longest_streak, is_active, created_at, updated_at FROM habits WHERE user_id = $1 ORDER BY created_at DESC`
	var models []HabitModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}
	habits := make([]*domain.Habit, len(models))
	for i, m := range models {
		habits[i] = m.ToDomain()
	}
	return habits, nil
}

func (r *postgresRepository) GetActiveHabits(ctx context.Context, userID int) ([]*domain.Habit, error) {
	query := `SELECT id, user_id, life_area_id, name, icon, description, frequency, frequency_config, target_count, time_of_day, reminder_time, current_streak, longest_streak, is_active, created_at, updated_at FROM habits WHERE user_id = $1 AND is_active = true ORDER BY name`
	var models []HabitModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}
	habits := make([]*domain.Habit, len(models))
	for i, m := range models {
		habits[i] = m.ToDomain()
	}
	return habits, nil
}

func (r *postgresRepository) Update(ctx context.Context, habit *domain.Habit) error {
	query := `UPDATE habits SET name = $1, icon = $2, description = $3, frequency = $4, frequency_config = $5, target_count = $6, time_of_day = $7, reminder_time = $8, current_streak = $9, longest_streak = $10, is_active = $11, life_area_id = $12, updated_at = $13 WHERE id = $14`
	model := FromDomain(habit)
	_, err := r.db.ExecContext(ctx, query,
		model.Name, model.Icon, model.Description, model.Frequency, model.FrequencyConfig,
		model.TargetCount, model.TimeOfDay, model.ReminderTime,
		model.CurrentStreak, model.LongestStreak, model.IsActive, model.LifeAreaID, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM habits WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) LogHabit(ctx context.Context, habitID int, logDate time.Time, count int, notes string) error {
	// Check if a log already exists for this date
	existingLog, err := r.GetLogsForDate(ctx, habitID, logDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	
	// If log exists and is completed or skipped, don't allow update
	if existingLog != nil {
		if existingLog.IsCompleted {
			return errors.New("habit already completed today")
		}
		if existingLog.Skipped {
			return errors.New("habit already skipped today")
		}
	}
	
	query := `
		INSERT INTO habit_logs (habit_id, log_date, count, notes, is_completed, skipped, created_at)
		VALUES ($1, $2, $3, $4, $5, false, $6)
		ON CONFLICT (habit_id, log_date) DO UPDATE SET count = $3, notes = $4, is_completed = $5, skipped = false
	`
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}
	_, err = r.db.ExecContext(ctx, query, habitID, logDate.Format("2006-01-02"), count, notesPtr, count > 0, time.Now())
	return err
}

func (r *postgresRepository) SkipHabit(ctx context.Context, habitID int, logDate time.Time, notes string) error {
	// Check if a log already exists for this date
	existingLog, err := r.GetLogsForDate(ctx, habitID, logDate)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	
	// If log exists and is completed or skipped, don't allow update
	if existingLog != nil {
		if existingLog.IsCompleted {
			return errors.New("habit already completed today - cannot skip")
		}
		if existingLog.Skipped {
			return errors.New("habit already skipped today")
		}
	}
	
	query := `
		INSERT INTO habit_logs (habit_id, log_date, skipped, notes, is_completed, created_at)
		VALUES ($1, $2, true, $3, false, $4)
		ON CONFLICT (habit_id, log_date) DO UPDATE SET skipped = true, notes = $3, is_completed = false
	`
	var notesPtr *string
	if notes != "" {
		notesPtr = &notes
	}
	_, err = r.db.ExecContext(ctx, query, habitID, logDate.Format("2006-01-02"), notesPtr, time.Now())
	return err
}

func (r *postgresRepository) GetLogsForDate(ctx context.Context, habitID int, date time.Time) (*HabitLogModel, error) {
	query := `SELECT id, habit_id, log_date, count, notes, is_completed, skipped, created_at FROM habit_logs WHERE habit_id = $1 AND log_date = $2`
	var model HabitLogModel
	err := r.db.GetContext(ctx, &model, query, habitID, date.Format("2006-01-02"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &model, nil
}

func (r *postgresRepository) GetLogsByDateRange(ctx context.Context, habitID int, start, end time.Time) ([]*HabitLogModel, error) {
	query := `SELECT id, habit_id, log_date, count, notes, is_completed, created_at FROM habit_logs WHERE habit_id = $1 AND log_date BETWEEN $2 AND $3 ORDER BY log_date`
	var models []*HabitLogModel
	err := r.db.SelectContext(ctx, &models, query, habitID, start, end)
	if err != nil {
		return nil, err
	}
	return models, nil
}

func (r *postgresRepository) HasLogForToday(ctx context.Context, habitID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = $1 AND log_date = CURRENT_DATE AND is_completed = true)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, habitID)
	return exists, err
}

func (r *postgresRepository) HasSkippedToday(ctx context.Context, habitID int) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = $1 AND log_date = CURRENT_DATE AND skipped = true)`
	var exists bool
	err := r.db.GetContext(ctx, &exists, query, habitID)
	return exists, err
}
