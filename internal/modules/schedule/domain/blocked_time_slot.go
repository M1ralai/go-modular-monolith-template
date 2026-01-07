package domain

import "time"

type BlockedTimeSlot struct {
	ID            int
	UserID        int
	SourceType    string
	SourceID      *int
	StartDatetime time.Time
	EndDatetime   time.Time
	Reason        string
	IsFlexible    bool
	CreatedAt     time.Time
}

func (b *BlockedTimeSlot) OverlapsWith(start, end time.Time) bool {
	return b.StartDatetime.Before(end) && start.Before(b.EndDatetime)
}

func (b *BlockedTimeSlot) Duration() time.Duration {
	return b.EndDatetime.Sub(b.StartDatetime)
}
