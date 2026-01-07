package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
)

type habitService struct {
	repo        repository.HabitRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewHabitService(repo repository.HabitRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) HabitService {
	return &habitService{repo: repo, logger: logger, broadcaster: broadcaster}
}

func (s *habitService) Create(ctx context.Context, req *dto.CreateHabitRequest, userID int) (*dto.HabitResponse, error) {
	s.logger.Info("Creating habit", map[string]interface{}{"user_id": userID, "name": req.Name, "action": "CREATE_HABIT"})
	now := time.Now()
	targetCount := req.TargetCount
	if req.TargetDays > 0 {
		targetCount = req.TargetDays
	}

	config := make(map[string]interface{})
	if req.FrequencyDays != nil {
		config["days"] = req.FrequencyDays
	}
	if req.IntervalDays != nil {
		config["interval"] = *req.IntervalDays
	}

	habit := &domain.Habit{
		UserID:          userID,
		LifeAreaID:      req.LifeAreaID,
		Name:            req.Name,
		Icon:            req.Icon,
		Description:     req.Description,
		Frequency:       req.Frequency,
		FrequencyConfig: config,
		TargetCount:     targetCount,
		TimeOfDay:       req.TimeOfDay,
		ReminderTime:    req.ReminderTime,
		IsActive:        true,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	created, err := s.repo.Create(ctx, habit)
	if err != nil {
		s.logger.Error("Failed to create habit", err, map[string]interface{}{"user_id": userID, "action": "CREATE_HABIT_FAILED"})
		return nil, err
	}
	s.logger.Info("Habit created", map[string]interface{}{"user_id": userID, "habit_id": created.ID, "action": "CREATE_HABIT_SUCCESS"})
	
	response := dto.ToHabitResponse(created, false, false)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventHabitCreated, map[string]interface{}{
			"habit_id": created.ID,
			"habit":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventHabitCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return response, nil
}

func (s *habitService) GetByID(ctx context.Context, id, userID int) (*dto.HabitResponse, error) {
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}
	if habit.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	completedToday, _ := s.repo.HasLogForToday(ctx, id)
	skippedToday, _ := s.repo.HasSkippedToday(ctx, id)
	return dto.ToHabitResponse(habit, completedToday, skippedToday), nil
}

func (s *habitService) GetAll(ctx context.Context, userID int) ([]*dto.HabitResponse, error) {
	habits, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.HabitResponse, len(habits))
	for i, h := range habits {
		completedToday, _ := s.repo.HasLogForToday(ctx, h.ID)
		skippedToday, _ := s.repo.HasSkippedToday(ctx, h.ID)
		result[i] = dto.ToHabitResponse(h, completedToday, skippedToday)
	}
	return result, nil
}

func (s *habitService) GetActive(ctx context.Context, userID int) ([]*dto.HabitResponse, error) {
	habits, err := s.repo.GetActiveHabits(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*dto.HabitResponse, len(habits))
	for i, h := range habits {
		completedToday, _ := s.repo.HasLogForToday(ctx, h.ID)
		skippedToday, _ := s.repo.HasSkippedToday(ctx, h.ID)
		result[i] = dto.ToHabitResponse(h, completedToday, skippedToday)
	}
	return result, nil
}

func (s *habitService) Update(ctx context.Context, id int, req *dto.UpdateHabitRequest, userID int) (*dto.HabitResponse, error) {
	s.logger.Info("Updating habit", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "UPDATE_HABIT"})
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, errors.New("habit not found")
	}
	if habit.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Name != nil {
		habit.Name = *req.Name
	}
	if req.Description != nil {
		habit.Description = *req.Description
	}
	if req.Frequency != nil {
		habit.Frequency = *req.Frequency
	}
	if req.TargetCount != nil {
		habit.TargetCount = *req.TargetCount
	}
	if req.IsActive != nil {
		habit.IsActive = *req.IsActive
	}
	if req.LifeAreaID != nil {
		habit.LifeAreaID = req.LifeAreaID
	}
	if req.Icon != nil {
		habit.Icon = *req.Icon
	}
	if req.TimeOfDay != nil {
		habit.TimeOfDay = *req.TimeOfDay
	}
	if req.ReminderTime != nil {
		habit.ReminderTime = *req.ReminderTime
	}
	if req.TargetDays != nil {
		habit.TargetCount = *req.TargetDays
	}

	if req.FrequencyDays != nil || req.IntervalDays != nil {
		if habit.FrequencyConfig == nil {
			habit.FrequencyConfig = make(map[string]interface{})
		}
		if req.FrequencyDays != nil {
			habit.FrequencyConfig["days"] = req.FrequencyDays
		}
		if req.IntervalDays != nil {
			habit.FrequencyConfig["interval"] = *req.IntervalDays
		}
	}
	habit.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, habit); err != nil {
		s.logger.Error("Failed to update habit", err, map[string]interface{}{"user_id": userID, "habit_id": id, "action": "UPDATE_HABIT_FAILED"})
		return nil, err
	}
	s.logger.Info("Habit updated", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "UPDATE_HABIT_SUCCESS"})
	
	completedToday, _ := s.repo.HasLogForToday(ctx, id)
	skippedToday, _ := s.repo.HasSkippedToday(ctx, id)
	response := dto.ToHabitResponse(habit, completedToday, skippedToday)
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventHabitUpdated, map[string]interface{}{
			"habit_id": id,
			"habit":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventHabitUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return response, nil
}

func (s *habitService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting habit", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "DELETE_HABIT"})
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}
	if habit.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete habit", err, map[string]interface{}{"user_id": userID, "habit_id": id, "action": "DELETE_HABIT_FAILED"})
		return err
	}
	s.logger.Info("Habit deleted", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "DELETE_HABIT_SUCCESS"})
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventHabitDeleted, map[string]interface{}{
			"habit_id": id,
			"title":    habit.Name,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventHabitDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return nil
}

func (s *habitService) LogHabit(ctx context.Context, id int, req *dto.LogHabitRequest, userID int) error {
	s.logger.Info("Logging habit", map[string]interface{}{"user_id": userID, "habit_id": id, "count": req.Count, "action": "LOG_HABIT"})
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}
	if habit.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Check if habit is already completed or skipped today - prevent multiple actions in the same day
	alreadyCompleted, err := s.repo.HasLogForToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadyCompleted {
		return errors.New("habit already completed today")
	}
	
	alreadySkipped, err := s.repo.HasSkippedToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadySkipped {
		return errors.New("habit already skipped today")
	}
	
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.repo.LogHabit(ctx, id, today, req.Count, req.Notes); err != nil {
		s.logger.Error("Failed to log habit", err, map[string]interface{}{"user_id": userID, "habit_id": id, "action": "LOG_HABIT_FAILED"})
		return err
	}
	
	// Check if habit was completed (count >= target)
	wasCompleted := req.Count >= habit.TargetCount
	oldStreak := habit.CurrentStreak
	
	if wasCompleted {
		habit.IncrementStreak()
		if err := s.repo.Update(ctx, habit); err != nil {
			return err
		}
		s.logger.Info("Habit streak incremented", map[string]interface{}{"habit_id": id, "streak": habit.CurrentStreak, "action": "STREAK_INCREMENT"})
	}

	// Always broadcast habit completion event
	if s.broadcaster != nil {
		completedToday, _ := s.repo.HasLogForToday(ctx, id)
		skippedToday, _ := s.repo.HasSkippedToday(ctx, id)
		habitResponse := dto.ToHabitResponse(habit, completedToday, skippedToday)
		
		s.broadcaster.Publish(userID, notification.EventHabitCompleted, map[string]interface{}{
			"habit_id": id,
			"habit":    habitResponse,
			"streak":   habit.CurrentStreak,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventHabitCompleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})

		// Broadcast streak increased event if streak changed
		if wasCompleted && habit.CurrentStreak > oldStreak {
			s.broadcaster.Publish(userID, notification.EventStreakIncreased, map[string]interface{}{
				"habit_id": id,
				"habit":    habitResponse,
				"streak":   habit.CurrentStreak,
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventStreakIncreased,
				"user_id":    userID,
				"entity_id":  id,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}

		// Broadcast milestone event if reached milestone
		if wasCompleted && habit.CurrentStreak%10 == 0 && habit.CurrentStreak > 0 {
			s.broadcaster.Publish(userID, notification.EventStreakMilestone, map[string]interface{}{
				"habit_id":  id,
				"habit":     habitResponse,
				"streak":    habit.CurrentStreak,
				"milestone": notification.GetMilestoneName(habit.CurrentStreak),
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventStreakMilestone,
				"user_id":    userID,
				"entity_id":  id,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}
	}
	
	s.logger.Info("Habit logged", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "LOG_HABIT_SUCCESS"})
	return nil
}

func (s *habitService) Complete(ctx context.Context, id int, req *dto.LogHabitRequest, userID int) error {
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}
	if habit.UserID != userID {
		return errors.New("unauthorized")
	}

	// Check if habit is already completed or skipped today - prevent multiple actions in the same day
	alreadyCompleted, err := s.repo.HasLogForToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadyCompleted {
		return errors.New("habit already completed today")
	}
	
	alreadySkipped, err := s.repo.HasSkippedToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadySkipped {
		return errors.New("habit already skipped today - cannot complete")
	}

	// If count is not provided or less than target, set it to target
	count := req.Count
	if count < habit.TargetCount {
		count = habit.TargetCount
	}

	newReq := *req
	newReq.Count = count
	return s.LogHabit(ctx, id, &newReq, userID)
}

func (s *habitService) SkipHabit(ctx context.Context, id int, userID int) error {
	s.logger.Info("Skipping habit", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "SKIP_HABIT"})
	habit, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if habit == nil {
		return errors.New("habit not found")
	}
	if habit.UserID != userID {
		return errors.New("unauthorized")
	}
	
	// Check if habit is already completed or skipped today - prevent multiple actions in the same day
	alreadyCompleted, err := s.repo.HasLogForToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadyCompleted {
		return errors.New("habit already completed today - cannot skip")
	}
	
	alreadySkipped, err := s.repo.HasSkippedToday(ctx, id)
	if err != nil {
		return err
	}
	if alreadySkipped {
		return errors.New("habit already skipped today")
	}
	
	today := time.Now().Truncate(24 * time.Hour)
	if err := s.repo.SkipHabit(ctx, id, today, ""); err != nil {
		s.logger.Error("Failed to skip habit", err, map[string]interface{}{"user_id": userID, "habit_id": id, "action": "SKIP_HABIT_FAILED"})
		return err
	}
	
	s.logger.Info("Habit skipped", map[string]interface{}{"user_id": userID, "habit_id": id, "action": "SKIP_HABIT_SUCCESS"})
	
	// Broadcast skip event
	if s.broadcaster != nil {
		completedToday, _ := s.repo.HasLogForToday(ctx, id)
		skippedToday, _ := s.repo.HasSkippedToday(ctx, id)
		habitResponse := dto.ToHabitResponse(habit, completedToday, skippedToday)
		
		s.broadcaster.Publish(userID, notification.EventHabitSkipped, map[string]interface{}{
			"habit_id": id,
			"habit":    habitResponse,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventHabitSkipped,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	
	return nil
}
