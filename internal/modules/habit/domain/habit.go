package domain

import "time"

type Habit struct {
	ID              int
	UserID          int
	LifeAreaID      *int
	Name            string
	Icon            string
	Description     string
	Frequency       string
	FrequencyConfig map[string]interface{}
	TargetCount     int
	TimeOfDay       string
	ReminderTime    string
	CurrentStreak   int
	LongestStreak   int
	IsActive        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (h *Habit) IncrementStreak() {
	h.CurrentStreak++
	if h.CurrentStreak > h.LongestStreak {
		h.LongestStreak = h.CurrentStreak
	}
	h.UpdatedAt = time.Now()
}

func (h *Habit) ResetStreak() {
	h.CurrentStreak = 0
	h.UpdatedAt = time.Now()
}

func (h *Habit) IsDaily() bool {
	return h.Frequency == "daily"
}
