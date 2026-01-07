package jobs

import (
	"context"
	"time"
)

// JobStatus represents the execution status of a job
type JobStatus string

const (
	JobStatusPending   JobStatus = "pending"
	JobStatusRunning   JobStatus = "running"
	JobStatusCompleted JobStatus = "completed"
	JobStatusFailed    JobStatus = "failed"
)

// RetryPolicy configures retry behavior for jobs
type RetryPolicy struct {
	MaxRetries int           // Maximum number of retry attempts
	Delay      time.Duration // Initial delay before first retry
	Backoff    time.Duration // Additional delay for each subsequent retry
}

// DefaultRetryPolicy returns sensible defaults for retry behavior
func DefaultRetryPolicy() *RetryPolicy {
	return &RetryPolicy{
		MaxRetries: 3,
		Delay:      30 * time.Second,
		Backoff:    1 * time.Minute,
	}
}

// Job is the interface that all background jobs must implement
type Job interface {
	// Name returns unique job identifier (e.g., "calendar_sync", "streak_calculation")
	Name() string

	// Execute runs the job with the given context
	// Context may be cancelled if timeout is reached
	Execute(ctx context.Context) error

	// Schedule returns cron expression for automatic scheduling
	// Return empty string for manual-only jobs
	// Examples:
	//   "*/15 * * * *" - every 15 minutes
	//   "0 * * * *"    - every hour at :00
	//   "0 0 * * *"    - daily at midnight
	//   "0 0 * * 0"    - weekly on Sunday at midnight
	Schedule() string

	// Timeout returns maximum execution duration
	// Job will be cancelled if it exceeds this duration
	Timeout() time.Duration

	// RetryPolicy returns retry configuration for failed jobs
	// Return nil to disable retries
	RetryPolicy() *RetryPolicy
}

// JobResult contains the result of a job execution
type JobResult struct {
	JobName     string
	Status      JobStatus
	StartedAt   time.Time
	CompletedAt time.Time
	Duration    time.Duration
	Error       error
	Result      interface{}
}

// JobProgressEvent represents a progress update from a running job
type JobProgressEvent struct {
	JobName   string
	Progress  float64 // 0-100
	Message   string
	UpdatedAt time.Time
}

// JobEventEmitter is used by jobs to emit WebSocket events
type JobEventEmitter interface {
	EmitJobStarted(ctx context.Context, jobName string)
	EmitJobProgress(ctx context.Context, jobName string, progress float64, message string)
	EmitJobCompleted(ctx context.Context, jobName string, result interface{})
	EmitJobFailed(ctx context.Context, jobName string, err error)
}

// BaseJob provides common functionality for all jobs
type BaseJob struct {
	name     string
	schedule string
	timeout  time.Duration
	retry    *RetryPolicy
}

// NewBaseJob creates a new base job with the given configuration
func NewBaseJob(name, schedule string, timeout time.Duration, retry *RetryPolicy) BaseJob {
	if retry == nil {
		retry = DefaultRetryPolicy()
	}
	return BaseJob{
		name:     name,
		schedule: schedule,
		timeout:  timeout,
		retry:    retry,
	}
}

func (b BaseJob) Name() string              { return b.name }
func (b BaseJob) Schedule() string          { return b.schedule }
func (b BaseJob) Timeout() time.Duration    { return b.timeout }
func (b BaseJob) RetryPolicy() *RetryPolicy { return b.retry }
