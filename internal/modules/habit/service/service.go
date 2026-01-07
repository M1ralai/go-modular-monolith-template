package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/dto"
)

type HabitService interface {
	Create(ctx context.Context, req *dto.CreateHabitRequest, userID int) (*dto.HabitResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.HabitResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.HabitResponse, error)
	GetActive(ctx context.Context, userID int) ([]*dto.HabitResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateHabitRequest, userID int) (*dto.HabitResponse, error)
	Delete(ctx context.Context, id, userID int) error
	LogHabit(ctx context.Context, id int, req *dto.LogHabitRequest, userID int) error
	Complete(ctx context.Context, id int, req *dto.LogHabitRequest, userID int) error
	SkipHabit(ctx context.Context, id int, userID int) error
}
