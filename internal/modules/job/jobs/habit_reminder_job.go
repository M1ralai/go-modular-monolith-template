package jobimpl

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
)

// HabitReminderJob sends habit reminders every 5 minutes
type HabitReminderJob struct {
	jobs.BaseJob
	logger       *logger.ZapLogger
	eventEmitter jobs.JobEventEmitter
}

// NewHabitReminderJob creates a new habit reminder job
func NewHabitReminderJob(logger *logger.ZapLogger, emitter jobs.JobEventEmitter) *HabitReminderJob {
	return &HabitReminderJob{
		BaseJob:      jobs.NewBaseJob("habit_reminder", "*/5 * * * *", 2*time.Minute, nil),
		logger:       logger,
		eventEmitter: emitter,
	}
}

func (j *HabitReminderJob) Execute(ctx context.Context) error {
	j.logger.Info("Habit reminder job started", map[string]interface{}{
		"job":    j.Name(),
		"action": "HABIT_REMINDER_STARTED",
	})

	if j.eventEmitter != nil {
		j.eventEmitter.EmitJobStarted(ctx, j.Name())
	}

	// TODO: Implement habit reminder logic
	// 1. Get habits with reminder times matching current time window
	// 2. Check if habit was already completed today
	// 3. Send reminder notification via WebSocket

	j.logger.Info("Habit reminder job completed", map[string]interface{}{
		"job":    j.Name(),
		"action": "HABIT_REMINDER_COMPLETED",
	})

	if j.eventEmitter != nil {
		j.eventEmitter.EmitJobCompleted(ctx, j.Name(), map[string]interface{}{
			"reminders_sent": 0,
		})
	}

	return nil
}
