package jobs

import (
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/robfig/cron/v3"
)

// Scheduler wraps robfig/cron to provide job scheduling
type Scheduler struct {
	cron     *cron.Cron
	pool     *WorkerPool
	logger   *logger.ZapLogger
	jobs     map[string]Job
	entryIDs map[string]cron.EntryID
}

// NewScheduler creates a new cron-based scheduler
func NewScheduler(pool *WorkerPool, logger *logger.ZapLogger) *Scheduler {
	return &Scheduler{
		cron:     cron.New(cron.WithSeconds()),
		pool:     pool,
		logger:   logger,
		jobs:     make(map[string]Job),
		entryIDs: make(map[string]cron.EntryID),
	}
}

// Register adds a job to the scheduler
// If the job has a schedule, it will be automatically executed
func (s *Scheduler) Register(job Job) error {
	s.jobs[job.Name()] = job

	schedule := job.Schedule()
	if schedule == "" {
		s.logger.Info("Job registered (manual only)", map[string]interface{}{
			"job":    job.Name(),
			"action": "JOB_REGISTERED_MANUAL",
		})
		return nil
	}

	// Add cron job
	entryID, err := s.cron.AddFunc(schedule, func() {
		s.pool.SubmitAsync(job)
	})
	if err != nil {
		s.logger.Error("Failed to schedule job", err, map[string]interface{}{
			"job":      job.Name(),
			"schedule": schedule,
			"action":   "JOB_SCHEDULE_FAILED",
		})
		return err
	}

	s.entryIDs[job.Name()] = entryID

	s.logger.Info("Job registered", map[string]interface{}{
		"job":      job.Name(),
		"schedule": schedule,
		"action":   "JOB_REGISTERED",
	})

	return nil
}

// Start begins the scheduler
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("Scheduler started", map[string]interface{}{
		"jobs_count": len(s.jobs),
		"action":     "SCHEDULER_STARTED",
	})
}

// Stop gracefully shuts down the scheduler
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("Scheduler stopped", map[string]interface{}{
		"action": "SCHEDULER_STOPPED",
	})
}

// TriggerJob manually triggers a job by name
func (s *Scheduler) TriggerJob(jobName string) *JobResult {
	job, exists := s.jobs[jobName]
	if !exists {
		return &JobResult{
			JobName: jobName,
			Status:  JobStatusFailed,
			Error:   ErrJobNotFound,
		}
	}

	return s.pool.Submit(job)
}

// TriggerJobAsync manually triggers a job without waiting
func (s *Scheduler) TriggerJobAsync(jobName string) error {
	job, exists := s.jobs[jobName]
	if !exists {
		return ErrJobNotFound
	}

	return s.pool.SubmitAsync(job)
}

// GetJob returns a job by name
func (s *Scheduler) GetJob(jobName string) (Job, bool) {
	job, exists := s.jobs[jobName]
	return job, exists
}

// ListJobs returns all registered jobs
func (s *Scheduler) ListJobs() []Job {
	jobs := make([]Job, 0, len(s.jobs))
	for _, job := range s.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// GetEntries returns all cron entries for monitoring
func (s *Scheduler) GetEntries() []cron.Entry {
	return s.cron.Entries()
}
