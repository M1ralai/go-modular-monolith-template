package jobimpl

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/repository"
)

// HabitSkipJob skips a habit asynchronously
type HabitSkipJob struct {
	jobs.BaseJob
	logger      *logger.ZapLogger
	repo        repository.HabitRepository
	broadcaster *notifService.Broadcaster
	habitID     int
	userID      int
}

// LockKey returns a unique lock key for this job instance
func (j *HabitSkipJob) LockKey() int64 {
	return jobs.LockKey(j.Name(), j.habitID)
}

// NewHabitSkipJob creates a new habit skip job
func NewHabitSkipJob(
	logger *logger.ZapLogger,
	repo repository.HabitRepository,
	broadcaster *notifService.Broadcaster,
	habitID, userID int,
) *HabitSkipJob {
	return &HabitSkipJob{
		BaseJob:     jobs.NewBaseJob("habit_skip", "", 30*time.Second, nil),
		logger:      logger,
		repo:        repo,
		broadcaster: broadcaster,
		habitID:     habitID,
		userID:      userID,
	}
}

func (j *HabitSkipJob) Execute(ctx context.Context) error {
	j.logger.Info("Habit skip job started", map[string]interface{}{
		"habit_id": j.habitID,
		"user_id":  j.userID,
		"lock_key": j.LockKey(),
		"action":   "HABIT_SKIP_JOB_STARTED",
	})

	// Use context with longer timeout for database operations
	dbCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Get habit
	habit, err := j.repo.GetByID(dbCtx, j.habitID)
	if err != nil {
		j.logger.Error("Failed to get habit in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_SKIP_JOB_FAILED",
		})
		return err
	}

	if habit == nil {
		err := errors.New("habit not found")
		j.logger.Error("Habit not found in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_SKIP_JOB_NOT_FOUND",
		})
		return err
	}

	if habit.UserID != j.userID {
		err := errors.New("unauthorized")
		j.logger.Error("Unauthorized habit skip in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_SKIP_JOB_UNAUTHORIZED",
		})
		return err
	}

	// Skip habit for today
	today := time.Now().Truncate(24 * time.Hour)
	if err := j.repo.SkipHabit(dbCtx, j.habitID, today, ""); err != nil {
		j.logger.Error("Failed to skip habit in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_SKIP_JOB_SKIP_FAILED",
		})
		return err
	}

	// Broadcast WebSocket message
	if j.broadcaster != nil {
		j.broadcaster.Publish(j.userID, "habit.skipped", map[string]interface{}{
			"habit_id": j.habitID,
			"title":    habit.Name,
		})
	}

	j.logger.Info("Habit skip job completed", map[string]interface{}{
		"habit_id": j.habitID,
		"user_id":  j.userID,
		"action":   "HABIT_SKIP_JOB_COMPLETED",
	})

	return nil
}
