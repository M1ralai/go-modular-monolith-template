package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/domain"
)

type GoalRepository interface {
	Create(ctx context.Context, goal *domain.Goal) (*domain.Goal, error)
	GetByID(ctx context.Context, id int) (*domain.Goal, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Goal, error)
	Update(ctx context.Context, goal *domain.Goal) error
	Delete(ctx context.Context, id int) error
	CountMilestones(ctx context.Context, goalID int) (total int, completed int, err error)
}
