package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/domain"
)

type PersonRepository interface {
	Create(ctx context.Context, person *domain.Person) (*domain.Person, error)
	GetByID(ctx context.Context, id int) (*domain.Person, error)
	GetByUserID(ctx context.Context, userID int) ([]*domain.Person, error)
	SearchByTag(ctx context.Context, userID int, tag string) ([]*domain.Person, error)
	Search(ctx context.Context, userID int, query string) ([]*domain.Person, error)
	Update(ctx context.Context, person *domain.Person) error
	Delete(ctx context.Context, id int) error
}
