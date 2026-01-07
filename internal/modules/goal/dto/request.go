package dto

import "time"

type CreateGoalRequest struct {
	LifeAreaID  *int       `json:"life_area_id,omitempty"`
	Title       string     `json:"title" validate:"required,min=1,max=255"`
	Description string     `json:"description,omitempty"`
	TargetDate  *time.Time `json:"target_date,omitempty"`
	Priority    string     `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
}

type UpdateGoalRequest struct {
	LifeAreaID  *int       `json:"life_area_id,omitempty"`
	Title       *string    `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string    `json:"description,omitempty"`
	TargetDate  *time.Time `json:"target_date,omitempty"`
	Priority    *string    `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	IsCompleted *bool      `json:"is_completed,omitempty"`
}
