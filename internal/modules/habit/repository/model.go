package repository

import (
	"encoding/json"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/domain"
)

type HabitModel struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	LifeAreaID      *int      `db:"life_area_id"`
	Name            string    `db:"name"`
	Icon            *string   `db:"icon"`
	Description     *string   `db:"description"`
	Frequency       string    `db:"frequency"`
	FrequencyConfig []byte    `db:"frequency_config"`
	TargetCount     int       `db:"target_count"`
	TimeOfDay       *string   `db:"time_of_day"`
	ReminderTime    *string   `db:"reminder_time"`
	CurrentStreak   int       `db:"current_streak"`
	LongestStreak   int       `db:"longest_streak"`
	IsActive        bool      `db:"is_active"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func (m *HabitModel) ToDomain() *domain.Habit {
	if m == nil {
		return nil
	}
	desc := ""
	if m.Description != nil {
		desc = *m.Description
	}
	icon := ""
	if m.Icon != nil {
		icon = *m.Icon
	}
	tod := ""
	if m.TimeOfDay != nil {
		tod = *m.TimeOfDay
	}
	rem := ""
	if m.ReminderTime != nil {
		rem = *m.ReminderTime
	}

	var config map[string]interface{}
	if len(m.FrequencyConfig) > 0 {
		json.Unmarshal(m.FrequencyConfig, &config)
	}

	return &domain.Habit{
		ID:              m.ID,
		UserID:          m.UserID,
		LifeAreaID:      m.LifeAreaID,
		Name:            m.Name,
		Icon:            icon,
		Description:     desc,
		Frequency:       m.Frequency,
		FrequencyConfig: config,
		TargetCount:     m.TargetCount,
		TimeOfDay:       tod,
		ReminderTime:    rem,
		CurrentStreak:   m.CurrentStreak,
		LongestStreak:   m.LongestStreak,
		IsActive:        m.IsActive,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}

func FromDomain(h *domain.Habit) *HabitModel {
	if h == nil {
		return nil
	}
	var desc *string
	if h.Description != "" {
		desc = &h.Description
	}
	var icon *string
	if h.Icon != "" {
		icon = &h.Icon
	}
	var tod *string
	if h.TimeOfDay != "" {
		tod = &h.TimeOfDay
	}
	var rem *string
	if h.ReminderTime != "" {
		rem = &h.ReminderTime
	}

	config, _ := json.Marshal(h.FrequencyConfig)

	return &HabitModel{
		ID:              h.ID,
		UserID:          h.UserID,
		LifeAreaID:      h.LifeAreaID,
		Name:            h.Name,
		Icon:            icon,
		Description:     desc,
		Frequency:       h.Frequency,
		FrequencyConfig: config,
		TargetCount:     h.TargetCount,
		TimeOfDay:       tod,
		ReminderTime:    rem,
		CurrentStreak:   h.CurrentStreak,
		LongestStreak:   h.LongestStreak,
		IsActive:        h.IsActive,
		CreatedAt:       h.CreatedAt,
		UpdatedAt:       h.UpdatedAt,
	}
}

type HabitLogModel struct {
	ID          int       `db:"id"`
	HabitID     int       `db:"habit_id"`
	LogDate     time.Time `db:"log_date"`
	Count       int       `db:"count"`
	Notes       *string   `db:"notes"`
	IsCompleted bool      `db:"is_completed"`
	Skipped     bool      `db:"skipped"`
	CreatedAt   time.Time `db:"created_at"`
}
