package dto

type CreateLifeAreaRequest struct {
	Name         string `json:"name" validate:"required,min=1,max=100"`
	Icon         string `json:"icon,omitempty" validate:"omitempty,max=50"`
	Color        string `json:"color,omitempty" validate:"omitempty,max=7"`
	DisplayOrder int    `json:"display_order,omitempty"`
}

type UpdateLifeAreaRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	Icon         *string `json:"icon,omitempty" validate:"omitempty,max=50"`
	Color        *string `json:"color,omitempty" validate:"omitempty,max=7"`
	DisplayOrder *int    `json:"display_order,omitempty"`
}
