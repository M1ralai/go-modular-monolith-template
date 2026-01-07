package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/domain"
)

type IntegrationResponse struct {
	ID         int        `json:"id"`
	Provider   string     `json:"provider"`
	IsActive   bool       `json:"is_active"`
	CalendarID string     `json:"calendar_id,omitempty"`
	LastSyncAt *time.Time `json:"last_sync_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

func ToIntegrationResponse(c *domain.CalendarIntegration) *IntegrationResponse {
	if c == nil {
		return nil
	}
	return &IntegrationResponse{ID: c.ID, Provider: c.Provider, IsActive: c.IsActive, CalendarID: c.CalendarID, LastSyncAt: c.LastSyncAt, CreatedAt: c.CreatedAt}
}

type SyncStatusResponse struct {
	Integrations []*IntegrationResponse `json:"integrations"`
}

type AuthURLResponse struct {
	AuthURL string `json:"auth_url"`
}

type RecurringEventResponse struct {
	ID               int        `json:"id"`
	Title            string     `json:"title"`
	DayOfWeek        string     `json:"day_of_week"`
	StartTime        string     `json:"start_time"`
	EndTime          string     `json:"end_time"`
	Location         string     `json:"location,omitempty"`
	StartDate        time.Time  `json:"start_date"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	CourseScheduleID *int       `json:"course_schedule_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

func ToRecurringEventResponse(e *domain.RecurringEvent) *RecurringEventResponse {
	if e == nil {
		return nil
	}
	return &RecurringEventResponse{ID: e.ID, Title: e.Title, DayOfWeek: e.DayOfWeek, StartTime: e.StartTime, EndTime: e.EndTime, Location: e.Location, StartDate: e.StartDate, EndDate: e.EndDate, CourseScheduleID: e.CourseScheduleID, CreatedAt: e.CreatedAt}
}
