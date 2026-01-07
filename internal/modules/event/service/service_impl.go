package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
)

type eventService struct {
	repo        repository.EventRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewEventService(repo repository.EventRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) EventService {
	return &eventService{repo: repo, logger: logger, broadcaster: broadcaster}
}

func (s *eventService) Create(ctx context.Context, req *dto.CreateEventRequest, userID int) (*dto.EventResponse, error) {
	s.logger.Info("Creating event", map[string]interface{}{"user_id": userID, "title": req.Title, "action": "CREATE_EVENT"})
	now := time.Now()
	event := &domain.Event{UserID: userID, LifeAreaID: req.LifeAreaID, Title: req.Title, Description: req.Description, StartTime: req.StartTime, EndTime: req.EndTime, Location: req.Location, IsAllDay: req.IsAllDay, IsRecurring: req.IsRecurring, Recurrence: req.Recurrence, CreatedAt: now, UpdatedAt: now}
	created, err := s.repo.Create(ctx, event)
	if err != nil {
		s.logger.Error("Failed to create event", err, map[string]interface{}{"user_id": userID, "action": "CREATE_EVENT_FAILED"})
		return nil, err
	}
	s.logger.Info("Event created", map[string]interface{}{"user_id": userID, "event_id": created.ID, "action": "CREATE_EVENT_SUCCESS"})
	response := dto.ToEventResponse(created)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventEventCreated, map[string]interface{}{
			"event_id": created.ID,
			"event":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventEventCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	return response, nil
}

func (s *eventService) GetByID(ctx context.Context, id, userID int) (*dto.EventResponse, error) {
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}
	if event.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return dto.ToEventResponse(event), nil
}

func (s *eventService) GetAll(ctx context.Context, userID int) ([]*dto.EventResponse, error) {
	events, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToEventResponseList(events), nil
}

func (s *eventService) GetByDateRange(ctx context.Context, userID int, start, end time.Time) ([]*dto.EventResponse, error) {
	events, err := s.repo.GetByDateRange(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	return dto.ToEventResponseList(events), nil
}

func (s *eventService) Update(ctx context.Context, id int, req *dto.UpdateEventRequest, userID int) (*dto.EventResponse, error) {
	s.logger.Info("Updating event", map[string]interface{}{"user_id": userID, "event_id": id, "action": "UPDATE_EVENT"})
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event not found")
	}
	if event.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Title != nil {
		event.Title = *req.Title
	}
	if req.Description != nil {
		event.Description = *req.Description
	}
	if req.StartTime != nil {
		event.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		event.EndTime = req.EndTime
	}
	if req.Location != nil {
		event.Location = *req.Location
	}
	if req.IsAllDay != nil {
		event.IsAllDay = *req.IsAllDay
	}
	if req.LifeAreaID != nil {
		event.LifeAreaID = req.LifeAreaID
	}
	event.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, event); err != nil {
		s.logger.Error("Failed to update event", err, map[string]interface{}{"user_id": userID, "event_id": id, "action": "UPDATE_EVENT_FAILED"})
		return nil, err
	}
	s.logger.Info("Event updated", map[string]interface{}{"user_id": userID, "event_id": id, "action": "UPDATE_EVENT_SUCCESS"})
	response := dto.ToEventResponse(event)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventEventUpdated, map[string]interface{}{
			"event_id": id,
			"event":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventEventUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	return response, nil
}

func (s *eventService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting event", map[string]interface{}{"user_id": userID, "event_id": id, "action": "DELETE_EVENT"})
	event, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if event == nil {
		return errors.New("event not found")
	}
	if event.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete event", err, map[string]interface{}{"user_id": userID, "event_id": id, "action": "DELETE_EVENT_FAILED"})
		return err
	}
	s.logger.Info("Event deleted", map[string]interface{}{"user_id": userID, "event_id": id, "action": "DELETE_EVENT_SUCCESS"})
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventEventDeleted, map[string]interface{}{
			"event_id": id,
			"title":    event.Title,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventEventDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}
	return nil
}
