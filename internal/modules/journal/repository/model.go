package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/domain"
)

type JournalModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	EntryDate   time.Time `db:"entry_date"`
	Content     *string   `db:"content"`
	Mood        *string   `db:"mood"`
	EnergyLevel *int      `db:"energy_level"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *JournalModel) ToDomain() *domain.JournalEntry {
	if m == nil {
		return nil
	}
	content, mood := "", ""
	energyLevel := 0
	if m.Content != nil {
		content = *m.Content
	}
	if m.Mood != nil {
		mood = *m.Mood
	}
	if m.EnergyLevel != nil {
		energyLevel = *m.EnergyLevel
	}
	return &domain.JournalEntry{ID: m.ID, UserID: m.UserID, EntryDate: m.EntryDate, Content: content, Mood: mood, EnergyLevel: energyLevel, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func FromDomain(j *domain.JournalEntry) *JournalModel {
	if j == nil {
		return nil
	}
	var content, mood *string
	var energyLevel *int
	if j.Content != "" {
		content = &j.Content
	}
	if j.Mood != "" {
		mood = &j.Mood
	}
	if j.EnergyLevel > 0 {
		energyLevel = &j.EnergyLevel
	}
	return &JournalModel{ID: j.ID, UserID: j.UserID, EntryDate: j.EntryDate, Content: content, Mood: mood, EnergyLevel: energyLevel, CreatedAt: j.CreatedAt, UpdatedAt: j.UpdatedAt}
}
