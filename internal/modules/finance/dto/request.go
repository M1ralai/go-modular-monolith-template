package dto

import "time"

type CreateTransactionRequest struct {
	Amount      float64   `json:"amount" validate:"required"`
	Type        string    `json:"type" validate:"required,oneof=income expense"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date" validate:"required"`
}

type UpdateTransactionRequest struct {
	Amount      *float64   `json:"amount,omitempty"`
	Type        *string    `json:"type,omitempty" validate:"omitempty,oneof=income expense"`
	Category    *string    `json:"category,omitempty"`
	Description *string    `json:"description,omitempty"`
	Date        *time.Time `json:"date,omitempty"`
}
