package domain

import "time"

type JournalEntry struct {
	ID          int
	UserID      int
	EntryDate   time.Time
	Content     string
	Mood        string
	EnergyLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
