package domain

import "time"

type LifeArea struct {
	ID           int
	UserID       int
	Name         string
	Icon         string
	Color        string
	DisplayOrder int
	CreatedAt    time.Time
}
