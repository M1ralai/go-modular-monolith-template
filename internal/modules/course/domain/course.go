package domain

import "time"

type Course struct {
	ID          int
	UserID      int
	Name        string
	Code        string
	Instructor  string
	Credits     float64
	Semester    string
	Type        string
	Color       string
	SyllabusURL string
	FinalGrade  string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Components  []*Component
	Schedules   []*Schedule
}

type Component struct {
	ID             int
	CourseID       int
	Type           string
	Name           string
	Weight         float64
	MaxScore       float64
	AchievedScore  *float64
	DueDate        *time.Time
	CompletionDate *time.Time
	IsCompleted    bool
	Notes          string
	DisplayOrder   int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Schedule struct {
	ID        int
	CourseID  int
	DayOfWeek string
	StartTime string // HH:MM format
	EndTime   string // HH:MM format
	Location  string
	CreatedAt time.Time
}

func (c *Course) IsCompleted() bool {
	return c.FinalGrade != ""
}
