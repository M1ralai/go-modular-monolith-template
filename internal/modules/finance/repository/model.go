package repository

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/domain"
)

type TransactionModel struct {
	ID          int       `db:"id"`
	UserID      int       `db:"user_id"`
	Amount      float64   `db:"amount"`
	Type        string    `db:"type"`
	Category    *string   `db:"category"`
	Description *string   `db:"description"`
	Date        time.Time `db:"date"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (m *TransactionModel) ToDomain() *domain.Transaction {
	if m == nil {
		return nil
	}
	cat, desc := "", ""
	if m.Category != nil {
		cat = *m.Category
	}
	if m.Description != nil {
		desc = *m.Description
	}
	return &domain.Transaction{ID: m.ID, UserID: m.UserID, Amount: m.Amount, Type: m.Type, Category: cat, Description: desc, Date: m.Date, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt}
}

func FromDomain(t *domain.Transaction) *TransactionModel {
	if t == nil {
		return nil
	}
	var cat, desc *string
	if t.Category != "" {
		cat = &t.Category
	}
	if t.Description != "" {
		desc = &t.Description
	}
	return &TransactionModel{ID: t.ID, UserID: t.UserID, Amount: t.Amount, Type: t.Type, Category: cat, Description: desc, Date: t.Date, CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}
