package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/note/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) NoteRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, note *domain.Note) (*domain.Note, error) {
	query := `
		INSERT INTO notes (user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	model := FromDomain(note)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.UserID,
		model.CourseID,
		model.ComponentID,
		model.LifeAreaID,
		model.Title,
		model.Content,
		model.IsFavorite,
		now,
		now,
	).Scan(&model.ID, &model.CreatedAt, &model.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.Note, error) {
	query := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE id = $1
	`

	var model NoteModel
	err := r.db.GetContext(ctx, &model, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetByUserID(ctx context.Context, userID int) ([]*domain.Note, error) {
	query := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY updated_at DESC
	`

	var models []NoteModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, len(models))
	for i, m := range models {
		notes[i] = m.ToDomain()
	}

	return notes, nil
}

func (r *postgresRepository) GetByCourseID(ctx context.Context, courseID int) ([]*domain.Note, error) {
	query := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE course_id = $1
		ORDER BY updated_at DESC
	`

	var models []NoteModel
	err := r.db.SelectContext(ctx, &models, query, courseID)
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, len(models))
	for i, m := range models {
		notes[i] = m.ToDomain()
	}

	return notes, nil
}

func (r *postgresRepository) GetByLifeAreaID(ctx context.Context, lifeAreaID int) ([]*domain.Note, error) {
	query := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE life_area_id = $1
		ORDER BY updated_at DESC
	`

	var models []NoteModel
	err := r.db.SelectContext(ctx, &models, query, lifeAreaID)
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, len(models))
	for i, m := range models {
		notes[i] = m.ToDomain()
	}

	return notes, nil
}

func (r *postgresRepository) GetFavorites(ctx context.Context, userID int) ([]*domain.Note, error) {
	query := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE user_id = $1 AND is_favorite = true
		ORDER BY updated_at DESC
	`

	var models []NoteModel
	err := r.db.SelectContext(ctx, &models, query, userID)
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, len(models))
	for i, m := range models {
		notes[i] = m.ToDomain()
	}

	return notes, nil
}

func (r *postgresRepository) Search(ctx context.Context, userID int, query string) ([]*domain.Note, error) {
	sqlQuery := `
		SELECT id, user_id, course_id, component_id, life_area_id, title, content, is_favorite, created_at, updated_at
		FROM notes
		WHERE user_id = $1 AND (title ILIKE $2 OR content ILIKE $2)
		ORDER BY updated_at DESC
	`

	searchPattern := "%" + query + "%"
	var models []NoteModel
	err := r.db.SelectContext(ctx, &models, sqlQuery, userID, searchPattern)
	if err != nil {
		return nil, err
	}

	notes := make([]*domain.Note, len(models))
	for i, m := range models {
		notes[i] = m.ToDomain()
	}

	return notes, nil
}

func (r *postgresRepository) Update(ctx context.Context, note *domain.Note) error {
	query := `
		UPDATE notes
		SET course_id = $1, component_id = $2, life_area_id = $3, title = $4, content = $5, is_favorite = $6, updated_at = $7
		WHERE id = $8
	`

	model := FromDomain(note)
	_, err := r.db.ExecContext(
		ctx, query,
		model.CourseID,
		model.ComponentID,
		model.LifeAreaID,
		model.Title,
		model.Content,
		model.IsFavorite,
		time.Now(),
		model.ID,
	)

	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM notes WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) CreateLink(ctx context.Context, link *domain.NoteLink) (*domain.NoteLink, error) {
	query := `
		INSERT INTO note_links (source_note_id, target_note_id, link_text, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	now := time.Now()
	model := NoteLinkFromDomain(link)

	err := r.db.QueryRowxContext(
		ctx, query,
		model.SourceNoteID,
		model.TargetNoteID,
		model.LinkText,
		now,
	).Scan(&model.ID, &model.CreatedAt)

	if err != nil {
		return nil, err
	}

	return model.ToDomain(), nil
}

func (r *postgresRepository) GetOutgoingLinks(ctx context.Context, noteID int) ([]*domain.NoteLink, error) {
	query := `
		SELECT id, source_note_id, target_note_id, link_text, created_at
		FROM note_links
		WHERE source_note_id = $1
		ORDER BY created_at
	`

	var models []NoteLinkModel
	err := r.db.SelectContext(ctx, &models, query, noteID)
	if err != nil {
		return nil, err
	}

	links := make([]*domain.NoteLink, len(models))
	for i, m := range models {
		links[i] = m.ToDomain()
	}

	return links, nil
}

func (r *postgresRepository) GetBacklinks(ctx context.Context, noteID int) ([]*domain.NoteLink, error) {
	query := `
		SELECT id, source_note_id, target_note_id, link_text, created_at
		FROM note_links
		WHERE target_note_id = $1
		ORDER BY created_at
	`

	var models []NoteLinkModel
	err := r.db.SelectContext(ctx, &models, query, noteID)
	if err != nil {
		return nil, err
	}

	links := make([]*domain.NoteLink, len(models))
	for i, m := range models {
		links[i] = m.ToDomain()
	}

	return links, nil
}

func (r *postgresRepository) DeleteLink(ctx context.Context, id int) error {
	query := `DELETE FROM note_links WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *postgresRepository) DeleteLinksByNote(ctx context.Context, noteID int) error {
	query := `DELETE FROM note_links WHERE source_note_id = $1 OR target_note_id = $1`
	_, err := r.db.ExecContext(ctx, query, noteID)
	return err
}
