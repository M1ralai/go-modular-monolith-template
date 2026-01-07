package service

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/journal/repository"
)

type journalService struct {
	repo   repository.JournalRepository
	logger *logger.ZapLogger
}

func NewJournalService(repo repository.JournalRepository, logger *logger.ZapLogger) JournalService {
	return &journalService{repo: repo, logger: logger}
}

func (s *journalService) Create(ctx context.Context, req *dto.CreateJournalRequest, userID int) (*dto.JournalResponse, error) {
	s.logger.Info("Creating journal entry", map[string]interface{}{"user_id": userID, "date": req.EntryDate, "action": "CREATE_JOURNAL"})
	now := time.Now()
	entry := &domain.JournalEntry{UserID: userID, EntryDate: req.EntryDate, Content: req.Content, Mood: req.Mood, EnergyLevel: req.EnergyLevel, CreatedAt: now, UpdatedAt: now}
	created, err := s.repo.Create(ctx, entry)
	if err != nil {
		s.logger.Error("Failed to create journal", err, map[string]interface{}{"user_id": userID, "action": "CREATE_JOURNAL_FAILED"})
		return nil, err
	}
	s.logger.Info("Journal created", map[string]interface{}{"user_id": userID, "journal_id": created.ID, "action": "CREATE_JOURNAL_SUCCESS"})
	return dto.ToJournalResponse(created), nil
}

func (s *journalService) GetByID(ctx context.Context, id, userID int) (*dto.JournalResponse, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("journal entry not found")
	}
	if entry.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	return dto.ToJournalResponse(entry), nil
}

func (s *journalService) GetAll(ctx context.Context, userID int) ([]*dto.JournalResponse, error) {
	entries, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return dto.ToJournalResponseList(entries), nil
}

func (s *journalService) Update(ctx context.Context, id int, req *dto.UpdateJournalRequest, userID int) (*dto.JournalResponse, error) {
	s.logger.Info("Updating journal", map[string]interface{}{"user_id": userID, "journal_id": id, "action": "UPDATE_JOURNAL"})
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, errors.New("journal entry not found")
	}
	if entry.UserID != userID {
		return nil, errors.New("unauthorized")
	}
	if req.Content != nil {
		entry.Content = *req.Content
	}
	if req.Mood != nil {
		entry.Mood = *req.Mood
	}
	if req.EnergyLevel != nil {
		entry.EnergyLevel = *req.EnergyLevel
	}
	entry.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, entry); err != nil {
		s.logger.Error("Failed to update journal", err, map[string]interface{}{"user_id": userID, "journal_id": id, "action": "UPDATE_JOURNAL_FAILED"})
		return nil, err
	}
	s.logger.Info("Journal updated", map[string]interface{}{"user_id": userID, "journal_id": id, "action": "UPDATE_JOURNAL_SUCCESS"})
	return dto.ToJournalResponse(entry), nil
}

func (s *journalService) Delete(ctx context.Context, id, userID int) error {
	s.logger.Info("Deleting journal", map[string]interface{}{"user_id": userID, "journal_id": id, "action": "DELETE_JOURNAL"})
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if entry == nil {
		return errors.New("journal entry not found")
	}
	if entry.UserID != userID {
		return errors.New("unauthorized")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete journal", err, map[string]interface{}{"user_id": userID, "journal_id": id, "action": "DELETE_JOURNAL_FAILED"})
		return err
	}
	s.logger.Info("Journal deleted", map[string]interface{}{"user_id": userID, "journal_id": id, "action": "DELETE_JOURNAL_SUCCESS"})
	return nil
}
