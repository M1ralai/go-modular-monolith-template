package repository

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/domain"
)

type EventRepository interface {
	Create(ctx context.Context, event *domain.Event) (*domain.Event, error)
	GetByID(ctx context.Context, id int) (*domain.Event, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Event, error)
	GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.Event, error)
	Update(ctx context.Context, event *domain.Event) error
	Delete(ctx context.Context, id int) error
}
