package domain

import "time"

type RecurringEvent struct {
	ID               int
	UserID           int
	CourseScheduleID *int
	Title            string
	DayOfWeek        string
	StartTime        string
	EndTime          string
	Location         string
	StartDate        time.Time
	EndDate          *time.Time
	ExcludeDates     []time.Time
	GoogleEventID    string
	AppleEventID     string
	LastSyncedAt     *time.Time
	CreatedAt        time.Time
}

func (r *RecurringEvent) ShouldExclude(date time.Time) bool {
	for _, d := range r.ExcludeDates {
		if d.Year() == date.Year() && d.Month() == date.Month() && d.Day() == date.Day() {
			return true
		}
	}
	return false
}
