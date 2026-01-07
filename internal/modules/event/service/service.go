package service

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/dto"
)

type EventService interface {
	Create(ctx context.Context, req *dto.CreateEventRequest, userID int) (*dto.EventResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.EventResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.EventResponse, error)
	GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*dto.EventResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateEventRequest, userID int) (*dto.EventResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
