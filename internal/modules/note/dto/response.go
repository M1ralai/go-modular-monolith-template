package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/domain"
)

type NoteResponse struct {
	ID            int             `json:"id"`
	UserID        int             `json:"user_id"`
	CourseID      *int            `json:"course_id,omitempty"`
	ComponentID   *int            `json:"component_id,omitempty"`
	LifeAreaID    *int            `json:"life_area_id,omitempty"`
	Title         string          `json:"title"`
	Content       string          `json:"content,omitempty"`
	IsFavorite    bool            `json:"is_favorite"`
	OutgoingLinks []*NoteLinkInfo `json:"outgoing_links,omitempty"`
	Backlinks     []*NoteLinkInfo `json:"backlinks,omitempty"`
	BacklinkCount int             `json:"backlink_count"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type NoteLinkInfo struct {
	ID       int    `json:"id"`
	NoteID   int    `json:"note_id"`
	Title    string `json:"title"`
	LinkText string `json:"link_text,omitempty"`
}

type NoteLinkResponse struct {
	ID           int       `json:"id"`
	SourceNoteID int       `json:"source_note_id"`
	TargetNoteID int       `json:"target_note_id"`
	LinkText     string    `json:"link_text,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToNoteResponse(n *domain.Note) *NoteResponse {
	if n == nil {
		return nil
	}

	return &NoteResponse{
		ID:          n.ID,
		UserID:      n.UserID,
		CourseID:    n.CourseID,
		ComponentID: n.ComponentID,
		LifeAreaID:  n.LifeAreaID,
		Title:       n.Title,
		Content:     n.Content,
		IsFavorite:  n.IsFavorite,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}

func ToNoteResponseList(notes []*domain.Note) []*NoteResponse {
	result := make([]*NoteResponse, len(notes))
	for i, n := range notes {
		result[i] = ToNoteResponse(n)
	}
	return result
}

func ToNoteLinkResponse(l *domain.NoteLink) *NoteLinkResponse {
	if l == nil {
		return nil
	}

	return &NoteLinkResponse{
		ID:           l.ID,
		SourceNoteID: l.SourceNoteID,
		TargetNoteID: l.TargetNoteID,
		LinkText:     l.LinkText,
		CreatedAt:    l.CreatedAt,
	}
}
