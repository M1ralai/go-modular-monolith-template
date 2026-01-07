package service

import (
	"context"
	"sort"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/repository"
)

type scheduleService struct {
	blockedSlotRepo repository.BlockedTimeSlotRepository
	logger          *logger.ZapLogger
	broadcaster     *notifService.Broadcaster
}

func NewScheduleService(blockedSlotRepo repository.BlockedTimeSlotRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) ScheduleService {
	return &scheduleService{blockedSlotRepo: blockedSlotRepo, logger: logger, broadcaster: broadcaster}
}

// CheckConflict implements the conflict detection algorithm
func (s *scheduleService) CheckConflict(ctx context.Context, userID int, start, end time.Time) (*dto.ConflictResponse, error) {
	s.logger.Info("Checking time conflict", map[string]interface{}{"user_id": userID, "start": start, "end": end, "action": "CHECK_CONFLICT"})

	blockedSlots, err := s.blockedSlotRepo.GetByUserAndTimeRange(ctx, userID, start, end)
	if err != nil {
		s.logger.Error("Failed to check conflict", err, map[string]interface{}{"user_id": userID, "action": "CHECK_CONFLICT_FAILED"})
		return nil, err
	}

	for _, slot := range blockedSlots {
		if slot.OverlapsWith(start, end) && !slot.IsFlexible {
			s.logger.Info("Time conflict detected", map[string]interface{}{"user_id": userID, "reason": slot.Reason, "action": "CONFLICT_DETECTED"})

			suggestions, _ := s.GetFreeSlots(ctx, userID, start, int(end.Sub(start).Minutes()))
			suggestionDTOs := make([]*dto.TimeSlotResponse, len(suggestions))
			copy(suggestionDTOs, suggestions)

			if s.broadcaster != nil {
				s.broadcaster.Publish(userID, notification.EventConflictDetected, map[string]interface{}{
					"reason": slot.Reason,
					"start":  start,
					"end":    end,
				})
			}

			return &dto.ConflictResponse{HasConflict: true, Reason: slot.Reason, Suggestions: suggestionDTOs}, nil
		}
	}

	s.logger.Info("No conflict found", map[string]interface{}{"user_id": userID, "action": "CHECK_CONFLICT_SUCCESS"})
	return &dto.ConflictResponse{HasConflict: false}, nil
}

// GetFreeSlots implements the free slot finder algorithm
func (s *scheduleService) GetFreeSlots(ctx context.Context, userID int, date time.Time, durationMinutes int) ([]*dto.TimeSlotResponse, error) {
	s.logger.Info("Finding free slots", map[string]interface{}{"user_id": userID, "date": date, "duration": durationMinutes, "action": "GET_FREE_SLOTS"})

	blockedSlots, err := s.blockedSlotRepo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	// Sort by start time
	sort.Slice(blockedSlots, func(i, j int) bool { return blockedSlots[i].StartDatetime.Before(blockedSlots[j].StartDatetime) })

	// Day boundaries (08:00 - 22:00)
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 8, 0, 0, 0, date.Location())
	dayEnd := time.Date(date.Year(), date.Month(), date.Day(), 22, 0, 0, 0, date.Location())

	var freeSlots []*dto.TimeSlotResponse
	currentStart := dayStart

	for _, slot := range blockedSlots {
		if slot.StartDatetime.After(currentStart) {
			gap := slot.StartDatetime.Sub(currentStart)
			if int(gap.Minutes()) >= durationMinutes {
				freeSlots = append(freeSlots, &dto.TimeSlotResponse{Start: currentStart, End: slot.StartDatetime, DurationMinutes: int(gap.Minutes())})
			}
		}
		if slot.EndDatetime.After(currentStart) {
			currentStart = slot.EndDatetime
		}
	}

	// Check remaining time until day end
	if dayEnd.After(currentStart) {
		gap := dayEnd.Sub(currentStart)
		if int(gap.Minutes()) >= durationMinutes {
			freeSlots = append(freeSlots, &dto.TimeSlotResponse{Start: currentStart, End: dayEnd, DurationMinutes: int(gap.Minutes())})
		}
	}

	s.logger.Info("Free slots found", map[string]interface{}{"user_id": userID, "count": len(freeSlots), "action": "GET_FREE_SLOTS_SUCCESS"})
	return freeSlots, nil
}

// GetBlockedSlots returns blocked slots for a date
func (s *scheduleService) GetBlockedSlots(ctx context.Context, userID int, date time.Time) ([]*dto.BlockedSlotResponse, error) {
	slots, err := s.blockedSlotRepo.GetByUserAndDate(ctx, userID, date)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.BlockedSlotResponse, len(slots))
	for i, slot := range slots {
		result[i] = dto.ToBlockedSlotResponse(slot)
	}
	return result, nil
}

// CreateBlockedSlot creates a new blocked time slot
func (s *scheduleService) CreateBlockedSlot(ctx context.Context, userID int, req *dto.CreateBlockedSlotRequest) (*dto.BlockedSlotResponse, error) {
	s.logger.Info("Creating blocked slot", map[string]interface{}{"user_id": userID, "reason": req.Reason, "action": "CREATE_BLOCKED_SLOT"})

	slot := &domain.BlockedTimeSlot{
		UserID:        userID,
		SourceType:    req.SourceType,
		SourceID:      req.SourceID,
		StartDatetime: req.StartDatetime,
		EndDatetime:   req.EndDatetime,
		Reason:        req.Reason,
		IsFlexible:    req.IsFlexible,
		CreatedAt:     time.Now(),
	}

	created, err := s.blockedSlotRepo.Create(ctx, slot)
	if err != nil {
		s.logger.Error("Failed to create blocked slot", err, map[string]interface{}{"user_id": userID, "action": "CREATE_BLOCKED_SLOT_FAILED"})
		return nil, err
	}

	s.logger.Info("Blocked slot created", map[string]interface{}{"user_id": userID, "slot_id": created.ID, "action": "CREATE_BLOCKED_SLOT_SUCCESS"})
	return dto.ToBlockedSlotResponse(created), nil
}

// DeleteBlockedSlot deletes a blocked slot
func (s *scheduleService) DeleteBlockedSlot(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting blocked slot", map[string]interface{}{"user_id": userID, "slot_id": id, "action": "DELETE_BLOCKED_SLOT"})

	slot, err := s.blockedSlotRepo.GetByID(ctx, id)
	if err != nil || slot == nil || slot.UserID != userID {
		return err
	}

	if err := s.blockedSlotRepo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete blocked slot", err, map[string]interface{}{"user_id": userID, "slot_id": id, "action": "DELETE_BLOCKED_SLOT_FAILED"})
		return err
	}

	s.logger.Info("Blocked slot deleted", map[string]interface{}{"user_id": userID, "slot_id": id, "action": "DELETE_BLOCKED_SLOT_SUCCESS"})
	return nil
}

// GenerateEventsForSchedule generates semester events from a course schedule
func (s *scheduleService) GenerateEventsForSchedule(ctx context.Context, userID int, req *dto.GenerateEventsRequest) (*dto.GenerateEventsResponse, error) {
	s.logger.Info("Generating semester events", map[string]interface{}{
		"user_id":     userID,
		"title":       req.Title,
		"day_of_week": req.DayOfWeek,
		"start_date":  req.SemesterStartDate,
		"end_date":    req.SemesterEndDate,
		"action":      "GENERATE_EVENTS",
	})

	// Parse times
	startTime, _ := time.Parse("15:04", req.StartTime)
	endTime, _ := time.Parse("15:04", req.EndTime)

	// Map day of week to time.Weekday
	dayMap := map[string]time.Weekday{"Sunday": 0, "Monday": 1, "Tuesday": 2, "Wednesday": 3, "Thursday": 4, "Friday": 5, "Saturday": 6}
	targetDay := dayMap[req.DayOfWeek]

	// Find first matching weekday
	current := req.SemesterStartDate
	for current.Weekday() != targetDay {
		current = current.AddDate(0, 0, 1)
	}

	eventsGenerated := 0
	blockedSlotsCreated := 0

	// Loop weekly until semester end
	for !current.After(req.SemesterEndDate) {
		// Skip excluded dates
		excluded := false
		for _, excl := range req.ExcludeDates {
			if current.Year() == excl.Year() && current.Month() == excl.Month() && current.Day() == excl.Day() {
				excluded = true
				break
			}
		}

		if !excluded {
			// Create blocked time slot for each occurrence
			eventStart := time.Date(current.Year(), current.Month(), current.Day(), startTime.Hour(), startTime.Minute(), 0, 0, current.Location())
			eventEnd := time.Date(current.Year(), current.Month(), current.Day(), endTime.Hour(), endTime.Minute(), 0, 0, current.Location())

			slot := &domain.BlockedTimeSlot{
				UserID:        userID,
				SourceType:    "course",
				SourceID:      req.CourseScheduleID,
				StartDatetime: eventStart,
				EndDatetime:   eventEnd,
				Reason:        req.Title + " (" + req.Location + ")",
				IsFlexible:    false,
			}

			if _, err := s.blockedSlotRepo.Create(ctx, slot); err == nil {
				blockedSlotsCreated++
				eventsGenerated++
			}
		}

		current = current.AddDate(0, 0, 7) // Next week
	}

	s.logger.Info("Semester events generated", map[string]interface{}{
		"user_id":               userID,
		"events_generated":      eventsGenerated,
		"blocked_slots_created": blockedSlotsCreated,
		"action":                "GENERATE_EVENTS_SUCCESS",
	})

	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventEventsGenerated, map[string]interface{}{
			"title":          req.Title,
			"events_created": eventsGenerated,
			"semester_start": req.SemesterStartDate,
			"semester_end":   req.SemesterEndDate,
		})
	}

	return &dto.GenerateEventsResponse{EventsGenerated: eventsGenerated, BlockedSlotsCreated: blockedSlotsCreated}, nil
}
