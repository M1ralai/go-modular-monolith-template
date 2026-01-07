package dto

type CreatePersonRequest struct {
	Name         string   `json:"name" validate:"required,min=1,max=255"`
	Email        string   `json:"email,omitempty" validate:"omitempty,email"`
	Phone        string   `json:"phone,omitempty"`
	Company      string   `json:"company,omitempty"`
	Relationship string   `json:"relationship,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Notes        string   `json:"notes,omitempty"`
}

type UpdatePersonRequest struct {
	Name         *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Email        *string  `json:"email,omitempty" validate:"omitempty,email"`
	Phone        *string  `json:"phone,omitempty"`
	Company      *string  `json:"company,omitempty"`
	Relationship *string  `json:"relationship,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	Notes        *string  `json:"notes,omitempty"`
}
