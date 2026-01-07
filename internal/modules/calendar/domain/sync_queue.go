package domain

import "time"

type SyncQueue struct {
	ID           int
	UserID       int
	EventID      *int
	Provider     string
	Action       string
	Status       string
	RetryCount   int
	ErrorMessage string
	CreatedAt    time.Time
	SyncedAt     *time.Time
}

func (s *SyncQueue) IsPending() bool { return s.Status == "pending" }
func (s *SyncQueue) IsFailed() bool  { return s.Status == "failed" }
func (s *SyncQueue) CanRetry() bool  { return s.RetryCount < 5 }
