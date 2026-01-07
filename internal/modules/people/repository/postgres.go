package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/people/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct{ db *sqlx.DB }

func NewPostgresRepository(db *sqlx.DB) PersonRepository { return &postgresRepository{db: db} }

func (r *postgresRepository) Create(ctx context.Context, person *domain.Person) (*domain.Person, error) {
	query := `INSERT INTO people (user_id, name, email, phone, company, relationship, tags, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id, created_at, updated_at`
	now := time.Now()
	model := FromDomain(person)
	err := r.db.QueryRowxContext(ctx, query, model.UserID, model.Name, model.Email, model.Phone, model.Company, model.Relationship, model.Tags, model.Notes, now, now).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Person, error) {
	query := `SELECT id, user_id, name, email, phone, company, relationship, tags, notes, created_at, updated_at FROM people WHERE id = $1`
	var model PersonModel
	if err := r.db.GetContext(ctx, &model, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Person, error) {
	query := `SELECT id, user_id, name, email, phone, company, relationship, tags, notes, created_at, updated_at FROM people WHERE user_id = $1 ORDER BY name`
	var models []PersonModel
	if err := r.db.SelectContext(ctx, &models, query, userID); err != nil {
		return nil, err
	}
	people := make([]*domain.Person, len(models))
	for i, m := range models {
		people[i] = m.ToDomain()
	}
	return people, nil
}

func (r *postgresRepository) SearchByTag(ctx context.Context, userID int, tag string) ([]*domain.Person, error) {
	query := `SELECT id, user_id, name, email, phone, company, relationship, tags, notes, created_at, updated_at FROM people WHERE user_id = $1 AND $2 = ANY(tags) ORDER BY name`
	var models []PersonModel
	if err := r.db.SelectContext(ctx, &models, query, userID, tag); err != nil {
		return nil, err
	}
	people := make([]*domain.Person, len(models))
	for i, m := range models {
		people[i] = m.ToDomain()
	}
	return people, nil
}

func (r *postgresRepository) Search(ctx context.Context, userID int, query string) ([]*domain.Person, error) {
	sqlQuery := `SELECT id, user_id, name, email, phone, company, relationship, tags, notes, created_at, updated_at FROM people WHERE user_id = $1 AND (name ILIKE $2 OR company ILIKE $2 OR email ILIKE $2) ORDER BY name`
	var models []PersonModel
	if err := r.db.SelectContext(ctx, &models, sqlQuery, userID, "%"+query+"%"); err != nil {
		return nil, err
	}
	people := make([]*domain.Person, len(models))
	for i, m := range models {
		people[i] = m.ToDomain()
	}
	return people, nil
}

func (r *postgresRepository) Update(ctx context.Context, person *domain.Person) error {
	query := `UPDATE people SET name = $1, email = $2, phone = $3, company = $4, relationship = $5, tags = $6, notes = $7, updated_at = $8 WHERE id = $9`
	model := FromDomain(person)
	_, err := r.db.ExecContext(ctx, query, model.Name, model.Email, model.Phone, model.Company, model.Relationship, model.Tags, model.Notes, time.Now(), model.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM people WHERE id = $1`, id)
	return err
}
