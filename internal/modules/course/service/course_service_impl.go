package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/repository"
)

type courseService struct {
	repo   repository.CourseRepository
	logger *logger.ZapLogger
}

func NewCourseService(repo repository.CourseRepository, logger *logger.ZapLogger) CourseService {
	return &courseService{
		repo:   repo,
		logger: logger,
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

	return dto.ToCourseResponse(created), nil
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

	return dto.ToCourseResponse(course), nil
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

	return dto.ToComponentResponse(created), nil
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

	return dto.ToComponentResponse(component), nil
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

	schedule := &domain.Schedule{
		CourseID:  req.CourseID,
		DayOfWeek: req.DayOfWeek,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
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

	return dto.ToScheduleResponse(created), nil
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
		schedule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		schedule.EndTime = *req.EndTime
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

	return dto.ToScheduleResponse(schedule), nil
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

	return nil
}
