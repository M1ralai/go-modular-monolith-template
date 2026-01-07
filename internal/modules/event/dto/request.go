package dto

import "time"

type CreateEventRequest struct {
	LifeAreaID  *int       `json:"life_area_id,omitempty"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description,omitempty"`
	StartTime   time.Time  `json:"start_time" validate:"required"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Location    string     `json:"location,omitempty"`
	IsAllDay    bool       `json:"is_all_day,omitempty"`
	IsRecurring bool       `json:"is_recurring,omitempty"`
	Recurrence  string     `json:"recurrence,omitempty"`
}

type UpdateEventRequest struct {
	LifeAreaID  *int       `json:"life_area_id,omitempty"`
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Location    *string    `json:"location,omitempty"`
	IsAllDay    *bool      `json:"is_all_day,omitempty"`
}
