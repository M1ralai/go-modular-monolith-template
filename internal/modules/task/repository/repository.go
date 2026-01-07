package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) (*domain.Task, error)
	GetByID(ctx context.Context, id int) (*domain.Task, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Task, error)
	GetSubtasks(ctx context.Context, parentID int) ([]*domain.Task, error)
	GetParentTasks(ctx context.Context, userID int) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id int) error
	CountSubtasks(ctx context.Context, parentID int) (total int, completed int, err error)
	GetStats(ctx context.Context, userID int) (completedToday, dueToday, dueTomorrow, overdue int, err error)
}
