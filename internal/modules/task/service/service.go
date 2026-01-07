package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/dto"
)

type TaskService interface {
	Create(ctx context.Context, req *dto.CreateTaskRequest, userID int) (*dto.TaskResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.TaskResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.TaskResponse, error)
	GetParentTasks(ctx context.Context, userID int) ([]*dto.TaskResponse, error)
	GetSubtasks(ctx context.Context, parentID, userID int) ([]*dto.TaskResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateTaskRequest, userID int) (*dto.TaskResponse, error)
	Delete(ctx context.Context, id, userID int) error
	CompleteSubtask(ctx context.Context, subtaskID, userID int) error
	GetStats(ctx context.Context, userID int) (*dto.TaskStatsResponse, error)
}
