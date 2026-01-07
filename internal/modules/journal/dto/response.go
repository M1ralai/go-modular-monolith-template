package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/domain"
)

type JournalResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	EntryDate   time.Time `json:"entry_date"`
	Content     string    `json:"content,omitempty"`
	Mood        string    `json:"mood,omitempty"`
	EnergyLevel int       `json:"energy_level,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func ToJournalResponse(j *domain.JournalEntry) *JournalResponse {
	if j == nil {
		return nil
	}
	return &JournalResponse{ID: j.ID, UserID: j.UserID, EntryDate: j.EntryDate, Content: j.Content, Mood: j.Mood, EnergyLevel: j.EnergyLevel, CreatedAt: j.CreatedAt, UpdatedAt: j.UpdatedAt}
}

func ToJournalResponseList(entries []*domain.JournalEntry) []*JournalResponse {
	result := make([]*JournalResponse, len(entries))
	for i, e := range entries {
		result[i] = ToJournalResponse(e)
	}
	return result
}
