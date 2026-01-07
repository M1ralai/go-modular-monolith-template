package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/event/domain"
)

type EventModel struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	LifeAreaID  *int       `db:"life_area_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	StartTime   time.Time  `db:"start_time"`
	EndTime     *time.Time `db:"end_time"`
	Location    *string    `db:"location"`
	IsAllDay    bool       `db:"is_all_day"`
	IsRecurring bool       `db:"is_recurring"`
	Recurrence  *string    `db:"recurrence"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (m *EventModel) ToDomain() *domain.Event {
	if m == nil {
		return nil
	}
	desc, loc, rec := "", "", ""
	if m.Description != nil {
		desc = *m.Description
	}
	if m.Location != nil {
		loc = *m.Location
	}
	if m.Recurrence != nil {
		rec = *m.Recurrence
	}
	return &domain.Event{ID: m.ID, UserID: m.UserID, LifeAreaID: m.LifeAreaID, Title: m.Title, Description: desc, StartTime: m.StartTime, EndTime: m.EndTime, Location: loc, IsAllDay: m.IsAllDay, IsRecurring: m.IsRecurring, Recurrence: rec, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func FromDomain(e *domain.Event) *EventModel {
	if e == nil {
		return nil
	}
	var desc, loc, rec *string
	if e.Description != "" {
		desc = &e.Description
	}
	if e.Location != "" {
		loc = &e.Location
	}
	if e.Recurrence != "" {
		rec = &e.Recurrence
	}
	return &EventModel{ID: e.ID, UserID: e.UserID, LifeAreaID: e.LifeAreaID, Title: e.Title, Description: desc, StartTime: e.StartTime, EndTime: e.EndTime, Location: loc, IsAllDay: e.IsAllDay, IsRecurring: e.IsRecurring, Recurrence: rec, CreatedAt: e.CreatedAt, UpdatedAt: e.UpdatedAt}
}
