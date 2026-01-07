package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/dto"
)

type PersonService interface {
	Create(ctx context.Context, req *dto.CreatePersonRequest, userID int) (*dto.PersonResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.PersonResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.PersonResponse, error)
	SearchByTag(ctx context.Context, userID int, tag string) ([]*dto.PersonResponse, error)
	Search(ctx context.Context, userID int, query string) ([]*dto.PersonResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdatePersonRequest, userID int) (*dto.PersonResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
