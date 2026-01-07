package dto

import "time"

type CreateJournalRequest struct {
	EntryDate   time.Time `json:"entry_date" validate:"required"`
	Content     string    `json:"content,omitempty"`
	Mood        string    `json:"mood,omitempty"`
	EnergyLevel int       `json:"energy_level,omitempty" validate:"min=0,max=10"`
}

type UpdateJournalRequest struct {
	Content     *string `json:"content,omitempty"`
	Mood        *string `json:"mood,omitempty"`
	EnergyLevel *int    `json:"energy_level,omitempty" validate:"omitempty,min=0,max=10"`
}
