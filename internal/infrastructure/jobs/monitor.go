package jobs

import (
	"sync"
	"time"
)

// JobMetrics tracks metrics for job execution
type JobMetrics struct {
	TotalExecutions  int64
	SuccessfulRuns   int64
	FailedRuns       int64
	TotalDuration    time.Duration
	AverageDuration  time.Duration
	LastRunAt        time.Time
	LastRunStatus    JobStatus
	LastRunDuration  time.Duration
	ConsecutiveFails int
}

// Monitor tracks job execution metrics
type Monitor struct {
	metrics map[string]*JobMetrics
	mu      sync.RWMutex
}

// NewMonitor creates a new job monitor
func NewMonitor() *Monitor {
	return &Monitor{
		metrics: make(map[string]*JobMetrics),
	}
}

// RecordExecution records a job execution result
func (m *Monitor) RecordExecution(result *JobResult) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, exists := m.metrics[result.JobName]
	if !exists {
		metrics = &JobMetrics{}
		m.metrics[result.JobName] = metrics
	}

	metrics.TotalExecutions++
	metrics.TotalDuration += result.Duration
	metrics.AverageDuration = time.Duration(int64(metrics.TotalDuration) / metrics.TotalExecutions)
	metrics.LastRunAt = result.StartedAt
	metrics.LastRunStatus = result.Status
	metrics.LastRunDuration = result.Duration

	if result.Status == JobStatusCompleted {
		metrics.SuccessfulRuns++
		metrics.ConsecutiveFails = 0
	} else if result.Status == JobStatusFailed {
		metrics.FailedRuns++
		metrics.ConsecutiveFails++
	}
}

// GetMetrics returns metrics for a specific job
func (m *Monitor) GetMetrics(jobName string) *JobMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.metrics[jobName]
	if !exists {
		return nil
	}

	// Return a copy
	return &JobMetrics{
		TotalExecutions:  metrics.TotalExecutions,
		SuccessfulRuns:   metrics.SuccessfulRuns,
		FailedRuns:       metrics.FailedRuns,
		TotalDuration:    metrics.TotalDuration,
		AverageDuration:  metrics.AverageDuration,
		LastRunAt:        metrics.LastRunAt,
		LastRunStatus:    metrics.LastRunStatus,
		LastRunDuration:  metrics.LastRunDuration,
		ConsecutiveFails: metrics.ConsecutiveFails,
	}
}

// GetAllMetrics returns metrics for all jobs
func (m *Monitor) GetAllMetrics() map[string]*JobMetrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]*JobMetrics, len(m.metrics))
	for name, metrics := range m.metrics {
		result[name] = &JobMetrics{
			TotalExecutions:  metrics.TotalExecutions,
			SuccessfulRuns:   metrics.SuccessfulRuns,
			FailedRuns:       metrics.FailedRuns,
			TotalDuration:    metrics.TotalDuration,
			AverageDuration:  metrics.AverageDuration,
			LastRunAt:        metrics.LastRunAt,
			LastRunStatus:    metrics.LastRunStatus,
			LastRunDuration:  metrics.LastRunDuration,
			ConsecutiveFails: metrics.ConsecutiveFails,
		}
	}
	return result
}

// GetJobsWithConsecutiveFailures returns jobs that have failed consecutively
func (m *Monitor) GetJobsWithConsecutiveFailures(threshold int) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var jobs []string
	for name, metrics := range m.metrics {
		if metrics.ConsecutiveFails >= threshold {
			jobs = append(jobs, name)
		}
	}
	return jobs
}

// GetSuccessRate returns the success rate for a job (0-100)
func (m *Monitor) GetSuccessRate(jobName string) float64 {
	m.mu.RLock()
	defer m.mu.RUnlock()

	metrics, exists := m.metrics[jobName]
	if !exists || metrics.TotalExecutions == 0 {
		return 0
	}

	return float64(metrics.SuccessfulRuns) / float64(metrics.TotalExecutions) * 100
}
