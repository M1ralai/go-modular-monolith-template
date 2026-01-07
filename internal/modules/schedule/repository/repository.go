package repository

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/domain"
)

type BlockedTimeSlotRepository interface {
	Create(ctx context.Context, slot *domain.BlockedTimeSlot) (*domain.BlockedTimeSlot, error)
	GetByID(ctx context.Context, id int) (*domain.BlockedTimeSlot, error)
	GetByUserAndTimeRange(ctx context.Context, userID int, start, end time.Time) ([]*domain.BlockedTimeSlot, error)
	GetByUserAndDate(ctx context.Context, userID int, date time.Time) ([]*domain.BlockedTimeSlot, error)
	GetBySource(ctx context.Context, sourceType string, sourceID int) ([]*domain.BlockedTimeSlot, error)
	Update(ctx context.Context, slot *domain.BlockedTimeSlot) error
	Delete(ctx context.Context, id int) error
	DeleteBySource(ctx context.Context, sourceType string, sourceID int) error
	DeleteOld(ctx context.Context, before time.Time) (int, error)
}
