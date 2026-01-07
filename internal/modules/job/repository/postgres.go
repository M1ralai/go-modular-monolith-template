package repository

import (
	"context"
	"database/sql"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/job/domain"
	"github.com/jmoiron/sqlx"
)

type postgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL job repository
func NewPostgresRepository(db *sqlx.DB) JobRepository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, execution *domain.JobExecution) (*domain.JobExecution, error) {
	query := `
		INSERT INTO job_executions (job_name, status, started_at, completed_at, error, result, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query,
		execution.JobName,
		execution.Status,
		execution.StartedAt,
		execution.CompletedAt,
		execution.Error,
		execution.Result,
		execution.DurationMs,
	).Scan(&execution.ID, &execution.CreatedAt)

	if err != nil {
		return nil, err
	}
	return execution, nil
}

func (r *postgresRepository) Update(ctx context.Context, execution *domain.JobExecution) error {
	query := `
		UPDATE job_executions
		SET status = $1, completed_at = $2, error = $3, result = $4, duration_ms = $5
		WHERE id = $6`

	_, err := r.db.ExecContext(ctx, query,
		execution.Status,
		execution.CompletedAt,
		execution.Error,
		execution.Result,
		execution.DurationMs,
		execution.ID,
	)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id int) (*domain.JobExecution, error) {
	var execution domain.JobExecution
	query := `SELECT * FROM job_executions WHERE id = $1`
	err := r.db.GetContext(ctx, &execution, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *postgresRepository) GetByJobName(ctx context.Context, jobName string, limit int) ([]*domain.JobExecution, error) {
	var executions []*domain.JobExecution
	query := `SELECT * FROM job_executions WHERE job_name = $1 ORDER BY started_at DESC LIMIT $2`
	err := r.db.SelectContext(ctx, &executions, query, jobName, limit)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (r *postgresRepository) GetLatestByJobName(ctx context.Context, jobName string) (*domain.JobExecution, error) {
	var execution domain.JobExecution
	query := `SELECT * FROM job_executions WHERE job_name = $1 ORDER BY started_at DESC LIMIT 1`
	err := r.db.GetContext(ctx, &execution, query, jobName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &execution, nil
}

func (r *postgresRepository) GetRunning(ctx context.Context) ([]*domain.JobExecution, error) {
	var executions []*domain.JobExecution
	query := `SELECT * FROM job_executions WHERE status = 'running' ORDER BY started_at DESC`
	err := r.db.SelectContext(ctx, &executions, query)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (r *postgresRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.JobExecution, error) {
	var executions []*domain.JobExecution
	query := `SELECT * FROM job_executions ORDER BY started_at DESC LIMIT $1 OFFSET $2`
	err := r.db.SelectContext(ctx, &executions, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return executions, nil
}

func (r *postgresRepository) DeleteOlderThan(ctx context.Context, days int) (int64, error) {
	query := `DELETE FROM job_executions WHERE created_at < NOW() - ($1 || ' days')::interval`
	result, err := r.db.ExecContext(ctx, query, days)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
