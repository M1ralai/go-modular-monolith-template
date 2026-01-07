package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/repository"
)

type taskService struct {
	repo        repository.TaskRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewTaskService(repo repository.TaskRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) TaskService {
	return &taskService{
		repo:        repo,
		logger:      logger,
		broadcaster: broadcaster,
	}
}

func (s *taskService) Create(ctx context.Context, req *dto.CreateTaskRequest, userID int) (*dto.TaskResponse, error) {
	s.logger.Info("Creating task", map[string]interface{}{
		"user_id": userID,
		"title":   req.Title,
		"action":  "CREATE_TASK",
	})

	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}

	now := time.Now()
	task := &domain.Task{
		UserID:             userID,
		ParentTaskID:       req.ParentTaskID,
		Title:              req.Title,
		Description:        req.Description,
		DueDate:            req.DueDate,
		EstimatedStart:     req.EstimatedStart,
		EstimatedEnd:       req.EstimatedEnd,
		Priority:           priority,
		IsCompleted:        false,
		ProgressPercentage: 0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	created, err := s.repo.Create(ctx, task)
	if err != nil {
		s.logger.Error("Failed to create task", err, map[string]interface{}{
			"user_id": userID,
			"title":   req.Title,
			"action":  "CREATE_TASK_FAILED",
		})
		return nil, err
	}

	s.logger.Info("Task created successfully", map[string]interface{}{
		"user_id": userID,
		"task_id": created.ID,
		"action":  "CREATE_TASK_SUCCESS",
	})

	response := dto.ToTaskResponse(created, 0, 0)
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventTaskCreated, map[string]interface{}{
			"task_id": created.ID,
			"task":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventTaskCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *taskService) GetByID(ctx context.Context, id, userID int) (*dto.TaskResponse, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	if task.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	total, completed, _ := s.repo.CountSubtasks(ctx, id)
	return dto.ToTaskResponse(task, total, completed), nil
}

func (s *taskService) GetAll(ctx context.Context, userID int) ([]*dto.TaskResponse, error) {
	tasks, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		total, completed, _ := s.repo.CountSubtasks(ctx, task.ID)
		result[i] = dto.ToTaskResponse(task, total, completed)
	}

	return result, nil
}

func (s *taskService) GetParentTasks(ctx context.Context, userID int) ([]*dto.TaskResponse, error) {
	tasks, err := s.repo.GetParentTasks(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		total, completed, _ := s.repo.CountSubtasks(ctx, task.ID)
		result[i] = dto.ToTaskResponse(task, total, completed)
	}

	return result, nil
}

func (s *taskService) GetSubtasks(ctx context.Context, parentID, userID int) ([]*dto.TaskResponse, error) {
	parent, err := s.repo.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if parent == nil {
		return nil, errors.New("parent task not found")
	}
	if parent.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	subtasks, err := s.repo.GetSubtasks(ctx, parentID)
	if err != nil {
		return nil, err
	}

	return dto.ToTaskResponseList(subtasks), nil
}

func (s *taskService) Update(ctx context.Context, id int, req *dto.UpdateTaskRequest, userID int) (*dto.TaskResponse, error) {
	s.logger.Info("Updating task", map[string]interface{}{
		"user_id": userID,
		"task_id": id,
		"action":  "UPDATE_TASK",
	})

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	if task.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}
	if req.EstimatedStart != nil {
		task.EstimatedStart = req.EstimatedStart
	}
	if req.EstimatedEnd != nil {
		task.EstimatedEnd = req.EstimatedEnd
	}
	if req.ActualStart != nil {
		task.ActualStart = req.ActualStart
	}
	if req.ActualEnd != nil {
		task.ActualEnd = req.ActualEnd
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	wasCompleted := task.IsCompleted
	if req.IsCompleted != nil {
		task.IsCompleted = *req.IsCompleted
		if *req.IsCompleted && !wasCompleted {
			// Mark as completed if transitioning from incomplete to complete
			task.MarkCompleted()
		} else if !*req.IsCompleted && wasCompleted {
			// Unmark completion
			task.CompletedAt = nil
		}
	}
	if req.CompletedAt != nil {
		task.CompletedAt = req.CompletedAt
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		s.logger.Error("Failed to update task", err, map[string]interface{}{
			"user_id": userID,
			"task_id": id,
			"action":  "UPDATE_TASK_FAILED",
		})
		return nil, err
	}

	s.logger.Info("Task updated successfully", map[string]interface{}{
		"user_id": userID,
		"task_id": id,
		"action":  "UPDATE_TASK_SUCCESS",
	})

	total, completed, _ := s.repo.CountSubtasks(ctx, id)
	response := dto.ToTaskResponse(task, total, completed)
	
	if s.broadcaster != nil {
		// If task was just completed, send completion event
		if req.IsCompleted != nil && *req.IsCompleted && !wasCompleted {
			s.broadcaster.Publish(userID, notification.EventTaskCompleted, map[string]interface{}{
				"task_id": id,
				"task":    response,
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventTaskCompleted,
				"user_id":    userID,
				"entity_id":  id,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}
		
		// Always send update event
		s.broadcaster.Publish(userID, notification.EventTaskUpdated, map[string]interface{}{
			"task_id": id,
			"task":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventTaskUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *taskService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting task", map[string]interface{}{
		"user_id": userID,
		"task_id": id,
		"action":  "DELETE_TASK",
	})

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if task == nil {
		return errors.New("task not found")
	}
	if task.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete task", err, map[string]interface{}{
			"user_id": userID,
			"task_id": id,
			"action":  "DELETE_TASK_FAILED",
		})
		return err
	}

	s.logger.Info("Task deleted successfully", map[string]interface{}{
		"user_id": userID,
		"task_id": id,
		"action":  "DELETE_TASK_SUCCESS",
	})

	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventTaskDeleted, map[string]interface{}{
			"task_id": id,
			"title":   task.Title,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventTaskDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return nil
}

func (s *taskService) CompleteSubtask(ctx context.Context, subtaskID, userID int) error {
	s.logger.Info("Completing subtask", map[string]interface{}{
		"user_id":    userID,
		"subtask_id": subtaskID,
		"action":     "COMPLETE_SUBTASK",
	})

	subtask, err := s.repo.GetByID(ctx, subtaskID)
	if err != nil {
		return err
	}
	if subtask == nil {
		return errors.New("subtask not found")
	}
	if subtask.UserID != userID {
		return errors.New("unauthorized")
	}

	if subtask.IsCompleted {
		return nil
	}

	subtask.MarkCompleted()

	if err := s.repo.Update(ctx, subtask); err != nil {
		s.logger.Error("Failed to complete subtask", err, map[string]interface{}{
			"user_id":    userID,
			"subtask_id": subtaskID,
			"action":     "COMPLETE_SUBTASK_FAILED",
		})
		return err
	}

	s.logger.Info("Subtask completed successfully", map[string]interface{}{
		"user_id":    userID,
		"subtask_id": subtaskID,
		"action":     "COMPLETE_SUBTASK_SUCCESS",
	})

	total, completed, _ := s.repo.CountSubtasks(ctx, subtaskID)
	response := dto.ToTaskResponse(subtask, total, completed)
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventTaskCompleted, map[string]interface{}{
			"task_id": subtaskID,
			"task":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventTaskCompleted,
			"user_id":    userID,
			"entity_id":  subtaskID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	if subtask.ParentTaskID != nil {
		return s.checkAndCompleteParent(ctx, *subtask.ParentTaskID, userID)
	}

	return nil
}

func (s *taskService) checkAndCompleteParent(ctx context.Context, parentID, userID int) error {
	total, completed, err := s.repo.CountSubtasks(ctx, parentID)
	if err != nil {
		return err
	}

	if total == 0 {
		return nil
	}

	parent, err := s.repo.GetByID(ctx, parentID)
	if err != nil || parent == nil {
		return err
	}

	parent.ProgressPercentage = float64(completed) / float64(total) * 100
	parent.UpdatedAt = time.Now()

	if completed == total {
		parent.MarkCompleted()

		s.logger.Info("Parent task auto-completed", map[string]interface{}{
			"user_id":   userID,
			"parent_id": parentID,
			"action":    "AUTO_COMPLETE_PARENT",
		})

		if s.broadcaster != nil {
			parentTotal, parentCompleted, _ := s.repo.CountSubtasks(ctx, parentID)
			parentResponse := dto.ToTaskResponse(parent, parentTotal, parentCompleted)
			
			s.broadcaster.Publish(userID, notification.EventTaskCompleted, map[string]interface{}{
				"task_id":        parentID,
				"task":           parentResponse,
				"auto_completed": true,
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventTaskCompleted,
				"user_id":    userID,
				"entity_id":  parentID,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}
	}

	return s.repo.Update(ctx, parent)
}

func (s *taskService) GetStats(ctx context.Context, userID int) (*dto.TaskStatsResponse, error) {
	completedToday, dueToday, dueTomorrow, overdue, err := s.repo.GetStats(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get task stats", err, map[string]interface{}{
			"user_id": userID,
			"action":  "GET_TASK_STATS_FAILED",
		})
		return nil, err
	}

	return &dto.TaskStatsResponse{
		CompletedToday: completedToday,
		DueToday:       dueToday,
		DueTomorrow:    dueTomorrow,
		Overdue:        overdue,
	}, nil
}
