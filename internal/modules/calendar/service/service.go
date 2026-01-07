package service

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/calendar/dto"
)

type CalendarService interface {
	// OAuth Flow
	GetGoogleAuthURL(ctx context.Context, userID int) (string, error)
	HandleGoogleCallback(ctx context.Context, userID int, code string) (*dto.IntegrationResponse, error)
	DisconnectGoogle(ctx context.Context, userID int) error

	// Sync Operations
	SyncGoogle(ctx context.Context, userID int) error
	GetSyncStatus(ctx context.Context, userID int) (*dto.SyncStatusResponse, error)

	// Integration Management
	GetIntegrations(ctx context.Context, userID int) ([]*dto.IntegrationResponse, error)

	// Queue Operations
	QueueSync(ctx context.Context, userID int, eventID int, action string) error
	ProcessSyncQueue(ctx context.Context, limit int) (int, error)
}
