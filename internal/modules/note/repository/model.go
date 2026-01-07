package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/domain"
)

type NoteModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	CourseID    *int      `db:"course_id"`
	ComponentID *int      `db:"component_id"`
	LifeAreaID  *int      `db:"life_area_id"`
	Title       string    `db:"title"`
	Content     *string   `db:"content"`
	IsFavorite  bool      `db:"is_favorite"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *NoteModel) ToDomain() *domain.Note {
	if m == nil {
		return nil
	}

	content := ""
	if m.Content != nil {
		content = *m.Content
	}

	return &domain.Note{
		ID:          m.ID,
		UserID:      m.UserID,
		CourseID:    m.CourseID,
		ComponentID: m.ComponentID,
		LifeAreaID:  m.LifeAreaID,
		Title:       m.Title,
		Content:     content,
		IsFavorite:  m.IsFavorite,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromDomain(n *domain.Note) *NoteModel {
	if n == nil {
		return nil
	}

	var content *string
	if n.Content != "" {
		content = &n.Content
	}

	return &NoteModel{
		ID:          n.ID,
		UserID:      n.UserID,
		CourseID:    n.CourseID,
		ComponentID: n.ComponentID,
		LifeAreaID:  n.LifeAreaID,
		Title:       n.Title,
		Content:     content,
		IsFavorite:  n.IsFavorite,
		CreatedAt:   n.CreatedAt,
		UpdatedAt:   n.UpdatedAt,
	}
}

type NoteLinkModel struct {
	ID           int       `db:"id"`
	SourceNoteID int       `db:"source_note_id"`
	TargetNoteID int       `db:"target_note_id"`
	LinkText     *string   `db:"link_text"`
	CreatedAt    time.Time `db:"created_at"`
}

func (m *NoteLinkModel) ToDomain() *domain.NoteLink {
	if m == nil {
		return nil
	}

	linkText := ""
	if m.LinkText != nil {
		linkText = *m.LinkText
	}

	return &domain.NoteLink{
		ID:           m.ID,
		SourceNoteID: m.SourceNoteID,
		TargetNoteID: m.TargetNoteID,
		LinkText:     linkText,
		CreatedAt:    m.CreatedAt,
	}
}

func NoteLinkFromDomain(l *domain.NoteLink) *NoteLinkModel {
	if l == nil {
		return nil
	}

	var linkText *string
	if l.LinkText != "" {
		linkText = &l.LinkText
	}

	return &NoteLinkModel{
		ID:           l.ID,
		SourceNoteID: l.SourceNoteID,
		TargetNoteID: l.TargetNoteID,
		LinkText:     linkText,
		CreatedAt:    l.CreatedAt,
	}
}
