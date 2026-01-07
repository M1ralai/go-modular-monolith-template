package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) LifeAreaRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, lifeArea *domain.LifeArea) (*domain.LifeArea, error) {
	query := `
		INSERT INTO life_areas (user_id, name, icon, color, display_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	now := time.Now()
	model := FromDomain(lifeArea)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.UserID,
		model.Name,
		model.Icon,
		model.Color,
		model.DisplayOrder,
		now,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.LifeArea, error) {
	query := `
		SELECT id, user_id, name, icon, color, display_order, created_at
		FROM life_areas
		WHERE id = $1
	`

	var model LifeAreaModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.LifeArea, error) {
	query := `
		SELECT id, user_id, name, icon, color, display_order, created_at
		FROM life_areas
		WHERE user_id = $1
		ORDER BY display_order ASC, created_at ASC
	`

	var models []LifeAreaModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	areas := make([]*domain.LifeArea, len(models))
	for i, m := range models {
		areas[i] = m.ToDomain()
	}

	return areas, nil
}

func (r *postgresRepository) Update(ctx context.Context, lifeArea *domain.LifeArea) error {
	query := `
		UPDATE life_areas
		SET name = $1, icon = $2, color = $3, display_order = $4
		WHERE id = $5
	`

	model := FromDomain(lifeArea)
	_, err := r.db.ExecContext(
		ctx, query,
		model.Name,
		model.Icon,
		model.Color,
		model.DisplayOrder,
		model.ID,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM life_areas WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
