package jobs

import (
	"errors"
	"sync"
)

// Common errors
var (
	ErrJobNotFound      = errors.New("job not found")
	ErrJobAlreadyExists = errors.New("job already exists")
)

// Registry manages job registration and lookup
type Registry struct {
	jobs map[string]Job
	mu   sync.RWMutex
}

// NewRegistry creates a new job registry
func NewRegistry() *Registry {
	return &Registry{
		jobs: make(map[string]Job),
	}
}

// Register adds a job to the registry
func (r *Registry) Register(job Job) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[job.Name()]; exists {
		return ErrJobAlreadyExists
	}

	r.jobs[job.Name()] = job
	return nil
}

// Get returns a job by name
func (r *Registry) Get(name string) (Job, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	job, exists := r.jobs[name]
	if !exists {
		return nil, ErrJobNotFound
	}

	return job, nil
}

// List returns all registered jobs
func (r *Registry) List() []Job {
	r.mu.RLock()
	defer r.mu.RUnlock()

	jobs := make([]Job, 0, len(r.jobs))
	for _, job := range r.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}

// Unregister removes a job from the registry
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.jobs[name]; !exists {
		return ErrJobNotFound
	}

	delete(r.jobs, name)
	return nil
}

// Count returns the number of registered jobs
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.jobs)
}
