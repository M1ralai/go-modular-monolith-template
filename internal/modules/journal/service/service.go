package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/dto"
)

type JournalService interface {
	Create(ctx context.Context, req *dto.CreateJournalRequest, userID int) (*dto.JournalResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.JournalResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.JournalResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateJournalRequest, userID int) (*dto.JournalResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
