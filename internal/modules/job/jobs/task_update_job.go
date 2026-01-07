package jobimpl

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/repository"
)

// TaskUpdateJob updates a task asynchronously
type TaskUpdateJob struct {
	jobs.BaseJob
	logger      *logger.ZapLogger
	repo        repository.TaskRepository
	broadcaster *notifService.Broadcaster
	taskID      int
	userID      int
	updates     *dto.UpdateTaskRequest
}

// NewTaskUpdateJob creates a new task update job
func NewTaskUpdateJob(
	logger *logger.ZapLogger,
	repo repository.TaskRepository,
	broadcaster *notifService.Broadcaster,
	taskID, userID int,
	updates *dto.UpdateTaskRequest,
) *TaskUpdateJob {
	return &TaskUpdateJob{
		BaseJob:     jobs.NewBaseJob("task_update", "", 30*time.Second, nil),
		logger:      logger,
		repo:        repo,
		broadcaster: broadcaster,
		taskID:      taskID,
		userID:      userID,
		updates:     updates,
	}
}

func (j *TaskUpdateJob) Execute(ctx context.Context) error {
	j.logger.Info("Task update job started", map[string]interface{}{
		"task_id": j.taskID,
		"user_id": j.userID,
		"action":  "TASK_UPDATE_JOB_STARTED",
	})

	// Get task
	task, err := j.repo.GetByID(ctx, j.taskID)
	if err != nil {
		j.logger.Error("Failed to get task in job", err, map[string]interface{}{
			"task_id": j.taskID,
			"user_id": j.userID,
			"action":  "TASK_UPDATE_JOB_FAILED",
		})
		return err
	}

	if task == nil {
		err := context.DeadlineExceeded // Use as "not found" error
		j.logger.Error("Task not found in job", err, map[string]interface{}{
			"task_id": j.taskID,
			"user_id": j.userID,
			"action":  "TASK_UPDATE_JOB_NOT_FOUND",
		})
		return err
	}

	if task.UserID != j.userID {
		err := context.DeadlineExceeded // Use as "unauthorized" error
		j.logger.Error("Unauthorized task update in job", err, map[string]interface{}{
			"task_id": j.taskID,
			"user_id": j.userID,
			"action":  "TASK_UPDATE_JOB_UNAUTHORIZED",
		})
		return err
	}

	// Apply updates
	if j.updates.Title != nil {
		task.Title = *j.updates.Title
	}
	if j.updates.Description != nil {
		task.Description = *j.updates.Description
	}
	if j.updates.DueDate != nil {
		task.DueDate = j.updates.DueDate
	}
	if j.updates.Priority != nil {
		task.Priority = *j.updates.Priority
	}
	if j.updates.IsCompleted != nil {
		task.IsCompleted = *j.updates.IsCompleted
	}
	if j.updates.CompletedAt != nil {
		task.CompletedAt = j.updates.CompletedAt
	}

	task.UpdatedAt = time.Now()

	// Update in database
	if err := j.repo.Update(ctx, task); err != nil {
		j.logger.Error("Failed to update task in job", err, map[string]interface{}{
			"task_id": j.taskID,
			"user_id": j.userID,
			"action":  "TASK_UPDATE_JOB_UPDATE_FAILED",
		})
		return err
	}

	// Broadcast WebSocket message
	if j.broadcaster != nil {
		total, completed, _ := j.repo.CountSubtasks(ctx, j.taskID)
		response := dto.ToTaskResponse(task, total, completed)

		j.broadcaster.Publish(j.userID, "task.updated", map[string]interface{}{
			"task": response,
		})

		// If task was completed, send completion notification
		if j.updates.IsCompleted != nil && *j.updates.IsCompleted {
			j.broadcaster.Publish(j.userID, "task.completed", map[string]interface{}{
				"task_id": j.taskID,
				"title":   task.Title,
			})
		}
	}

	j.logger.Info("Task update job completed", map[string]interface{}{
		"task_id": j.taskID,
		"user_id": j.userID,
		"action":  "TASK_UPDATE_JOB_COMPLETED",
	})

	return nil
}
