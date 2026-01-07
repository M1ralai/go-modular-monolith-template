package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/domain"
)

type TransactionResponse struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Amount      float64   `json:"amount"`
	Type        string    `json:"type"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SummaryResponse struct {
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Balance float64 `json:"balance"`
}

func ToTransactionResponse(t *domain.Transaction) *TransactionResponse {
	if t == nil {
		return nil
	}
	return &TransactionResponse{ID: t.ID, UserID: t.UserID, Amount: t.Amount, Type: t.Type, Category: t.Category, Description: t.Description, Date: t.Date, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}

func ToTransactionResponseList(txs []*domain.Transaction) []*TransactionResponse {
	result := make([]*TransactionResponse, len(txs))
	for i, t := range txs {
		result[i] = ToTransactionResponse(t)
	}
	return result
}
