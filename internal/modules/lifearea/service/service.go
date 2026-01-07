package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/dto"
)

type LifeAreaService interface {
	Create(ctx context.Context, req *dto.CreateLifeAreaRequest, userID int) (*dto.LifeAreaResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.LifeAreaResponse, error)
	GetByUserID(ctx context.Context, userID int) ([]*dto.LifeAreaResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateLifeAreaRequest, userID int) (*dto.LifeAreaResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
