package domain

import "time"

type TimeSlot struct {
	Start           time.Time
	End             time.Time
	DurationMinutes int
}

func (t *TimeSlot) CanFit(durationMinutes int) bool {
	return t.DurationMinutes >= durationMinutes
}

type Availability struct {
	Date      time.Time
	FreeSlots []TimeSlot
}

type ConflictResult struct {
	HasConflict bool
	Reason      string
	Suggestions []TimeSlot
}
