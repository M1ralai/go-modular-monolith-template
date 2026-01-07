package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/course/domain"
)

type CourseModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Name        string    `db:"name"`
	Code        *string   `db:"code"`
	Instructor  *string   `db:"instructor"`
	Credits     *float64  `db:"credits"`
	Semester    *string   `db:"semester"`
	Type        *string   `db:"type"`
	Color       *string   `db:"color"`
	SyllabusURL *string   `db:"syllabus_url"`
	FinalGrade  *string   `db:"final_grade"`
	IsActive    bool      `db:"is_active"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *CourseModel) ToDomain() *domain.Course {
	return &domain.Course{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Code:        derefString(m.Code),
		Instructor:  derefString(m.Instructor),
		Credits:     derefFloat(m.Credits),
		Semester:    derefString(m.Semester),
		Type:        derefString(m.Type),
		Color:       derefString(m.Color),
		SyllabusURL: derefString(m.SyllabusURL),
		FinalGrade:  derefString(m.FinalGrade),
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromDomain(c *domain.Course) *CourseModel {
	return &CourseModel{
		ID:          c.ID,
		UserID:      c.UserID,
		Name:        c.Name,
		Code:        refString(c.Code),
		Instructor:  refString(c.Instructor),
		Credits:     refFloat(c.Credits),
		Semester:    refString(c.Semester),
		Type:        refString(c.Type),
		Color:       refString(c.Color),
		SyllabusURL: refString(c.SyllabusURL),
		FinalGrade:  refString(c.FinalGrade),
		IsActive:    c.IsActive,
		CreatedAt:   c.CreatedAt,
	}
}

type ComponentModel struct {
	ID             int        `db:"id"`
	CourseID       int        `db:"course_id"`
	Type           string     `db:"type"`
	Name           string     `db:"name"`
	Weight         float64    `db:"weight"`
	MaxScore       float64    `db:"max_score"`
	AchievedScore  *float64   `db:"achieved_score"`
	DueDate        *time.Time `db:"due_date"`
	CompletionDate *time.Time `db:"completion_date"`
	IsCompleted    bool       `db:"is_completed"`
	Notes          *string    `db:"notes"`
	DisplayOrder   int        `db:"display_order"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

func (m *ComponentModel) ToDomain() *domain.Component {
	return &domain.Component{
		ID:             m.ID,
		CourseID:       m.CourseID,
		Type:           m.Type,
		Name:           m.Name,
		Weight:         m.Weight,
		MaxScore:       m.MaxScore,
		AchievedScore:  m.AchievedScore,
		DueDate:        m.DueDate,
		CompletionDate: m.CompletionDate,
		IsCompleted:    m.IsCompleted,
		Notes:          derefString(m.Notes),
		DisplayOrder:   m.DisplayOrder,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}

func FromDomainComponent(c *domain.Component) *ComponentModel {
	return &ComponentModel{
		ID:             c.ID,
		CourseID:       c.CourseID,
		Type:           c.Type,
		Name:           c.Name,
		Weight:         c.Weight,
		MaxScore:       c.MaxScore,
		AchievedScore:  c.AchievedScore,
		DueDate:        c.DueDate,
		CompletionDate: c.CompletionDate,
		IsCompleted:    c.IsCompleted,
		Notes:          refString(c.Notes),
		DisplayOrder:   c.DisplayOrder,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func refString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func derefFloat(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func refFloat(f float64) *float64 {
	if f == 0 {
		return nil
	}
	return &f
}

type ScheduleModel struct {
	ID        int       `db:"id"`
	CourseID  int       `db:"course_id"`
	DayOfWeek string    `db:"day_of_week"`
	StartTime string    `db:"start_time"` // TIME type returns as string
	EndTime   string    `db:"end_time"`   // TIME type returns as string
	Location  *string   `db:"location"`
	CreatedAt time.Time `db:"created_at"`
}

func (m *ScheduleModel) ToDomain() *domain.Schedule {
	return &domain.Schedule{
		ID:        m.ID,
		CourseID:  m.CourseID,
		DayOfWeek: m.DayOfWeek,
		StartTime: m.StartTime,
		EndTime:   m.EndTime,
		Location:  derefString(m.Location),
		CreatedAt: m.CreatedAt,
	}
}

func FromDomainSchedule(s *domain.Schedule) *ScheduleModel {
	return &ScheduleModel{
		ID:        s.ID,
		CourseID:  s.CourseID,
		DayOfWeek: s.DayOfWeek,
		StartTime: s.StartTime,
		EndTime:   s.EndTime,
		Location:  refString(s.Location),
		CreatedAt: s.CreatedAt,
	}
}
