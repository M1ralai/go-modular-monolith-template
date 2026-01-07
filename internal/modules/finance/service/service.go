package service

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/dto"
)

type FinanceService interface {
	Create(ctx context.Context, req *dto.CreateTransactionRequest, userID int) (*dto.TransactionResponse, error)
	GetByID(ctx context.Context, id, userID int) (*dto.TransactionResponse, error)
	GetAll(ctx context.Context, userID int) ([]*dto.TransactionResponse, error)
	GetSummary(ctx context.Context, userID int, start, end time.Time) (*dto.SummaryResponse, error)
	Update(ctx context.Context, id int, req *dto.UpdateTransactionRequest, userID int) (*dto.TransactionResponse, error)
	Delete(ctx context.Context, id, userID int) error
}
