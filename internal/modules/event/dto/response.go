package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/domain"
)

type EventResponse struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	LifeAreaID  *int       `json:"life_area_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	Location    string     `json:"location,omitempty"`
	IsAllDay    bool       `json:"is_all_day"`
	IsRecurring bool       `json:"is_recurring"`
	Recurrence  string     `json:"recurrence,omitempty"`
	DurationMin int        `json:"duration_min,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToEventResponse(e *domain.Event) *EventResponse {
	if e == nil {
		return nil
	}
	durationMin := int(e.Duration().Minutes())
	return &EventResponse{ID: e.ID, UserID: e.UserID, LifeAreaID: e.LifeAreaID, Title: e.Title, Description: e.Description, StartTime: e.StartTime, EndTime: e.EndTime, Location: e.Location, IsAllDay: e.IsAllDay, IsRecurring: e.IsRecurring, Recurrence: e.Recurrence, DurationMin: durationMin, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt}
}

func ToEventResponseList(events []*domain.Event) []*EventResponse {
	result := make([]*EventResponse, len(events))
	for i, e := range events {
		result[i] = ToEventResponse(e)
	}
	return result
}
