package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/domain"
)

type CourseRepository interface {
	Create(ctx context.Context, course *domain.Course) (*domain.Course, error)
	GetByID(ctx context.Context, id int) (*domain.Course, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Course, error)
	GetActiveCourses(ctx context.Context, userID int) ([]*domain.Course, error)
	Update(ctx context.Context, course *domain.Course) error
	Delete(ctx context.Context, id int) error

	CreateComponent(ctx context.Context, comp *domain.Component) (*domain.Component, error)
	UpdateComponent(ctx context.Context, comp *domain.Component) error
	DeleteComponent(ctx context.Context, id int) error
	GetComponents(ctx context.Context, courseID int) ([]*domain.Component, error)
	GetComponentByID(ctx context.Context, id int) (*domain.Component, error)

	CreateSchedule(ctx context.Context, sched *domain.Schedule) (*domain.Schedule, error)
	UpdateSchedule(ctx context.Context, sched *domain.Schedule) error
	DeleteSchedule(ctx context.Context, id int) error
	GetSchedules(ctx context.Context, courseID int) ([]*domain.Schedule, error)
	GetScheduleByID(ctx context.Context, id int) (*domain.Schedule, error)
}
