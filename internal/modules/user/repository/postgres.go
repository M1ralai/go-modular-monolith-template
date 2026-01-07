package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/user/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) UserRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `
		INSERT INTO users (email, password_hash, full_name, avatar_url, timezone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	model := FromDomain(user)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.Email,
		model.PasswordHash,
		model.FullName,
		model.AvatarURL,
		model.Timezone,
		now,
		now,
	).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, timezone, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var model UserModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, timezone, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var model UserModel
	err := r.db.GetContext(ctx, &model, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, email, password_hash, full_name, avatar_url, timezone, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
	`

	var models []UserModel
	err := r.db.SelectContext(ctx, &models, query)
	if err != nil {
		return nil, err
	}

	users := make([]*domain.User, len(models))
	for i, m := range models {
		users[i] = m.ToDomain()
	}

	return users, nil
}

func (r *postgresRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET email = $1, password_hash = $2, full_name = $3, avatar_url = $4, timezone = $5, updated_at = $6
		WHERE id = $7
	`

	model := FromDomain(user)
	_, err := r.db.ExecContext(
		ctx, query,
		model.Email,
		model.PasswordHash,
		model.FullName,
		model.AvatarURL,
		model.Timezone,
		time.Now(),
		model.ID,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
