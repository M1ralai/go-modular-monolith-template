package domain

import "time"

type Event struct {
	ID          int
	UserID      int
	LifeAreaID  *int
	Title       string
	Description string
	StartTime   time.Time
	EndTime     *time.Time
	Location    string
	IsAllDay    bool
	IsRecurring bool
	Recurrence  string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *Event) Duration() time.Duration {
	if e.EndTime == nil {
		return 0
	}
	return e.EndTime.Sub(e.StartTime)
}
