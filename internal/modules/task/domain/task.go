package domain

import "time"

type Task struct {
	ID                 int
	UserID             int
	ParentTaskID       *int
	Title              string
	Description        string
	DueDate            *time.Time
	EstimatedStart     *time.Time
	EstimatedEnd       *time.Time
	ActualStart        *time.Time
	ActualEnd          *time.Time
	Priority           string
	IsCompleted        bool
	CompletedAt        *time.Time
	ProgressPercentage float64
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.IsCompleted {
		return false
	}
	return time.Now().After(*t.DueDate)
}

func (t *Task) CalculateEstimatedDurationMinutes() int {
	if t.EstimatedStart == nil || t.EstimatedEnd == nil {
		return 0
	}
	return int(t.EstimatedEnd.Sub(*t.EstimatedStart).Minutes())
}

func (t *Task) CalculateActualDurationMinutes() int {
	if t.ActualStart == nil || t.ActualEnd == nil {
		return 0
	}
	return int(t.ActualEnd.Sub(*t.ActualStart).Minutes())
}

func (t *Task) IsSubtask() bool {
	return t.ParentTaskID != nil
}

func (t *Task) MarkCompleted() {
	now := time.Now()
	t.IsCompleted = true
	t.CompletedAt = &now
	t.ProgressPercentage = 100
	t.UpdatedAt = now
}
