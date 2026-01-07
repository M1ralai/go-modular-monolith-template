package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct{ db *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) JournalRepository { return &postgresRepository{db: db} }

func (r *postgresRepository) Create(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error) {
	query := `INSERT INTO journal_entries (user_id, entry_date, content, mood, energy_level, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := FromDomain(entry)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.EntryDate, model.Content, model.Mood, model.EnergyLevel, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.JournalEntry, error) {
	query := `SELECT id, user_id, entry_date, content, mood, energy_level, created_at, updated_at FROM journal_entries WHERE id = $1`
	var model JournalModel
	if err := r.db.GetContext(ctx, &model, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.JournalEntry, error) {
	query := `SELECT id, user_id, entry_date, content, mood, energy_level, created_at, updated_at FROM journal_entries WHERE user_id = $1 ORDER BY entry_date DESC`
	var models []JournalModel
	if err := r.db.SelectContext(ctx, &models, query, userID); err != nil {
		return nil, err
	}
	entries := make([]*domain.JournalEntry, len(models))
	for i, m := range models {
		entries[i] = m.ToDomain()
	}
	return entries, nil
}

func (r *postgresRepository) GetByDate(ctx context.Context, userID int, date time.Time) (*domain.JournalEntry, error) {
	query := `SELECT id, user_id, entry_date, content, mood, energy_level, created_at, updated_at FROM journal_entries WHERE user_id = $1 AND entry_date = $2`
	var model JournalModel
	if err := r.db.GetContext(ctx, &model, query, userID, date); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.JournalEntry, error) {
	query := `SELECT id, user_id, entry_date, content, mood, energy_level, created_at, updated_at FROM journal_entries WHERE user_id = $1 AND entry_date BETWEEN $2 AND $3 ORDER BY entry_date DESC`
	var models []JournalModel
	if err := r.db.SelectContext(ctx, &models, query, userID, start, end); err != nil {
		return nil, err
	}
	entries := make([]*domain.JournalEntry, len(models))
	for i, m := range models {
		entries[i] = m.ToDomain()
	}
	return entries, nil
}

func (r *postgresRepository) Update(ctx context.Context, entry *domain.JournalEntry) error {
	query := `UPDATE journal_entries SET content = $1, mood = $2, energy_level = $3, updated_at = $4 WHERE id = $5`
	model := FromDomain(entry)
	_, err := r.db.ExecContext(ctx, query, model.Content, model.Mood, model.EnergyLevel, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM journal_entries WHERE id = $1`, id)
	return err
}
