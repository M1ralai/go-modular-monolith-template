package domain

import "time"

type NoteLink struct {
	ID           int
	SourceNoteID int
	TargetNoteID int
	LinkText     string
	CreatedAt    time.Time
}
