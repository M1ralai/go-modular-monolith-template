package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/domain"
)

type ConflictResponse struct {
	HasConflict bool                `json:"has_conflict"`
	Reason      string              `json:"reason,omitempty"`
	Suggestions []*TimeSlotResponse `json:"suggestions,omitempty"`
}

type TimeSlotResponse struct {
	Start           time.Time `json:"start"`
	End             time.Time `json:"end"`
	DurationMinutes int       `json:"duration_minutes"`
}

type BlockedSlotResponse struct {
	ID            int       `json:"id"`
	SourceType    string    `json:"source_type,omitempty"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Reason        string    `json:"reason"`
	IsFlexible    bool      `json:"is_flexible"`
	CreatedAt     time.Time `json:"created_at"`
}

func ToBlockedSlotResponse(s *domain.BlockedTimeSlot) *BlockedSlotResponse {
	if s == nil {
		return nil
	}
	return &BlockedSlotResponse{ID: s.ID, SourceType: s.SourceType, StartDatetime: s.StartDatetime, EndDatetime: s.EndDatetime, Reason: s.Reason, IsFlexible: s.IsFlexible, CreatedAt: s.CreatedAt}
}

type GenerateEventsResponse struct {
	EventsGenerated     int `json:"events_generated"`
	BlockedSlotsCreated int `json:"blocked_slots_created"`
}
