package domain

import "time"

type Transaction struct {
	ID          int
	UserID      int
	Amount      float64
	Type        string
	Category    string
	Description string
	Date        time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *Transaction) IsExpense() bool { return t.Type == "expense" }
func (t *Transaction) IsIncome() bool  { return t.Type == "income" }
