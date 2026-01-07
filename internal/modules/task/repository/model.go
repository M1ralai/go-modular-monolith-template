package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/task/domain"
)

type TaskModel struct {
	ID                 int        `db:"id"`
	UserID             int        `db:"user_id"`
	ParentTaskID       *int       `db:"parent_task_id"`
	Title              string     `db:"title"`
	Description        *string    `db:"description"`
	DueDate            *time.Time `db:"due_date"`
	EstimatedStart     *time.Time `db:"estimated_start"`
	EstimatedEnd       *time.Time `db:"estimated_end"`
	ActualStart        *time.Time `db:"actual_start"`
	ActualEnd          *time.Time `db:"actual_end"`
	Priority           string     `db:"priority"`
	IsCompleted        bool       `db:"is_completed"`
	CompletedAt        *time.Time `db:"completed_at"`
	ProgressPercentage float64    `db:"progress_percentage"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

func (m *TaskModel) ToDomain() *domain.Task {
	if m == nil {
		return nil
	}

	description := ""
	if m.Description != nil {
		description = *m.Description
	}

	return &domain.Task{
		ID:                 m.ID,
		UserID:             m.UserID,
		ParentTaskID:       m.ParentTaskID,
		Title:              m.Title,
		Description:        description,
		DueDate:            m.DueDate,
		EstimatedStart:     m.EstimatedStart,
		EstimatedEnd:       m.EstimatedEnd,
		ActualStart:        m.ActualStart,
		ActualEnd:          m.ActualEnd,
		Priority:           m.Priority,
		IsCompleted:        m.IsCompleted,
		CompletedAt:        m.CompletedAt,
		ProgressPercentage: m.ProgressPercentage,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}
}

func FromDomain(t *domain.Task) *TaskModel {
	if t == nil {
		return nil
	}

	var description *string
	if t.Description != "" {
		description = &t.Description
	}

	return &TaskModel{
		ID:                 t.ID,
		UserID:             t.UserID,
		ParentTaskID:       t.ParentTaskID,
		Title:              t.Title,
		Description:        description,
		DueDate:            t.DueDate,
		EstimatedStart:     t.EstimatedStart,
		EstimatedEnd:       t.EstimatedEnd,
		ActualStart:        t.ActualStart,
		ActualEnd:          t.ActualEnd,
		Priority:           t.Priority,
		IsCompleted:        t.IsCompleted,
		CompletedAt:        t.CompletedAt,
		ProgressPercentage: t.ProgressPercentage,
		CreatedAt:          t.CreatedAt,
		UpdatedAt:          t.UpdatedAt,
	}
}
