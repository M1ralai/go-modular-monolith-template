package domain

import "time"

type Goal struct {
	ID          int
	UserID      int
	LifeAreaID  *int
	Title       string
	Description string
	TargetDate  *time.Time
	IsCompleted bool
	CompletedAt *time.Time
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (g *Goal) MarkCompleted() {
	now := time.Now()
	g.IsCompleted = true
	g.CompletedAt = &now
	g.UpdatedAt = now
}
