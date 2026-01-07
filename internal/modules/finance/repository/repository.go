package repository

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	GetByID(ctx context.Context, id int) (*domain.Transaction, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Transaction, error)
	GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.Transaction, error)
	GetSummary(ctx context.Context, userID int, start, end time.Time) (income float64, expense float64, err error)
	Update(ctx context.Context, tx *domain.Transaction) error
	Delete(ctx context.Context, id int) error
}
