package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/lifearea/domain"
)

type LifeAreaModel struct {
	ID           int       `db:"id"`
	UserID       int       `db:"user_id"`
	Name         string    `db:"name"`
	Icon         *string   `db:"icon"`
	Color        *string   `db:"color"`
	DisplayOrder int       `db:"display_order"`
	CreatedAt    time.Time `db:"created_at"`
}

func (m *LifeAreaModel) ToDomain() *domain.LifeArea {
	icon := ""
	if m.Icon != nil {
		icon = *m.Icon
	}
	color := ""
	if m.Color != nil {
		color = *m.Color
	}

	return &domain.LifeArea{
		ID:           m.ID,
		UserID:       m.UserID,
		Name:         m.Name,
		Icon:         icon,
		Color:        color,
		DisplayOrder: m.DisplayOrder,
		CreatedAt:    m.CreatedAt,
	}
}

func FromDomain(la *domain.LifeArea) *LifeAreaModel {
	var icon, color *string
	if la.Icon != "" {
		icon = &la.Icon
	}
	if la.Color != "" {
		color = &la.Color
	}

	return &LifeAreaModel{
		ID:           la.ID,
		UserID:       la.UserID,
		Name:         la.Name,
		Icon:         icon,
		Color:        color,
		DisplayOrder: la.DisplayOrder,
		CreatedAt:    la.CreatedAt,
	}
}
