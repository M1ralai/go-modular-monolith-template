package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
)

type TaskResponse struct {
	ID                   int        `json:"id"`
	UserID               int        `json:"user_id"`
	ParentTaskID         *int       `json:"parent_task_id,omitempty"`
	Title                string     `json:"title"`
	Description          string     `json:"description,omitempty"`
	DueDate              *time.Time `json:"due_date,omitempty"`
	EstimatedStart       *time.Time `json:"estimated_start,omitempty"`
	EstimatedEnd         *time.Time `json:"estimated_end,omitempty"`
	ActualStart          *time.Time `json:"actual_start,omitempty"`
	ActualEnd            *time.Time `json:"actual_end,omitempty"`
	Priority             string     `json:"priority"`
	IsCompleted          bool       `json:"is_completed"`
	CompletedAt          *time.Time `json:"completed_at,omitempty"`
	ProgressPercentage   float64    `json:"progress_percentage"`
	CompletedSubtasks    int        `json:"completed_subtasks"`
	TotalSubtasks        int        `json:"total_subtasks"`
	EstimatedDurationMin int        `json:"estimated_duration_min,omitempty"`
	ActualDurationMin    int        `json:"actual_duration_min,omitempty"`
	IsOverdue            bool       `json:"is_overdue"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

func ToTaskResponse(t *domain.Task, totalSubtasks, completedSubtasks int) *TaskResponse {
	if t == nil {
		return nil
	}

	progress := t.ProgressPercentage
	if totalSubtasks > 0 {
		progress = float64(completedSubtasks) / float64(totalSubtasks) * 100
	} else if t.IsCompleted {
		progress = 100
	}

	return &TaskResponse{
		ID:                   t.ID,
		UserID:               t.UserID,
		ParentTaskID:         t.ParentTaskID,
		Title:                t.Title,
		Description:          t.Description,
		DueDate:              t.DueDate,
		EstimatedStart:       t.EstimatedStart,
		EstimatedEnd:         t.EstimatedEnd,
		ActualStart:          t.ActualStart,
		ActualEnd:            t.ActualEnd,
		Priority:             t.Priority,
		IsCompleted:          t.IsCompleted,
		CompletedAt:          t.CompletedAt,
		ProgressPercentage:   progress,
		CompletedSubtasks:    completedSubtasks,
		TotalSubtasks:        totalSubtasks,
		EstimatedDurationMin: t.CalculateEstimatedDurationMinutes(),
		ActualDurationMin:    t.CalculateActualDurationMinutes(),
		IsOverdue:            t.IsOverdue(),
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
}

func ToSimpleTaskResponse(t *domain.Task) *TaskResponse {
	return ToTaskResponse(t, 0, 0)
}

func ToTaskResponseList(tasks []*domain.Task) []*TaskResponse {
	result := make([]*TaskResponse, len(tasks))
	for i, t := range tasks {
		result[i] = ToSimpleTaskResponse(t)
	}
	return result
}

type TaskStatsResponse struct {
	CompletedToday int `json:"completed_today"`
	DueToday       int `json:"due_today"`
	DueTomorrow    int `json:"due_tomorrow"`
	Overdue        int `json:"overdue"`
}
