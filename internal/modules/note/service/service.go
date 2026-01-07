package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/dto"
)

type NoteService interface {
	Create(ctx context.Context, req *dto.CreateNoteRequest, userID int) (*dto.NoteResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.NoteResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.NoteResponse, error)
	GetByCourse(ctx context.Context, courseID, userID int) ([]*dto.NoteResponse, error)
	GetByLifeArea(ctx context.Context, lifeAreaID, userID int) ([]*dto.NoteResponse, error)
	GetFavorites(ctx context.Context, userID int) ([]*dto.NoteResponse, error)
	Search(ctx context.Context, userID int, query string) ([]*dto.NoteResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateNoteRequest, userID int) (*dto.NoteResponse, error)
	Delete(ctx context.Context, id, userID int) error

	CreateLink(ctx context.Context, sourceNoteID int, req *dto.CreateNoteLinkRequest, userID int) (*dto.NoteLinkResponse, error)
	GetBacklinks(ctx context.Context, noteID, userID int) ([]*dto.NoteLinkInfo, error)
	DeleteLink(ctx context.Context, linkID, userID int) error
}
