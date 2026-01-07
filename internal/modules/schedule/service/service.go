package service

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/dto"
)

type ScheduleService interface {
	// Conflict Detection
	CheckConflict(ctx context.Context, userID int, start, end time.Time) (*dto.ConflictResponse, error)

	// Free Slots
	GetFreeSlots(ctx context.Context, userID int, date time.Time, durationMinutes int) ([]*dto.TimeSlotResponse, error)

	// Blocked Slots
	GetBlockedSlots(ctx context.Context, userID int, date time.Time) ([]*dto.BlockedSlotResponse, error)
	CreateBlockedSlot(ctx context.Context, userID int, req *dto.CreateBlockedSlotRequest) (*dto.BlockedSlotResponse, error)
	DeleteBlockedSlot(ctx context.Context, id, userID int) error

	// Event Generation
	GenerateEventsForSchedule(ctx context.Context, userID int, req *dto.GenerateEventsRequest) (*dto.GenerateEventsResponse, error)
}
