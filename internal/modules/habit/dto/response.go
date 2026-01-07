package dto

import (
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/domain"
)

type HabitResponse struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	LifeAreaID     *int      `json:"life_area_id,omitempty"`
	Title          string    `json:"title"`
	Icon           string    `json:"icon"`
	Description    string    `json:"description,omitempty"`
	Frequency      string    `json:"frequency"`
	FrequencyDays  []string  `json:"frequencyDays,omitempty"`
	IntervalDays   int       `json:"intervalDays,omitempty"`
	TargetCount    int       `json:"target_count"`
	CurrentStreak  int       `json:"current_streak"`
	BestStreak     int       `json:"bestStreak"`
	LongestStreak  int       `json:"longest_streak"`
	IsActive       bool      `json:"is_active"`
	CompletedToday bool      `json:"completed_today"`
	SkippedToday   bool      `json:"skipped_today"`
	TimeOfDay      string    `json:"timeOfDay,omitempty"`
	ReminderTime   string    `json:"reminderTime,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func ToHabitResponse(h *domain.Habit, completedToday bool, skippedToday bool) *HabitResponse {
	if h == nil {
		return nil
	}

	var freqDays []string
	if days, ok := h.FrequencyConfig["days"].([]interface{}); ok {
		for _, d := range days {
			if s, ok := d.(string); ok {
				freqDays = append(freqDays, s)
			}
		}
	} else if days, ok := h.FrequencyConfig["days"].([]string); ok {
		freqDays = days
	}

	var intervalDays int
	if interval, ok := h.FrequencyConfig["interval"].(float64); ok {
		intervalDays = int(interval)
	} else if interval, ok := h.FrequencyConfig["interval"].(int); ok {
		intervalDays = interval
	}

	return &HabitResponse{
		ID:             h.ID,
		UserID:         h.UserID,
		LifeAreaID:     h.LifeAreaID,
		Title:          h.Name,
		Icon:           h.Icon,
		Description:    h.Description,
		Frequency:      h.Frequency,
		FrequencyDays:  freqDays,
		IntervalDays:   intervalDays,
		TargetCount:    h.TargetCount,
		CurrentStreak:  h.CurrentStreak,
		BestStreak:     h.LongestStreak,
		LongestStreak:  h.LongestStreak,
		IsActive:       h.IsActive,
		CompletedToday: completedToday,
		SkippedToday:   skippedToday,
		TimeOfDay:      h.TimeOfDay,
		ReminderTime:   h.ReminderTime,
		CreatedAt:      h.CreatedAt,
		UpdatedAt:      h.UpdatedAt,
	}
}

func ToHabitResponseList(habits []*domain.Habit) []*HabitResponse {
	result := make([]*HabitResponse, len(habits))
	for i, h := range habits {
		result[i] = ToHabitResponse(h, false, false)
	}
	return result
}
