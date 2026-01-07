package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/dto"
)

type ComponentService interface {
	GetAll(ctx context.Context) ([]*dto.CreateComponentRequest, error)
	GetByID(ctx context.Context, id int) (*dto.CreateComponentRequest, error)
}
