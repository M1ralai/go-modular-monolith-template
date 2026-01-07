package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct{ db *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) TransactionRepository { return &postgresRepository{db: db} }

func (r *postgresRepository) Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	query := `INSERT INTO finance_transactions (user_id, amount, type, category, description, date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := FromDomain(tx)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.Amount, model.Type, model.Category, model.Description, model.Date, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Transaction, error) {
	query := `SELECT id, user_id, amount, type, category, description, date, created_at, updated_at FROM finance_transactions WHERE id = $1`
	var model TransactionModel
	if err := r.db.GetContext(ctx, &model, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Transaction, error) {
	query := `SELECT id, user_id, amount, type, category, description, date, created_at, updated_at FROM finance_transactions WHERE user_id = $1 ORDER BY date DESC`
	var models []TransactionModel
	if err := r.db.SelectContext(ctx, &models, query, userID); err != nil {
		return nil, err
	}
	txs := make([]*domain.Transaction, len(models))
	for i, m := range models {
		txs[i] = m.ToDomain()
	}
	return txs, nil
}

func (r *postgresRepository) GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.Transaction, error) {
	query := `SELECT id, user_id, amount, type, category, description, date, created_at, updated_at FROM finance_transactions WHERE user_id = $1 AND date BETWEEN $2 AND $3 ORDER BY date DESC`
	var models []TransactionModel
	if err := r.db.SelectContext(ctx, &models, query, userID, start, end); err != nil {
		return nil, err
	}
	txs := make([]*domain.Transaction, len(models))
	for i, m := range models {
		txs[i] = m.ToDomain()
	}
	return txs, nil
}

func (r *postgresRepository) GetSummary(ctx context.Context, userID int, start, end time.Time) (income float64, expense float64, err error) {
	query := `SELECT COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income, COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense FROM finance_transactions WHERE user_id = $1 AND date BETWEEN $2 AND $3`
	err = r.db.QueryRowxContext(ctx, query, userID, start, end).Scan(&income, &expense)
	return
}

func (r *postgresRepository) Update(ctx context.Context, tx *domain.Transaction) error {
	query := `UPDATE finance_transactions SET amount = $1, type = $2, category = $3, description = $4, date = $5, updated_at = $6 WHERE id = $7`
	model := FromDomain(tx)
	_, err := r.db.ExecContext(ctx, query, model.Amount, model.Type, model.Category, model.Description, model.Date, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM finance_transactions WHERE id = $1`, id)
	return err
}
