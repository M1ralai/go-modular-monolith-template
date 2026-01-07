package domain

import "time"

type CalendarIntegration struct {
	ID           int
	UserID       int
	Provider     string
	AccessToken  string
	RefreshToken string
	ExpiresAt    *time.Time
	CalendarID   string
	IsActive     bool
	LastSyncAt   *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c *CalendarIntegration) IsExpired() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*c.ExpiresAt)
}

func (c *CalendarIntegration) NeedsRefresh() bool {
	if c.ExpiresAt == nil {
		return false
	}
	return time.Now().Add(5 * time.Minute).After(*c.ExpiresAt)
}
