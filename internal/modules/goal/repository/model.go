package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/goal/domain"
)

type GoalModel struct {
	ID          int        `db:"id"`
	UserID      int        `db:"user_id"`
	LifeAreaID  *int       `db:"life_area_id"`
	Title       string     `db:"title"`
	Description *string    `db:"description"`
	TargetDate  *time.Time `db:"target_date"`
	IsCompleted bool       `db:"is_completed"`
	CompletedAt *time.Time `db:"completed_at"`
	Priority    string     `db:"priority"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
}

func (m *GoalModel) ToDomain() *domain.Goal {
	if m == nil {
		return nil
	}
	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}
	return &domain.Goal{ID: m.ID, UserID: m.UserID, LifeAreaID: m.LifeAreaID, Title: m.Title, Description: desc, TargetDate: m.TargetDate, IsCompleted: m.IsCompleted, CompletedAt: m.CompletedAt, Priority: m.Priority, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func FromDomain(g *domain.Goal) *GoalModel {
	if g == nil {
		return nil
	}
	var desc *string
	if g.Description != "" {
		desc = &g.Description
	}
	return &GoalModel{ID: g.ID, UserID: g.UserID, LifeAreaID: g.LifeAreaID, Title: g.Title, Description: desc, TargetDate: g.TargetDate, IsCompleted: g.IsCompleted, CompletedAt: g.CompletedAt, Priority: g.Priority, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt}
}

type MilestoneModel struct {
	ID          int        `db:"id"`
	GoalID      int        `db:"goal_id"`
	Title       string     `db:"title"`
	IsCompleted bool       `db:"is_completed"`
	CompletedAt *time.Time `db:"completed_at"`
	CreatedAt   time.Time  `db:"created_at"`
}
