package dto

import "time"

type CreateTaskRequest struct {
	ParentTaskID   *int       `json:"parent_task_id,omitempty"`
	Title          string     `json:"title" validate:"required,min=1,max=255"`
	Description    string     `json:"description,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedStart *time.Time `json:"estimated_start,omitempty"`
	EstimatedEnd   *time.Time `json:"estimated_end,omitempty"`
	Priority       string     `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
}

type UpdateTaskRequest struct {
	Title          *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description    *string    `json:"description,omitempty"`
	DueDate        *time.Time `json:"due_date,omitempty"`
	EstimatedStart *time.Time `json:"estimated_start,omitempty"`
	EstimatedEnd   *time.Time `json:"estimated_end,omitempty"`
	ActualStart    *time.Time `json:"actual_start,omitempty"`
	ActualEnd      *time.Time `json:"actual_end,omitempty"`
	Priority       *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	IsCompleted    *bool      `json:"is_completed,omitempty"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}
