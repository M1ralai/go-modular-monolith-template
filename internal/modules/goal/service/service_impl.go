package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
)

type goalService struct {
	repo        repository.GoalRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewGoalService(repo repository.GoalRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) GoalService {
	return &goalService{repo: repo, logger: logger, broadcaster: broadcaster}
}

func (s *goalService) Create(ctx context.Context, req *dto.CreateGoalRequest, userID int) (*dto.GoalResponse, error) {
	s.logger.Info("Creating goal", map[string]interface{}{"user_id": userID, "title": req.Title, "action": "CREATE_GOAL"})
	priority := req.Priority
	if priority == "" {
		priority = "medium"
	}
	now := time.Now()
	goal := &domain.Goal{UserID: userID, LifeAreaID: req.LifeAreaID, Title: req.Title, Description: req.Description, TargetDate: req.TargetDate, Priority: priority, CreatedAt: now, UpdatedAt: now}
	created, err := s.repo.Create(ctx, goal)
	if err != nil {
		s.logger.Error("Failed to create goal", err, map[string]interface{}{"user_id": userID, "action": "CREATE_GOAL_FAILED"})
		return nil, err
	}
	s.logger.Info("Goal created", map[string]interface{}{"user_id": userID, "goal_id": created.ID, "action": "CREATE_GOAL_SUCCESS"})
	
	response := dto.ToGoalResponse(created, 0, 0)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventGoalCreated, map[string]interface{}{
			"goal_id": created.ID,
			"goal":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventGoalCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return response, nil
}

func (s *goalService) GetByID(ctx context.Context, id, userID int) (*dto.GoalResponse, error) {
	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if goal == nil {
		return nil, errors.New("goal not found")
	}
	if goal.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	total, completed, _ := s.repo.CountMilestones(ctx, id)
	return dto.ToGoalResponse(goal, total, completed), nil
}

func (s *goalService) GetAll(ctx context.Context, userID int) ([]*dto.GoalResponse, error) {
	goals, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.GoalResponse, len(goals))
	for i, g := range goals {
		total, completed, _ := s.repo.CountMilestones(ctx, g.ID)
		result[i] = dto.ToGoalResponse(g, total, completed)
	}
	return result, nil
}

func (s *goalService) Update(ctx context.Context, id int, req *dto.UpdateGoalRequest, userID int) (*dto.GoalResponse, error) {
	s.logger.Info("Updating goal", map[string]interface{}{"user_id": userID, "goal_id": id, "action": "UPDATE_GOAL"})
	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if goal == nil {
		return nil, errors.New("goal not found")
	}
	if goal.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Title != nil {
		goal.Title = *req.Title
	}
	if req.Description != nil {
		goal.Description = *req.Description
	}
	if req.TargetDate != nil {
		goal.TargetDate = req.TargetDate
	}
	if req.Priority != nil {
		goal.Priority = *req.Priority
	}
	if req.LifeAreaID != nil {
		goal.LifeAreaID = req.LifeAreaID
	}
	if req.IsCompleted != nil && *req.IsCompleted {
		goal.MarkCompleted()
	}
	goal.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, goal); err != nil {
		s.logger.Error("Failed to update goal", err, map[string]interface{}{"user_id": userID, "goal_id": id, "action": "UPDATE_GOAL_FAILED"})
		return nil, err
	}
	s.logger.Info("Goal updated", map[string]interface{}{"user_id": userID, "goal_id": id, "action": "UPDATE_GOAL_SUCCESS"})
	
	total, completed, _ := s.repo.CountMilestones(ctx, id)
	response := dto.ToGoalResponse(goal, total, completed)
	
	// Check if goal was completed
	if req.IsCompleted != nil && *req.IsCompleted && goal.IsCompleted {
		if s.broadcaster != nil {
			s.broadcaster.Publish(userID, notification.EventGoalCompleted, map[string]interface{}{
				"goal_id": id,
				"goal":    response,
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventGoalCompleted,
				"user_id":    userID,
				"entity_id":  id,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}
	}
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventGoalUpdated, map[string]interface{}{
			"goal_id": id,
			"goal":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventGoalUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return response, nil
}

func (s *goalService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting goal", map[string]interface{}{"user_id": userID, "goal_id": id, "action": "DELETE_GOAL"})
	goal, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if goal == nil {
		return errors.New("goal not found")
	}
	if goal.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete goal", err, map[string]interface{}{"user_id": userID, "goal_id": id, "action": "DELETE_GOAL_FAILED"})
		return err
	}
	s.logger.Info("Goal deleted", map[string]interface{}{"user_id": userID, "goal_id": id, "action": "DELETE_GOAL_SUCCESS"})
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventGoalDeleted, map[string]interface{}{
			"goal_id": id,
			"title":   goal.Title,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventGoalDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return nil
}
