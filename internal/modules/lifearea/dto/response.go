package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/domain"
)

type LifeAreaResponse struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	Name         string    `json:"name"`
	Icon         string    `json:"icon,omitempty"`
	Color        string    `json:"color,omitempty"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
}

func ToLifeAreaResponse(la *domain.LifeArea) *LifeAreaResponse {
	return &LifeAreaResponse{
		ID:           la.ID,
		UserID:       la.UserID,
		Name:         la.Name,
		Icon:         la.Icon,
		Color:        la.Color,
		DisplayOrder: la.DisplayOrder,
		CreatedAt:    la.CreatedAt,
	}
}

func ToLifeAreaResponseList(areas []*domain.LifeArea) []*LifeAreaResponse {
	result := make([]*LifeAreaResponse, len(areas))
	for i, la := range areas {
		result[i] = ToLifeAreaResponse(la)
	}
	return result
}
