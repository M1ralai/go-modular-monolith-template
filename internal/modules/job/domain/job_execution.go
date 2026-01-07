package domain

import (
	"encoding/json"
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

// JobExecution represents a single job execution record
type JobExecution struct {
	ID          int             `db:"id"`
	JobName     string          `db:"job_name"`
	Status      JobStatus       `db:"status"`
	StartedAt   time.Time       `db:"started_at"`
	CompletedAt *time.Time      `db:"completed_at"`
	Error       *string         `db:"error"`
	Result      json.RawMessage `db:"result"`
	DurationMs  *int            `db:"duration_ms"`
	CreatedAt   time.Time       `db:"created_at"`
}

// NewJobExecution creates a new job execution record
func NewJobExecution(jobName string) *JobExecution {
	return &JobExecution{
		JobName:   jobName,
		Status:    JobStatusPending,
		StartedAt: time.Now(),
		CreatedAt: time.Now(),
	}
}

// MarkRunning marks the job as running
func (j *JobExecution) MarkRunning() {
	j.Status = JobStatusRunning
	now := time.Now()
	j.StartedAt = now
}

// MarkCompleted marks the job as completed
func (j *JobExecution) MarkCompleted(result interface{}) {
	j.Status = JobStatusCompleted
	now := time.Now()
	j.CompletedAt = &now
	durationMs := int(now.Sub(j.StartedAt).Milliseconds())
	j.DurationMs = &durationMs

	if result != nil {
		if data, err := json.Marshal(result); err == nil {
			j.Result = data
		}
	}
}

// MarkFailed marks the job as failed
func (j *JobExecution) MarkFailed(err error) {
	j.Status = JobStatusFailed
	now := time.Now()
	j.CompletedAt = &now
	durationMs := int(now.Sub(j.StartedAt).Milliseconds())
	j.DurationMs = &durationMs

	if err != nil {
		errStr := err.Error()
		j.Error = &errStr
	}
}

// Duration returns the execution duration
func (j *JobExecution) Duration() time.Duration {
	if j.DurationMs != nil {
		return time.Duration(*j.DurationMs) * time.Millisecond
	}
	if j.CompletedAt != nil {
		return j.CompletedAt.Sub(j.StartedAt)
	}
	return time.Since(j.StartedAt)
}
