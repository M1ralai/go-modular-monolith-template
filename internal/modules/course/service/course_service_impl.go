package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/repository"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/notification"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
)

// normalizeTime ensures time string is in HH:MM format
// Handles various input formats and converts to PostgreSQL TIME format
func normalizeTime(timeStr string) string {
	if timeStr == "" {
		return "00:00"
	}

	// Remove whitespace
	timeStr = strings.TrimSpace(timeStr)

	// Split by colon
	parts := strings.Split(timeStr, ":")
	if len(parts) < 2 {
		return "00:00"
	}

	// Parse hour and minute
	var hour, minute int
	fmt.Sscanf(parts[0], "%d", &hour)
	fmt.Sscanf(parts[1], "%d", &minute)

	// Validate ranges
	if hour < 0 || hour > 23 {
		hour = 0
	}
	if minute < 0 || minute > 59 {
		minute = 0
	}

	// Format as HH:MM
	return fmt.Sprintf("%02d:%02d", hour, minute)
}

type courseService struct {
	repo        repository.CourseRepository
	logger      *logger.ZapLogger
	broadcaster *notifService.Broadcaster
}

func NewCourseService(repo repository.CourseRepository, logger *logger.ZapLogger, broadcaster *notifService.Broadcaster) CourseService {
	return &courseService{
		repo:        repo,
		logger:      logger,
		broadcaster: broadcaster,
	}
}

func (s *courseService) Create(ctx context.Context, req *dto.CreateCourseRequest, userID int) (*dto.CourseResponse, error) {
	s.logger.Info("Creating course", map[string]interface{}{
		"user_id": userID,
		"name":    req.Name,
		"code":    req.Code,
		"action":  "CREATE_COURSE",
	})

	now := time.Now()
	course := &domain.Course{
		UserID:      userID,
		Name:        req.Name,
		Code:        req.Code,
		Instructor:  req.Instructor,
		Credits:     req.Credits,
		Semester:    req.Semester,
		Type:        req.Type,
		Color:       req.Color,
		SyllabusURL: req.SyllabusURL,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	created, err := s.repo.Create(ctx, course)
	if err != nil {
		s.logger.Error("failed to create course", err, map[string]interface{}{
			"user_id": userID,
			"name":    req.Name,
			"action":  "CREATE_COURSE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("course created", map[string]interface{}{
		"course_id": created.ID,
		"user_id":   userID,
		"action":    "CREATE_COURSE",
	})

	response := dto.ToCourseResponse(created)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventCourseCreated, map[string]interface{}{
			"course_id": created.ID,
			"course":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventCourseCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) GetByID(ctx context.Context, id, userID int) (*dto.CourseResponse, error) {
	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Load components and schedules
	components, err := s.repo.GetComponents(ctx, id)
	if err != nil {
		s.logger.Error("failed to load components", err, map[string]interface{}{
			"course_id": id,
		})
		components = []*domain.Component{}
	}
	course.Components = components

	schedules, err := s.repo.GetSchedules(ctx, id)
	if err != nil {
		s.logger.Error("failed to load schedules", err, map[string]interface{}{
			"course_id": id,
		})
		schedules = []*domain.Schedule{}
	}
	course.Schedules = schedules

	return dto.ToCourseResponse(course), nil
}

func (s *courseService) GetAll(ctx context.Context, userID int) ([]*dto.CourseResponse, error) {
	courses, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load components and schedules for all courses
	for _, course := range courses {
		components, err := s.repo.GetComponents(ctx, course.ID)
		if err != nil {
			s.logger.Error("failed to load components", err, map[string]interface{}{
				"course_id": course.ID,
			})
			components = []*domain.Component{}
		}
		course.Components = components

		schedules, err := s.repo.GetSchedules(ctx, course.ID)
		if err != nil {
			s.logger.Error("failed to load schedules", err, map[string]interface{}{
				"course_id": course.ID,
			})
			schedules = []*domain.Schedule{}
		}
		course.Schedules = schedules
	}

	return dto.ToCourseResponseList(courses), nil
}

func (s *courseService) GetActive(ctx context.Context, userID int) ([]*dto.CourseResponse, error) {
	courses, err := s.repo.GetActiveCourses(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load components and schedules for all active courses
	for _, course := range courses {
		components, err := s.repo.GetComponents(ctx, course.ID)
		if err != nil {
			s.logger.Error("failed to load components", err, map[string]interface{}{
				"course_id": course.ID,
			})
			components = []*domain.Component{}
		}
		course.Components = components

		schedules, err := s.repo.GetSchedules(ctx, course.ID)
		if err != nil {
			s.logger.Error("failed to load schedules", err, map[string]interface{}{
				"course_id": course.ID,
			})
			schedules = []*domain.Schedule{}
		}
		course.Schedules = schedules
	}

	return dto.ToCourseResponseList(courses), nil
}

func (s *courseService) Update(ctx context.Context, id int, req *dto.UpdateCourseRequest, userID int) (*dto.CourseResponse, error) {
	s.logger.Info("Updating course", map[string]interface{}{
		"user_id":   userID,
		"course_id": id,
		"action":    "UPDATE_COURSE",
	})

	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	if req.Name != nil {
		course.Name = *req.Name
	}
	if req.Code != nil {
		course.Code = *req.Code
	}
	if req.Instructor != nil {
		course.Instructor = *req.Instructor
	}
	if req.Credits != nil {
		course.Credits = *req.Credits
	}
	if req.Semester != nil {
		course.Semester = *req.Semester
	}
	if req.Type != nil {
		course.Type = *req.Type
	}
	if req.Color != nil {
		course.Color = *req.Color
	}
	if req.SyllabusURL != nil {
		course.SyllabusURL = *req.SyllabusURL
	}
	if req.FinalGrade != nil {
		course.FinalGrade = *req.FinalGrade
	}
	if req.IsActive != nil {
		course.IsActive = *req.IsActive
	}

	course.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, course); err != nil {
		s.logger.Error("failed to update course", err, map[string]interface{}{
			"course_id": id,
			"user_id":   userID,
			"action":    "UPDATE_COURSE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("course updated", map[string]interface{}{
		"course_id": id,
		"user_id":   userID,
		"action":    "UPDATE_COURSE",
	})

	response := dto.ToCourseResponse(course)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventCourseUpdated, map[string]interface{}{
			"course_id": id,
			"course":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventCourseUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting course", map[string]interface{}{
		"user_id":   userID,
		"course_id": id,
		"action":    "DELETE_COURSE",
	})

	course, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}
	if course.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete course", err, map[string]interface{}{
			"course_id": id,
			"user_id":   userID,
			"action":    "DELETE_COURSE_FAILED",
		})
		return err
	}

	s.logger.Info("course deleted", map[string]interface{}{
		"course_id": id,
		"user_id":   userID,
		"action":    "DELETE_COURSE",
	})

	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventCourseDeleted, map[string]interface{}{
			"course_id": id,
			"name":      course.Name,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventCourseDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return nil
}

func (s *courseService) CreateComponent(ctx context.Context, req *dto.CreateComponentRequest, userID int) (*dto.ComponentResponse, error) {
	s.logger.Info("Creating course component", map[string]interface{}{
		"user_id":   userID,
		"course_id": req.CourseID,
		"name":      req.Name,
		"type":      req.Type,
		"action":    "CREATE_COMPONENT",
	})

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	now := time.Now()
	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		parsed, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			return nil, errors.New("invalid due_date format, use YYYY-MM-DD")
		}
		dueDate = &parsed
	}

	var weight float64
	if req.Weight != nil {
		weight = *req.Weight
	}

	var maxScore float64
	if req.MaxScore != nil {
		maxScore = *req.MaxScore
	}

	component := &domain.Component{
		CourseID:      req.CourseID,
		Type:          req.Type,
		Name:          req.Name,
		Weight:        weight,
		MaxScore:      maxScore,
		AchievedScore: req.AchievedScore,
		DueDate:       dueDate,
		IsCompleted:   false,
		Notes:         req.Notes,
		DisplayOrder:  req.DisplayOrder,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	created, err := s.repo.CreateComponent(ctx, component)
	if err != nil {
		s.logger.Error("failed to create component", err, map[string]interface{}{
			"user_id":   userID,
			"course_id": req.CourseID,
			"action":    "CREATE_COMPONENT_FAILED",
		})
		return nil, err
	}

	s.logger.Info("component created", map[string]interface{}{
		"component_id": created.ID,
		"course_id":    req.CourseID,
		"user_id":      userID,
		"action":       "CREATE_COMPONENT",
	})

	response := dto.ToComponentResponse(created)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventComponentCreated, map[string]interface{}{
			"component_id": created.ID,
			"course_id":    req.CourseID,
			"component":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventComponentCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) GetComponents(ctx context.Context, courseID, userID int) ([]*dto.ComponentResponse, error) {
	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	components, err := s.repo.GetComponents(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return dto.ToComponentResponseList(components), nil
}

func (s *courseService) UpdateComponent(ctx context.Context, id int, req *dto.UpdateComponentRequest, userID int) (*dto.ComponentResponse, error) {
	s.logger.Info("Updating course component", map[string]interface{}{
		"user_id":      userID,
		"component_id": id,
		"action":       "UPDATE_COMPONENT",
	})

	component, err := s.repo.GetComponentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if component == nil {
		return nil, errors.New("component not found")
	}

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, component.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	if req.Type != nil {
		component.Type = *req.Type
	}
	if req.Name != nil {
		component.Name = *req.Name
	}
	if req.Weight != nil {
		component.Weight = *req.Weight
	}
	if req.MaxScore != nil {
		component.MaxScore = *req.MaxScore
	}
	if req.AchievedScore != nil {
		component.AchievedScore = req.AchievedScore
	}
	if req.DueDate != nil {
		if *req.DueDate == "" {
			component.DueDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.DueDate)
			if err != nil {
				return nil, errors.New("invalid due_date format, use YYYY-MM-DD")
			}
			component.DueDate = &parsed
		}
	}
	if req.CompletionDate != nil {
		if *req.CompletionDate == "" {
			component.CompletionDate = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.CompletionDate)
			if err != nil {
				return nil, errors.New("invalid completion_date format, use YYYY-MM-DD")
			}
			component.CompletionDate = &parsed
		}
	}
	if req.IsCompleted != nil {
		component.IsCompleted = *req.IsCompleted
	}
	if req.Notes != nil {
		component.Notes = *req.Notes
	}
	if req.DisplayOrder != nil {
		component.DisplayOrder = *req.DisplayOrder
	}

	component.UpdatedAt = time.Now()

	if err := s.repo.UpdateComponent(ctx, component); err != nil {
		s.logger.Error("failed to update component", err, map[string]interface{}{
			"component_id": id,
			"user_id":      userID,
			"action":       "UPDATE_COMPONENT_FAILED",
		})
		return nil, err
	}

	s.logger.Info("component updated", map[string]interface{}{
		"component_id": id,
		"user_id":      userID,
		"action":       "UPDATE_COMPONENT",
	})

	response := dto.ToComponentResponse(component)
	
	// Check if grade was updated - broadcast grade change
	if req.AchievedScore != nil && s.broadcaster != nil {
		// Recalculate course grade
		course, _ := s.repo.GetByID(ctx, component.CourseID)
		if course != nil {
			components, _ := s.repo.GetComponents(ctx, component.CourseID)
			var totalWeight, weightedScore float64
			for _, comp := range components {
				if comp.AchievedScore != nil && comp.IsCompleted {
					totalWeight += comp.Weight
					weightedScore += (*comp.AchievedScore / comp.MaxScore) * comp.Weight
				}
			}
			var newGrade float64
			if totalWeight > 0 {
				newGrade = (weightedScore / totalWeight) * 100
			}
			
			s.broadcaster.Publish(userID, notification.EventComponentGraded, map[string]interface{}{
				"component_id": id,
				"course_id":    component.CourseID,
				"component":     response,
				"new_grade":    newGrade,
			})
			s.logger.Info("WebSocket event published", map[string]interface{}{
				"event_type": notification.EventComponentGraded,
				"user_id":    userID,
				"entity_id":  id,
				"action":     "WS_EVENT_PUBLISHED",
			})
		}
	}
	
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventComponentUpdated, map[string]interface{}{
			"component_id": id,
			"course_id":   component.CourseID,
			"component":   response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventComponentUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) DeleteComponent(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting course component", map[string]interface{}{
		"user_id":      userID,
		"component_id": id,
		"action":       "DELETE_COMPONENT",
	})

	component, err := s.repo.GetComponentByID(ctx, id)
	if err != nil {
		return err
	}
	if component == nil {
		return errors.New("component not found")
	}

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, component.CourseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}
	if course.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.DeleteComponent(ctx, id); err != nil {
		s.logger.Error("failed to delete component", err, map[string]interface{}{
			"component_id": id,
			"user_id":      userID,
			"action":       "DELETE_COMPONENT_FAILED",
		})
		return err
	}

	s.logger.Info("component deleted", map[string]interface{}{
		"component_id": id,
		"user_id":      userID,
		"action":       "DELETE_COMPONENT",
	})

	return nil
}

func (s *courseService) CreateSchedule(ctx context.Context, req *dto.CreateScheduleRequest, userID int) (*dto.ScheduleResponse, error) {
	s.logger.Info("Creating course schedule", map[string]interface{}{
		"user_id":   userID,
		"course_id": req.CourseID,
		"day":       req.DayOfWeek,
		"action":    "CREATE_SCHEDULE",
	})

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, req.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Validate and normalize time format (HH:MM)
	startTime := normalizeTime(req.StartTime)
	endTime := normalizeTime(req.EndTime)

	schedule := &domain.Schedule{
		CourseID:  req.CourseID,
		DayOfWeek: req.DayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
		Location:  req.Location,
		CreatedAt: time.Now(),
	}

	created, err := s.repo.CreateSchedule(ctx, schedule)
	if err != nil {
		s.logger.Error("failed to create schedule", err, map[string]interface{}{
			"user_id":   userID,
			"course_id": req.CourseID,
			"action":    "CREATE_SCHEDULE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("schedule created", map[string]interface{}{
		"schedule_id": created.ID,
		"course_id":   req.CourseID,
		"user_id":     userID,
		"action":      "CREATE_SCHEDULE",
	})

	response := dto.ToScheduleResponse(created)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventScheduleCreated, map[string]interface{}{
			"schedule_id": created.ID,
			"course_id":   req.CourseID,
			"schedule":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventScheduleCreated,
			"user_id":    userID,
			"entity_id":  created.ID,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) GetSchedules(ctx context.Context, courseID, userID int) ([]*dto.ScheduleResponse, error) {
	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	schedules, err := s.repo.GetSchedules(ctx, courseID)
	if err != nil {
		return nil, err
	}

	return dto.ToScheduleResponseList(schedules), nil
}

func (s *courseService) UpdateSchedule(ctx context.Context, id int, req *dto.UpdateScheduleRequest, userID int) (*dto.ScheduleResponse, error) {
	s.logger.Info("Updating course schedule", map[string]interface{}{
		"user_id":     userID,
		"schedule_id": id,
		"action":      "UPDATE_SCHEDULE",
	})

	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if schedule == nil {
		return nil, errors.New("schedule not found")
	}

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, schedule.CourseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("course not found")
	}
	if course.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	// Update fields
	if req.DayOfWeek != nil {
		schedule.DayOfWeek = *req.DayOfWeek
	}
	if req.StartTime != nil {
		normalizedStart := normalizeTime(*req.StartTime)
		s.logger.Info("Normalizing start time", map[string]interface{}{
			"original":  *req.StartTime,
			"normalized": normalizedStart,
		})
		schedule.StartTime = normalizedStart
	}
	if req.EndTime != nil {
		normalizedEnd := normalizeTime(*req.EndTime)
		s.logger.Info("Normalizing end time", map[string]interface{}{
			"original":  *req.EndTime,
			"normalized": normalizedEnd,
		})
		schedule.EndTime = normalizedEnd
	}
	if req.Location != nil {
		schedule.Location = *req.Location
	}

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		s.logger.Error("failed to update schedule", err, map[string]interface{}{
			"schedule_id": id,
			"user_id":     userID,
			"action":      "UPDATE_SCHEDULE_FAILED",
		})
		return nil, err
	}

	s.logger.Info("schedule updated", map[string]interface{}{
		"schedule_id": id,
		"user_id":     userID,
		"action":      "UPDATE_SCHEDULE",
	})

	response := dto.ToScheduleResponse(schedule)
	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventScheduleUpdated, map[string]interface{}{
			"schedule_id": id,
			"course_id":   schedule.CourseID,
			"schedule":    response,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventScheduleUpdated,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return response, nil
}

func (s *courseService) DeleteSchedule(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting course schedule", map[string]interface{}{
		"user_id":     userID,
		"schedule_id": id,
		"action":      "DELETE_SCHEDULE",
	})

	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return err
	}
	if schedule == nil {
		return errors.New("schedule not found")
	}

	// Verify course belongs to user
	course, err := s.repo.GetByID(ctx, schedule.CourseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("course not found")
	}
	if course.UserID != userID {
		return errors.New("unauthorized")
	}

	if err := s.repo.DeleteSchedule(ctx, id); err != nil {
		s.logger.Error("failed to delete schedule", err, map[string]interface{}{
			"schedule_id": id,
			"user_id":     userID,
			"action":      "DELETE_SCHEDULE_FAILED",
		})
		return err
	}

	s.logger.Info("schedule deleted", map[string]interface{}{
		"schedule_id": id,
		"user_id":     userID,
		"action":      "DELETE_SCHEDULE",
	})

	if s.broadcaster != nil {
		s.broadcaster.Publish(userID, notification.EventScheduleDeleted, map[string]interface{}{
			"schedule_id": id,
			"course_id":   schedule.CourseID,
			"day_of_week": schedule.DayOfWeek,
		})
		s.logger.Info("WebSocket event published", map[string]interface{}{
			"event_type": notification.EventScheduleDeleted,
			"user_id":    userID,
			"entity_id":  id,
			"action":     "WS_EVENT_PUBLISHED",
		})
	}

	return nil
}
