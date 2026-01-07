package repository

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/domain"
)

type JournalRepository interface {
	Create(ctx context.Context, entry *domain.JournalEntry) (*domain.JournalEntry, error)
	GetByID(ctx context.Context, id int) (*domain.JournalEntry, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.JournalEntry, error)
	GetByDate(ctx context.Context, userID int, date time.Time) (*domain.JournalEntry, error)
	GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.JournalEntry, error)
	Update(ctx context.Context, entry *domain.JournalEntry) error
	Delete(ctx context.Context, id int) error
}
