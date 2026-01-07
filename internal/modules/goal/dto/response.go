package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/domain"
)

type GoalResponse struct {
	ID                  int        `json:"id"`
	UserID              int        `json:"user_id"`
	LifeAreaID          *int       `json:"life_area_id,omitempty"`
	Title               string     `json:"title"`
	Description         string     `json:"description,omitempty"`
	TargetDate          *time.Time `json:"target_date,omitempty"`
	IsCompleted         bool       `json:"is_completed"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"`
	Priority            string     `json:"priority"`
	ProgressPercentage  float64    `json:"progress_percentage"`
	TotalMilestones     int        `json:"total_milestones"`
	CompletedMilestones int        `json:"completed_milestones"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}

func ToGoalResponse(g *domain.Goal, total, completed int) *GoalResponse {
	if g == nil {
		return nil
	}
	progress := 0.0
	if total > 0 {
		progress = float64(completed) / float64(total) * 100
	}
	if g.IsCompleted {
		progress = 100
	}
	return &GoalResponse{ID: g.ID, UserID: g.UserID, LifeAreaID: g.LifeAreaID, Title: g.Title, Description: g.Description, TargetDate: g.TargetDate, IsCompleted: g.IsCompleted, CompletedAt: g.CompletedAt, Priority: g.Priority, ProgressPercentage: progress, TotalMilestones: total, CompletedMilestones: completed, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt}
}
