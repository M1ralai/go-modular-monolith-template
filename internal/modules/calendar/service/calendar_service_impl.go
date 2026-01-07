package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/domain"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/repository"
)

type calendarService struct {
	integrationRepo repository.CalendarIntegrationRepository
	syncQueueRepo   repository.SyncQueueRepository
	logger          *logger.ZapLogger
}

func NewCalendarService(
	integrationRepo repository.CalendarIntegrationRepository,
	syncQueueRepo repository.SyncQueueRepository,
	logger *logger.ZapLogger,
) CalendarService {
	return &calendarService{
		integrationRepo: integrationRepo,
		syncQueueRepo:   syncQueueRepo,
		logger:          logger,
	}
}

func (s *calendarService) GetGoogleAuthURL(ctx context.Context, userID int) (string, error) {
	s.logger.Info("Generating Google auth URL", map[string]interface{}{"user_id": userID, "action": "GET_GOOGLE_AUTH_URL"})

	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	redirectURI := os.Getenv("GOOGLE_REDIRECT_URI")
	if clientID == "" || redirectURI == "" {
		return "", errors.New("Google OAuth not configured")
	}

	scope := "https://www.googleapis.com/auth/calendar"
	authURL := "https://accounts.google.com/o/oauth2/v2/auth?client_id=" + clientID +
		"&redirect_uri=" + redirectURI +
		"&response_type=code&scope=" + scope +
		"&access_type=offline&prompt=consent&state=" + string(rune(userID))

	s.logger.Info("Google auth URL generated", map[string]interface{}{"user_id": userID, "action": "GET_GOOGLE_AUTH_URL_SUCCESS"})
	return authURL, nil
}

func (s *calendarService) HandleGoogleCallback(ctx context.Context, userID int, code string) (*dto.IntegrationResponse, error) {
	s.logger.Info("Handling Google OAuth callback", map[string]interface{}{"user_id": userID, "action": "GOOGLE_OAUTH_CALLBACK"})

	// In production: Exchange code for tokens using Google OAuth API
	// For now, simulate token exchange
	now := time.Now()
	expiresAt := now.Add(1 * time.Hour)

	existing, _ := s.integrationRepo.GetByUserAndProvider(ctx, userID, "google")
	if existing != nil {
		existing.AccessToken = "simulated_access_token_" + code
		existing.RefreshToken = "simulated_refresh_token"
		existing.ExpiresAt = &expiresAt
		existing.IsActive = true
		existing.UpdatedAt = now
		if err := s.integrationRepo.Update(ctx, existing); err != nil {
			s.logger.Error("Failed to update Google integration", err, map[string]interface{}{"user_id": userID, "action": "GOOGLE_OAUTH_CALLBACK_FAILED"})
			return nil, err
		}
		s.logger.Info("Google integration updated", map[string]interface{}{"user_id": userID, "integration_id": existing.ID, "action": "GOOGLE_OAUTH_CALLBACK_SUCCESS"})
		return dto.ToIntegrationResponse(existing), nil
	}

	integration := &domain.CalendarIntegration{
		UserID:       userID,
		Provider:     "google",
		AccessToken:  "simulated_access_token_" + code,
		RefreshToken: "simulated_refresh_token",
		ExpiresAt:    &expiresAt,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	created, err := s.integrationRepo.Create(ctx, integration)
	if err != nil {
		s.logger.Error("Failed to create Google integration", err, map[string]interface{}{"user_id": userID, "action": "GOOGLE_OAUTH_CALLBACK_FAILED"})
		return nil, err
	}

	s.logger.Info("Google integration created", map[string]interface{}{"user_id": userID, "integration_id": created.ID, "action": "GOOGLE_OAUTH_CALLBACK_SUCCESS"})
	return dto.ToIntegrationResponse(created), nil
}

func (s *calendarService) DisconnectGoogle(ctx context.Context, userID int) error {
	s.logger.Info("Disconnecting Google Calendar", map[string]interface{}{"user_id": userID, "action": "DISCONNECT_GOOGLE"})

	integration, err := s.integrationRepo.GetByUserAndProvider(ctx, userID, "google")
	if err != nil {
		return err
	}
	if integration == nil {
		return errors.New("no Google integration found")
	}

	integration.IsActive = false
	if err := s.integrationRepo.Update(ctx, integration); err != nil {
		s.logger.Error("Failed to disconnect Google", err, map[string]interface{}{"user_id": userID, "action": "DISCONNECT_GOOGLE_FAILED"})
		return err
	}

	s.logger.Info("Google Calendar disconnected", map[string]interface{}{"user_id": userID, "action": "DISCONNECT_GOOGLE_SUCCESS"})
	return nil
}

func (s *calendarService) SyncGoogle(ctx context.Context, userID int) error {
	s.logger.Info("Syncing with Google Calendar", map[string]interface{}{"user_id": userID, "action": "SYNC_GOOGLE"})

	integration, err := s.integrationRepo.GetByUserAndProvider(ctx, userID, "google")
	if err != nil {
		return err
	}
	if integration == nil || !integration.IsActive {
		return errors.New("Google Calendar not connected")
	}

	// In production: Call Google Calendar API to sync events
	// For now, simulate sync
	now := time.Now()
	integration.LastSyncAt = &now
	if err := s.integrationRepo.Update(ctx, integration); err != nil {
		s.logger.Error("Failed to sync Google", err, map[string]interface{}{"user_id": userID, "action": "SYNC_GOOGLE_FAILED"})
		return err
	}

	s.logger.Info("Google Calendar synced", map[string]interface{}{"user_id": userID, "action": "SYNC_GOOGLE_SUCCESS"})
	return nil
}

func (s *calendarService) GetSyncStatus(ctx context.Context, userID int) (*dto.SyncStatusResponse, error) {
	integrations, err := s.integrationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	status := &dto.SyncStatusResponse{Integrations: make([]*dto.IntegrationResponse, len(integrations))}
	for i, integration := range integrations {
		status.Integrations[i] = dto.ToIntegrationResponse(integration)
	}

	return status, nil
}

func (s *calendarService) GetIntegrations(ctx context.Context, userID int) ([]*dto.IntegrationResponse, error) {
	integrations, err := s.integrationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]*dto.IntegrationResponse, len(integrations))
	for i, integration := range integrations {
		result[i] = dto.ToIntegrationResponse(integration)
	}
	return result, nil
}

func (s *calendarService) QueueSync(ctx context.Context, userID int, eventID int, action string) error {
	s.logger.Info("Queueing sync operation", map[string]interface{}{"user_id": userID, "event_id": eventID, "action_type": action, "action": "QUEUE_SYNC"})

	integrations, err := s.integrationRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	for _, integration := range integrations {
		if !integration.IsActive {
			continue
		}

		queueItem := &domain.SyncQueue{
			UserID:   userID,
			EventID:  &eventID,
			Provider: integration.Provider,
			Action:   action,
			Status:   "pending",
		}
		if _, err := s.syncQueueRepo.Create(ctx, queueItem); err != nil {
			s.logger.Error("Failed to queue sync", err, map[string]interface{}{"user_id": userID, "event_id": eventID, "provider": integration.Provider, "action": "QUEUE_SYNC_FAILED"})
			return err
		}
	}

	s.logger.Info("Sync queued", map[string]interface{}{"user_id": userID, "event_id": eventID, "action": "QUEUE_SYNC_SUCCESS"})
	return nil
}

func (s *calendarService) ProcessSyncQueue(ctx context.Context, limit int) (int, error) {
	s.logger.Info("Processing sync queue", map[string]interface{}{"limit": limit, "action": "PROCESS_SYNC_QUEUE"})

	pending, err := s.syncQueueRepo.GetPending(ctx, limit)
	if err != nil {
		return 0, err
	}

	processed := 0
	for _, item := range pending {
		// In production: Call respective calendar API based on provider
		// For now, simulate processing
		if err := s.syncQueueRepo.UpdateStatus(ctx, item.ID, "completed", ""); err != nil {
			s.logger.Error("Failed to process sync item", err, map[string]interface{}{"sync_id": item.ID, "action": "PROCESS_SYNC_ITEM_FAILED"})
			s.syncQueueRepo.IncrementRetry(ctx, item.ID)
			continue
		}
		processed++
	}

	s.logger.Info("Sync queue processed", map[string]interface{}{"processed": processed, "action": "PROCESS_SYNC_QUEUE_SUCCESS"})
	return processed, nil
}
