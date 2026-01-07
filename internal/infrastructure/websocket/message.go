package websocket

import (
	"time"

	"github.com/google/uuid"
)

// Message Types - Connection
const (
	TypePing = "ping"
	TypePong = "pong"
	TypeAuth = "auth"
)

// Message Types - Tasks
const (
	TypeTaskCreated   = "task.created"
	TypeTaskUpdated   = "task.updated"
	TypeTaskCompleted = "task.completed"
	TypeTaskDeleted   = "task.deleted"
)

// Message Types - Habits
const (
	TypeHabitReminder   = "habit.reminder"
	TypeHabitMilestone  = "habit.milestone"
	TypeHabitStreakLost = "habit.streak_lost"
	TypeHabitCompleted  = "habit.completed"
)

// Message Types - Calendar
const (
	TypeCalendarSyncStatus      = "calendar.sync_status"
	TypeCalendarEventUpdated    = "calendar.event_updated"
	TypeCalendarConflict        = "calendar.conflict"
	TypeCalendarEventsGenerated = "calendar.events_generated"
)

// Message Types - Jobs
const (
	TypeJobStarted   = "job.started"
	TypeJobCompleted = "job.completed"
	TypeJobFailed    = "job.failed"
	TypeJobProgress  = "job.progress"
)

// Message Types - Sync
const (
	TypeDeviceSync = "device.sync"
	TypeConnected  = "connected"
	TypeError      = "error"
)

// Message represents a WebSocket message following the standard format
type Message struct {
	Type      string                 `json:"type"`
	Payload   map[string]interface{} `json:"payload"`
	Timestamp time.Time              `json:"timestamp"`
	MessageID string                 `json:"message_id,omitempty"`
	UserID    int                    `json:"-"` // Internal use only, not serialized
}

// NewMessage creates a new WebSocket message with auto-generated timestamp and message ID
func NewMessage(msgType string, userID int, payload map[string]interface{}) *Message {
	return &Message{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
		MessageID: uuid.New().String(),
		UserID:    userID,
	}
}

// NewMessageWithID creates a message with a specific message ID (for deduplication)
func NewMessageWithID(msgType string, userID int, payload map[string]interface{}, messageID string) *Message {
	return &Message{
		Type:      msgType,
		Payload:   payload,
		Timestamp: time.Now().UTC(),
		MessageID: messageID,
		UserID:    userID,
	}
}

// Task Event Payloads
type TaskPayload struct {
	TaskID      int        `json:"task_id"`
	Title       string     `json:"title"`
	Priority    string     `json:"priority,omitempty"`
	IsCompleted bool       `json:"is_completed,omitempty"`
	Progress    float64    `json:"progress,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	ParentID    *int       `json:"parent_id,omitempty"`
}

// Habit Event Payloads
type HabitPayload struct {
	HabitID       int    `json:"habit_id"`
	Title         string `json:"title"`
	CurrentStreak int    `json:"current_streak,omitempty"`
	BestStreak    int    `json:"best_streak,omitempty"`
	Milestone     string `json:"milestone,omitempty"`
	ReminderTime  string `json:"reminder_time,omitempty"`
}

// Calendar Event Payloads
type CalendarPayload struct {
	Provider       string `json:"provider,omitempty"`
	Progress       int    `json:"progress,omitempty"`
	EventsSynced   int    `json:"events_synced,omitempty"`
	ConflictReason string `json:"conflict_reason,omitempty"`
	EventID        int    `json:"event_id,omitempty"`
	Title          string `json:"title,omitempty"`
}

// Job Event Payloads
type JobPayload struct {
	JobType  string `json:"job_type"`
	JobID    string `json:"job_id,omitempty"`
	Progress int    `json:"progress,omitempty"`
	Message  string `json:"message,omitempty"`
	Error    string `json:"error,omitempty"`
}

// Helper functions to create typed payloads

func TaskToPayload(taskID int, title string, priority string) map[string]interface{} {
	return map[string]interface{}{
		"task_id":  taskID,
		"title":    title,
		"priority": priority,
	}
}

func HabitToPayload(habitID int, title string, streak int) map[string]interface{} {
	return map[string]interface{}{
		"habit_id":       habitID,
		"title":          title,
		"current_streak": streak,
	}
}

func CalendarSyncPayload(provider string, progress int) map[string]interface{} {
	return map[string]interface{}{
		"provider": provider,
		"progress": progress,
	}
}

func JobToPayload(jobType string, jobID string, progress int) map[string]interface{} {
	return map[string]interface{}{
		"job_type": jobType,
		"job_id":   jobID,
		"progress": progress,
	}
}
