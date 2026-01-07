package repository

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/job/domain"
)

// JobRepository defines the interface for job execution persistence
type JobRepository interface {
	// Create creates a new job execution record
	Create(ctx context.Context, execution *domain.JobExecution) (*domain.JobExecution, error)

	// Update updates an existing job execution record
	Update(ctx context.Context, execution *domain.JobExecution) error

	// GetByID returns a job execution by ID
	GetByID(ctx context.Context, id int) (*domain.JobExecution, error)

	// GetByJobName returns executions for a specific job
	GetByJobName(ctx context.Context, jobName string, limit int) ([]*domain.JobExecution, error)

	// GetLatestByJobName returns the most recent execution for a job
	GetLatestByJobName(ctx context.Context, jobName string) (*domain.JobExecution, error)

	// GetRunning returns all currently running job executions
	GetRunning(ctx context.Context) ([]*domain.JobExecution, error)

	// GetAll returns all job executions with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*domain.JobExecution, error)

	// DeleteOlderThan removes job executions older than a specified time
	DeleteOlderThan(ctx context.Context, days int) (int64, error)
}
