package dto

import "time"

type ConflictCheckRequest struct {
	Start time.Time `json:"start" validate:"required"`
	End   time.Time `json:"end" validate:"required"`
}

type FreeSlotsRequest struct {
	Date            time.Time `json:"date" validate:"required"`
	DurationMinutes int       `json:"duration_minutes" validate:"required,min=15"`
}

type CreateBlockedSlotRequest struct {
	SourceType    string    `json:"source_type,omitempty"`
	SourceID      *int      `json:"source_id,omitempty"`
	StartDatetime time.Time `json:"start_datetime" validate:"required"`
	EndDatetime   time.Time `json:"end_datetime" validate:"required"`
	Reason        string    `json:"reason" validate:"required"`
	IsFlexible    bool      `json:"is_flexible"`
}

type GenerateEventsRequest struct {
	CourseScheduleID  *int        `json:"course_schedule_id,omitempty"`
	Title             string      `json:"title" validate:"required"`
	DayOfWeek         string      `json:"day_of_week" validate:"required,oneof=Monday Tuesday Wednesday Thursday Friday Saturday Sunday"`
	StartTime         string      `json:"start_time" validate:"required"`
	EndTime           string      `json:"end_time" validate:"required"`
	Location          string      `json:"location,omitempty"`
	SemesterStartDate time.Time   `json:"semester_start_date" validate:"required"`
	SemesterEndDate   time.Time   `json:"semester_end_date" validate:"required"`
	ExcludeDates      []time.Time `json:"exclude_dates,omitempty"`
}
