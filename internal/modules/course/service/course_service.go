package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/dto"
)

type CourseService interface {
	Create(ctx context.Context, req *dto.CreateCourseRequest, userID int) (*dto.CourseResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.CourseResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.CourseResponse, error)
	GetActive(ctx context.Context, userID int) ([]*dto.CourseResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateCourseRequest, userID int) (*dto.CourseResponse, error)
	Delete(ctx context.Context, id, userID int) error

	// Component methods
	CreateComponent(ctx context.Context, req *dto.CreateComponentRequest, userID int) (*dto.ComponentResponse, error)
	GetComponents(ctx context.Context, courseID, userID int) ([]*dto.ComponentResponse, error)
	UpdateComponent(ctx context.Context, id int, req *dto.UpdateComponentRequest, userID int) (*dto.ComponentResponse, error)
	DeleteComponent(ctx context.Context, id, userID int) error

	// Schedule methods
	CreateSchedule(ctx context.Context, req *dto.CreateScheduleRequest, userID int) (*dto.ScheduleResponse, error)
	GetSchedules(ctx context.Context, courseID, userID int) ([]*dto.ScheduleResponse, error)
	UpdateSchedule(ctx context.Context, id int, req *dto.UpdateScheduleRequest, userID int) (*dto.ScheduleResponse, error)
	DeleteSchedule(ctx context.Context, id, userID int) error
}
