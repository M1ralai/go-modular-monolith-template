package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/schedule/domain"
)

type BlockedTimeSlotModel struct {
	ID            int       `db:"id"`
	UserID        int       `db:"user_id"`
	SourceType    *string   `db:"source_type"`
	SourceID      *int      `db:"source_id"`
	StartDatetime time.Time `db:"start_datetime"`
	EndDatetime   time.Time `db:"end_datetime"`
	Reason        *string   `db:"reason"`
	IsFlexible    bool      `db:"is_flexible"`
	CreatedAt     time.Time `db:"created_at"`
}

func (m *BlockedTimeSlotModel) ToDomain() *domain.BlockedTimeSlot {
	if m == nil {
		return nil
	}
	srcType, reason := "", ""
	if m.SourceType != nil {
		srcType = *m.SourceType
	}
	if m.Reason != nil {
		reason = *m.Reason
	}
	return &domain.BlockedTimeSlot{ID: m.ID, UserID: m.UserID, SourceType: srcType, SourceID: m.SourceID, StartDatetime: m.StartDatetime, EndDatetime: m.EndDatetime, Reason: reason, IsFlexible: m.IsFlexible, CreatedAt: m.CreatedAt}
}

func BlockedTimeSlotFromDomain(b *domain.BlockedTimeSlot) *BlockedTimeSlotModel {
	if b == nil {
		return nil
	}
	var srcType, reason *string
	if b.SourceType != "" {
		srcType = &b.SourceType
	}
	if b.Reason != "" {
		reason = &b.Reason
	}
	return &BlockedTimeSlotModel{ID: b.ID, UserID: b.UserID, SourceType: srcType, SourceID: b.SourceID, StartDatetime: b.StartDatetime, EndDatetime: b.EndDatetime, Reason: reason, IsFlexible: b.IsFlexible, CreatedAt: b.CreatedAt}
}
