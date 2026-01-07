package dto

type CreateHabitRequest struct {
	LifeAreaID    *int     `json:"life_area_id,omitempty"`
	Name          string   `json:"name" validate:"required,min=1,max=255"`
	Icon          string   `json:"icon,omitempty"`
	Description   string   `json:"description,omitempty"`
	Frequency     string   `json:"frequency" validate:"required,oneof=daily weekly custom"`
	FrequencyDays []string `json:"frequencyDays,omitempty"`
	IntervalDays  *int     `json:"intervalDays,omitempty"`
	TimeOfDay     string   `json:"timeOfDay,omitempty"`
	ReminderTime  string   `json:"reminderTime,omitempty"`
	TargetCount   int      `json:"target_count,omitempty"`
	TargetDays    int      `json:"targetDays,omitempty"`
}

type UpdateHabitRequest struct {
	LifeAreaID    *int     `json:"life_area_id,omitempty"`
	Name          *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Icon          *string  `json:"icon,omitempty"`
	Description   *string  `json:"description,omitempty"`
	Frequency     *string  `json:"frequency,omitempty" validate:"omitempty,oneof=daily weekly custom"`
	FrequencyDays []string `json:"frequencyDays,omitempty"`
	IntervalDays  *int     `json:"intervalDays,omitempty"`
	TargetCount   *int     `json:"target_count,omitempty" validate:"omitempty,min=1"`
	TargetDays    *int     `json:"targetDays,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
	TimeOfDay     *string  `json:"timeOfDay,omitempty"`
	ReminderTime  *string  `json:"reminderTime,omitempty"`
}

type LogHabitRequest struct {
	Count int    `json:"count" validate:"min=0"`
	Notes string `json:"notes,omitempty"`
}
