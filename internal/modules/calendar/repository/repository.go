package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/domain"
)

type CalendarIntegrationRepository interface {
	Create(ctx context.Context, integration *domain.CalendarIntegration) (*domain.CalendarIntegration, error)
	GetByID(ctx context.Context, id int) (*domain.CalendarIntegration, error)
	GetByUserAndProvider(ctx context.Context, userID int, provider string) (*domain.CalendarIntegration, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.CalendarIntegration, error)
	GetActiveByProvider(ctx context.Context, provider string) ([]*domain.CalendarIntegration, error)
	Update(ctx context.Context, integration *domain.CalendarIntegration) error
	Delete(ctx context.Context, id int) error
}

type RecurringEventRepository interface {
	Create(ctx context.Context, event *domain.RecurringEvent) (*domain.RecurringEvent, error)
	GetByID(ctx context.Context, id int) (*domain.RecurringEvent, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.RecurringEvent, error)
	GetByCourseScheduleID(ctx context.Context, scheduleID int) ([]*domain.RecurringEvent, error)
	Update(ctx context.Context, event *domain.RecurringEvent) error
	Delete(ctx context.Context, id int) error
}

type SyncQueueRepository interface {
	Create(ctx context.Context, item *domain.SyncQueue) (*domain.SyncQueue, error)
	GetPending(ctx context.Context, limit int) ([]*domain.SyncQueue, error)
	GetByEventID(ctx context.Context, eventID int) ([]*domain.SyncQueue, error)
	UpdateStatus(ctx context.Context, id int, status string, errorMsg string) error
	IncrementRetry(ctx context.Context, id int) error
	Delete(ctx context.Context, id int) error
}
