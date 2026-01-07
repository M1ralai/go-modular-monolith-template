package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/finance/repository"
)

type financeService struct {
	repo   repository.TransactionRepository
	logger *logger.ZapLogger
}

func NewFinanceService(repo repository.TransactionRepository, logger *logger.ZapLogger) FinanceService {
	return &financeService{repo: repo, logger: logger}
}

func (s *financeService) Create(ctx context.Context, req *dto.CreateTransactionRequest, userID int) (*dto.TransactionResponse, error) {
	s.logger.Info("Creating transaction", map[string]interface{}{"user_id": userID, "amount": req.Amount, "type": req.Type, "action": "CREATE_TRANSACTION"})
	now := time.Now()
	tx := &domain.Transaction{UserID: userID, Amount: req.Amount, Type: req.Type, Category: req.Category, Description: req.Description, Date: req.Date, CreatedAt: now, UpdatedAt: now}
	created, err := s.repo.Create(ctx, tx)
	if err != nil {
		s.logger.Error("Failed to create transaction", err, map[string]interface{}{"user_id": userID, "action": "CREATE_TRANSACTION_FAILED"})
		return nil, err
	}
	s.logger.Info("Transaction created", map[string]interface{}{"user_id": userID, "transaction_id": created.ID, "action": "CREATE_TRANSACTION_SUCCESS"})
	return dto.ToTransactionResponse(created), nil
}

func (s *financeService) GetByID(ctx context.Context, id, userID int) (*dto.TransactionResponse, error) {
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transaction not found")
	}
	if tx.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return dto.ToTransactionResponse(tx), nil
}

func (s *financeService) GetAll(ctx context.Context, userID int) ([]*dto.TransactionResponse, error) {
	txs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToTransactionResponseList(txs), nil
}

func (s *financeService) GetSummary(ctx context.Context, userID int, start, end time.Time) (*dto.SummaryResponse, error) {
	income, expense, err := s.repo.GetSummary(ctx, userID, start, end)
	if err != nil {
		return nil, err
	}
	return &dto.SummaryResponse{Income: income, Expense: expense, Balance: income - expense}, nil
}

func (s *financeService) Update(ctx context.Context, id int, req *dto.UpdateTransactionRequest, userID int) (*dto.TransactionResponse, error) {
	s.logger.Info("Updating transaction", map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "UPDATE_TRANSACTION"})
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, errors.New("transaction not found")
	}
	if tx.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Amount != nil {
		tx.Amount = *req.Amount
	}
	if req.Type != nil {
		tx.Type = *req.Type
	}
	if req.Category != nil {
		tx.Category = *req.Category
	}
	if req.Description != nil {
		tx.Description = *req.Description
	}
	if req.Date != nil {
		tx.Date = *req.Date
	}
	tx.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, tx); err != nil {
		s.logger.Error("Failed to update transaction", err, map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "UPDATE_TRANSACTION_FAILED"})
		return nil, err
	}
	s.logger.Info("Transaction updated", map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "UPDATE_TRANSACTION_SUCCESS"})
	return dto.ToTransactionResponse(tx), nil
}

func (s *financeService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting transaction", map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "DELETE_TRANSACTION"})
	tx, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if tx == nil {
		return errors.New("transaction not found")
	}
	if tx.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete transaction", err, map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "DELETE_TRANSACTION_FAILED"})
		return err
	}
	s.logger.Info("Transaction deleted", map[string]interface{}{"user_id": userID, "transaction_id": id, "action": "DELETE_TRANSACTION_SUCCESS"})
	return nil
}
