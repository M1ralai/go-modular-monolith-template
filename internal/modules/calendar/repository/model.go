package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/domain"
)

type CalendarIntegrationModel struct {
	ID           int        `db:"id"`
	UserID       int        `db:"user_id"`
	Provider     string     `db:"provider"`
	AccessToken  string     `db:"access_token"`
	RefreshToken *string    `db:"refresh_token"`
	ExpiresAt    *time.Time `db:"expires_at"`
	CalendarID   *string    `db:"calendar_id"`
	IsActive     bool       `db:"is_active"`
	LastSyncAt   *time.Time `db:"last_sync_at"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
}

func (m *CalendarIntegrationModel) ToDomain() *domain.CalendarIntegration {
	if m == nil {
		return nil
	}
	calID, refToken := "", ""
	if m.CalendarID != nil {
		calID = *m.CalendarID
	}
	if m.RefreshToken != nil {
		refToken = *m.RefreshToken
	}
	return &domain.CalendarIntegration{ID: m.ID, UserID: m.UserID, Provider: m.Provider, AccessToken: m.AccessToken, RefreshToken: refToken, ExpiresAt: m.ExpiresAt, CalendarID: calID, IsActive: m.IsActive, LastSyncAt: m.LastSyncAt, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func CalendarIntegrationFromDomain(c *domain.CalendarIntegration) *CalendarIntegrationModel {
	if c == nil {
		return nil
	}
	var calID, refToken *string
	if c.CalendarID != "" {
		calID = &c.CalendarID
	}
	if c.RefreshToken != "" {
		refToken = &c.RefreshToken
	}
	return &CalendarIntegrationModel{ID: c.ID, UserID: c.UserID, Provider: c.Provider, AccessToken: c.AccessToken, RefreshToken: refToken, ExpiresAt: c.ExpiresAt, CalendarID: calID, IsActive: c.IsActive, LastSyncAt: c.LastSyncAt, CreatedAt: c.CreatedAt, UpdatedAt: c.UpdatedAt}
}

type RecurringEventModel struct {
	ID               int        `db:"id"`
	UserID           int        `db:"user_id"`
	CourseScheduleID *int       `db:"course_schedule_id"`
	Title            string     `db:"title"`
	DayOfWeek        string     `db:"day_of_week"`
	StartTime        string     `db:"start_time"`
	EndTime          string     `db:"end_time"`
	Location         *string    `db:"location"`
	StartDate        time.Time  `db:"start_date"`
	EndDate          *time.Time `db:"end_date"`
	GoogleEventID    *string    `db:"google_event_id"`
	AppleEventID     *string    `db:"apple_event_id"`
	LastSyncedAt     *time.Time `db:"last_synced_at"`
	CreatedAt        time.Time  `db:"created_at"`
}

func (m *RecurringEventModel) ToDomain() *domain.RecurringEvent {
	if m == nil {
		return nil
	}
	loc, gID, aID := "", "", ""
	if m.Location != nil {
		loc = *m.Location
	}
	if m.GoogleEventID != nil {
		gID = *m.GoogleEventID
	}
	if m.AppleEventID != nil {
		aID = *m.AppleEventID
	}
	return &domain.RecurringEvent{ID: m.ID, UserID: m.UserID, CourseScheduleID: m.CourseScheduleID, Title: m.Title, DayOfWeek: m.DayOfWeek, StartTime: m.StartTime, EndTime: m.EndTime, Location: loc, StartDate: m.StartDate, EndDate: m.EndDate, GoogleEventID: gID, AppleEventID: aID, LastSyncedAt: m.LastSyncedAt, CreatedAt: m.CreatedAt}
}

type SyncQueueModel struct {
	ID           int        `db:"id"`
	UserID       int        `db:"user_id"`
	EventID      *int       `db:"event_id"`
	Provider     string     `db:"provider"`
	Action       string     `db:"action"`
	Status       string     `db:"status"`
	RetryCount   int        `db:"retry_count"`
	ErrorMessage *string    `db:"error_message"`
	CreatedAt    time.Time  `db:"created_at"`
	SyncedAt     *time.Time `db:"synced_at"`
}

func (m *SyncQueueModel) ToDomain() *domain.SyncQueue {
	if m == nil {
		return nil
	}
	errMsg := ""
	if m.ErrorMessage != nil {
		errMsg = *m.ErrorMessage
	}
	return &domain.SyncQueue{ID: m.ID, UserID: m.UserID, EventID: m.EventID, Provider: m.Provider, Action: m.Action, Status: m.Status, RetryCount: m.RetryCount, ErrorMessage: errMsg, CreatedAt: m.CreatedAt, SyncedAt: m.SyncedAt}
}
