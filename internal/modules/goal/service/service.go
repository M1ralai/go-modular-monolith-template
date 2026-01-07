package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/dto"
)

type GoalService interface {
	Create(ctx context.Context, req *dto.CreateGoalRequest, userID int) (*dto.GoalResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.GoalResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.GoalResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateGoalRequest, userID int) (*dto.GoalResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
