package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
)

// WorkerPool manages a pool of workers that execute jobs concurrently
type WorkerPool struct {
	workers      int
	jobQueue     chan jobExecution
	quit         chan bool
	wg           sync.WaitGroup
	logger       *logger.ZapLogger
	eventEmitter JobEventEmitter
	mu           sync.RWMutex
	running      bool
	lock         *DistributedLock
}

type jobExecution struct {
	job      Job
	ctx      context.Context
	resultCh chan *JobResult
}

// NewWorkerPool creates a new worker pool with the specified number of workers
func NewWorkerPool(workers, queueSize int, logger *logger.ZapLogger, emitter JobEventEmitter, lock *DistributedLock) *WorkerPool {
	return &WorkerPool{
		workers:      workers,
		jobQueue:     make(chan jobExecution, queueSize),
		quit:         make(chan bool),
		logger:       logger,
		eventEmitter: emitter,
		lock:         lock,
	}
}

// Start launches all workers
func (p *WorkerPool) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.running {
		return
	}

	p.running = true
	for i := 0; i < p.workers; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}

	p.logger.Info("Worker pool started", map[string]interface{}{
		"workers":    p.workers,
		"queue_size": cap(p.jobQueue),
		"action":     "WORKER_POOL_STARTED",
	})
}

// Stop gracefully shuts down the worker pool
func (p *WorkerPool) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	close(p.quit)

	// Wait for workers with timeout
	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		p.logger.Info("Worker pool stopped gracefully", map[string]interface{}{
			"action": "WORKER_POOL_STOPPED",
		})
	case <-time.After(30 * time.Second):
		p.logger.Error("Worker pool shutdown timed out", nil, map[string]interface{}{
			"action": "WORKER_POOL_SHUTDOWN_TIMEOUT",
		})
	}
}

// Submit adds a job to the queue for execution
func (p *WorkerPool) Submit(job Job) *JobResult {
	p.mu.RLock()
	if !p.running {
		p.mu.RUnlock()
		return &JobResult{
			JobName:   job.Name(),
			Status:    JobStatusFailed,
			Error:     fmt.Errorf("worker pool not running"),
			StartedAt: time.Now(),
		}
	}
	p.mu.RUnlock()

	resultCh := make(chan *JobResult, 1)
	exec := jobExecution{
		job:      job,
		ctx:      context.Background(),
		resultCh: resultCh,
	}

	select {
	case p.jobQueue <- exec:
		p.logger.Info("Job submitted", map[string]interface{}{
			"job":    job.Name(),
			"action": "JOB_SUBMITTED",
		})
	default:
		p.logger.Error("Job queue full, dropping job", nil, map[string]interface{}{
			"job":    job.Name(),
			"action": "JOB_QUEUE_FULL",
		})
		return &JobResult{
			JobName:   job.Name(),
			Status:    JobStatusFailed,
			Error:     fmt.Errorf("job queue full"),
			StartedAt: time.Now(),
		}
	}

	return <-resultCh
}

// SubmitAsync adds a job to the queue without waiting for result
func (p *WorkerPool) SubmitAsync(job Job) error {
	p.mu.RLock()
	if !p.running {
		p.mu.RUnlock()
		return fmt.Errorf("worker pool not running")
	}
	p.mu.RUnlock()

	exec := jobExecution{
		job:      job,
		ctx:      context.Background(),
		resultCh: nil, // No result channel for async
	}

	select {
	case p.jobQueue <- exec:
		p.logger.Info("Job submitted async", map[string]interface{}{
			"job":    job.Name(),
			"action": "JOB_SUBMITTED_ASYNC",
		})
		return nil
	default:
		p.logger.Error("Job queue full, dropping job", nil, map[string]interface{}{
			"job":    job.Name(),
			"action": "JOB_QUEUE_FULL",
		})
		return fmt.Errorf("job queue full")
	}
}

// worker is the main worker goroutine
func (p *WorkerPool) worker(id int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.quit:
			p.logger.Info("Worker shutting down", map[string]interface{}{
				"worker_id": id,
				"action":    "WORKER_SHUTDOWN",
			})
			return
		case exec := <-p.jobQueue:
			result := p.executeJob(id, exec.job, exec.ctx)
			if exec.resultCh != nil {
				exec.resultCh <- result
			}
		}
	}
}

// executeJob runs a job with timeout, retry, and panic recovery
func (p *WorkerPool) executeJob(workerID int, job Job, ctx context.Context) *JobResult {
	result := &JobResult{
		JobName:   job.Name(),
		StartedAt: time.Now(),
	}

	// Create timeout context - increase to 30 seconds minimum for database operations
	timeout := job.Timeout()
	if timeout <= 0 || timeout < 30*time.Second {
		timeout = 30 * time.Second // Minimum 30 seconds for database operations
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Try to acquire distributed lock if job supports it
	var lockKey int64
	if lockableJob, ok := job.(LockableJob); ok && p.lock != nil {
		lockKey = lockableJob.LockKey()
		acquired, err := p.lock.TryLock(ctx, lockKey)
		if err != nil {
			p.logger.Error("Failed to acquire lock", err, map[string]interface{}{
				"job":     job.Name(),
				"lock_key": lockKey,
				"action":  "LOCK_ACQUIRE_FAILED",
			})
			result.Status = JobStatusFailed
			result.Error = fmt.Errorf("failed to acquire lock: %w", err)
			result.CompletedAt = time.Now()
			return result
		}
		if !acquired {
			p.logger.Info("Job already running, skipping", map[string]interface{}{
				"job":      job.Name(),
				"lock_key": lockKey,
				"action":   "JOB_SKIPPED_LOCK_HELD",
			})
			result.Status = JobStatusFailed
			result.Error = fmt.Errorf("job already running (lock held)")
			result.CompletedAt = time.Now()
			return result
		}
		defer func() {
			if unlockErr := p.lock.Unlock(context.Background(), lockKey); unlockErr != nil {
				p.logger.Error("Failed to release lock", unlockErr, map[string]interface{}{
					"job":      job.Name(),
					"lock_key": lockKey,
					"action":   "LOCK_RELEASE_FAILED",
				})
			}
		}()
	}

	// Emit started event
	if p.eventEmitter != nil {
		p.eventEmitter.EmitJobStarted(ctx, job.Name())
	}

	p.logger.Info("Job started", map[string]interface{}{
		"job":       job.Name(),
		"worker_id": workerID,
		"timeout":   timeout.String(),
		"action":    "JOB_STARTED",
	})

	// Execute with retry
	var err error
	retryPolicy := job.RetryPolicy()
	maxAttempts := 1
	if retryPolicy != nil {
		maxAttempts = retryPolicy.MaxRetries + 1
	}

	for attempt := 0; attempt < maxAttempts; attempt++ {
		if attempt > 0 {
			delay := retryPolicy.Delay + time.Duration(attempt-1)*retryPolicy.Backoff
			p.logger.Info("Job retry", map[string]interface{}{
				"job":     job.Name(),
				"attempt": attempt + 1,
				"delay":   delay.String(),
				"action":  "JOB_RETRY",
			})
			time.Sleep(delay)
		}

		err = p.safeExecute(ctx, job)
		if err == nil {
			break
		}

		p.logger.Error("Job execution failed", err, map[string]interface{}{
			"job":     job.Name(),
			"attempt": attempt + 1,
			"action":  "JOB_EXECUTION_FAILED",
		})
	}

	result.CompletedAt = time.Now()
	result.Duration = result.CompletedAt.Sub(result.StartedAt)

	if err != nil {
		result.Status = JobStatusFailed
		result.Error = err

		if p.eventEmitter != nil {
			p.eventEmitter.EmitJobFailed(ctx, job.Name(), err)
		}

		p.logger.Error("Job failed", err, map[string]interface{}{
			"job":      job.Name(),
			"duration": result.Duration.String(),
			"action":   "JOB_FAILED",
		})
	} else {
		result.Status = JobStatusCompleted

		if p.eventEmitter != nil {
			p.eventEmitter.EmitJobCompleted(ctx, job.Name(), nil)
		}

		p.logger.Info("Job completed", map[string]interface{}{
			"job":      job.Name(),
			"duration": result.Duration.String(),
			"action":   "JOB_COMPLETED",
		})
	}

	return result
}

// safeExecute runs job.Execute with panic recovery
func (p *WorkerPool) safeExecute(ctx context.Context, job Job) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic in job %s: %v", job.Name(), r)
			p.logger.Error("Job panicked", err, map[string]interface{}{
				"job":    job.Name(),
				"panic":  r,
				"action": "JOB_PANIC",
			})
		}
	}()

	return job.Execute(ctx)
}

// QueueSize returns the current number of jobs waiting in queue
func (p *WorkerPool) QueueSize() int {
	return len(p.jobQueue)
}

// IsRunning returns whether the worker pool is running
func (p *WorkerPool) IsRunning() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.running
}
