package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/job/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/job/repository"
)

// JobService provides job management functionality
type JobService interface {
	// TriggerJob manually triggers a job by name
	TriggerJob(ctx context.Context, jobName string) (*domain.JobExecution, error)

	// GetJobStatus returns the current status of a job
	GetJobStatus(ctx context.Context, jobName string) (*domain.JobExecution, error)

	// GetJobHistory returns the execution history for a job
	GetJobHistory(ctx context.Context, jobName string, limit int) ([]*domain.JobExecution, error)

	// ListJobs returns all registered jobs
	ListJobs(ctx context.Context) ([]JobInfo, error)

	// GetRunningJobs returns all currently running jobs
	GetRunningJobs(ctx context.Context) ([]*domain.JobExecution, error)
}

// JobInfo represents basic job information
type JobInfo struct {
	Name           string `json:"name"`
	Schedule       string `json:"schedule"`
	TimeoutSeconds int    `json:"timeout_seconds"`
	LastRunStatus  string `json:"last_run_status,omitempty"`
	LastRunAt      string `json:"last_run_at,omitempty"`
}

type jobService struct {
	repo      repository.JobRepository
	scheduler *jobs.Scheduler
	logger    *logger.ZapLogger
}

// NewJobService creates a new job service
func NewJobService(repo repository.JobRepository, scheduler *jobs.Scheduler, logger *logger.ZapLogger) JobService {
	return &jobService{
		repo:      repo,
		scheduler: scheduler,
		logger:    logger,
	}
}

func (s *jobService) TriggerJob(ctx context.Context, jobName string) (*domain.JobExecution, error) {
	s.logger.Info("Triggering job manually", map[string]interface{}{
		"job":    jobName,
		"action": "JOB_TRIGGER",
	})

	// Create execution record
	execution := domain.NewJobExecution(jobName)
	execution.MarkRunning()

	created, err := s.repo.Create(ctx, execution)
	if err != nil {
		s.logger.Error("Failed to create job execution record", err, map[string]interface{}{
			"job":    jobName,
			"action": "JOB_RECORD_CREATE_FAILED",
		})
		return nil, err
	}

	// Trigger job asynchronously
	go func() {
		result := s.scheduler.TriggerJob(jobName)

		// Update execution record
		if result.Error != nil {
			created.MarkFailed(result.Error)
		} else {
			created.MarkCompleted(result.Result)
		}

		if err := s.repo.Update(context.Background(), created); err != nil {
			s.logger.Error("Failed to update job execution record", err, map[string]interface{}{
				"job":    jobName,
				"action": "JOB_RECORD_UPDATE_FAILED",
			})
		}
	}()

	return created, nil
}

func (s *jobService) GetJobStatus(ctx context.Context, jobName string) (*domain.JobExecution, error) {
	return s.repo.GetLatestByJobName(ctx, jobName)
}

func (s *jobService) GetJobHistory(ctx context.Context, jobName string, limit int) ([]*domain.JobExecution, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.repo.GetByJobName(ctx, jobName, limit)
}

func (s *jobService) ListJobs(ctx context.Context) ([]JobInfo, error) {
	registeredJobs := s.scheduler.ListJobs()
	result := make([]JobInfo, len(registeredJobs))

	for i, job := range registeredJobs {
		info := JobInfo{
			Name:           job.Name(),
			Schedule:       job.Schedule(),
			TimeoutSeconds: int(job.Timeout().Seconds()),
		}

		// Get last run info
		lastExec, err := s.repo.GetLatestByJobName(ctx, job.Name())
		if err == nil && lastExec != nil {
			info.LastRunStatus = string(lastExec.Status)
			info.LastRunAt = lastExec.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		result[i] = info
	}

	return result, nil
}

func (s *jobService) GetRunningJobs(ctx context.Context) ([]*domain.JobExecution, error) {
	return s.repo.GetRunning(ctx)
}
