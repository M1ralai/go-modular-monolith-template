package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/domain"
)

type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) (*domain.Note, error)
	GetByID(ctx context.Context, id int) (*domain.Note, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Note, error)
	GetByCourseID(ctx context.Context, courseID int) ([]*domain.Note, error)
	GetByLifeAreaID(ctx context.Context, lifeAreaID int) ([]*domain.Note, error)
	GetFavorites(ctx context.Context, userID int) ([]*domain.Note, error)
	Search(ctx context.Context, userID int, query string) ([]*domain.Note, error)
	Update(ctx context.Context, note *domain.Note) error
	Delete(ctx context.Context, id int) error

	CreateLink(ctx context.Context, link *domain.NoteLink) (*domain.NoteLink, error)
	GetOutgoingLinks(ctx context.Context, noteID int) ([]*domain.NoteLink, error)
	GetBacklinks(ctx context.Context, noteID int) ([]*domain.NoteLink, error)
	DeleteLink(ctx context.Context, id int) error
	DeleteLinksByNote(ctx context.Context, noteID int) error
}
