package jobimpl

import (
	"context"
	"errors"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	notifService "github.com/M1ralai/go-modular-monolith-template/internal/modules/notification/service"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/dto"
	"github.com/M1ralai/go-modular-monolith-template/internal/modules/habit/repository"
)

// HabitCompleteJob completes a habit asynchronously
type HabitCompleteJob struct {
	jobs.BaseJob
	logger      *logger.ZapLogger
	repo        repository.HabitRepository
	broadcaster *notifService.Broadcaster
	habitID     int
	userID      int
	request     *dto.LogHabitRequest
}

// LockKey returns a unique lock key for this job instance
func (j *HabitCompleteJob) LockKey() int64 {
	return jobs.LockKey(j.Name(), j.habitID)
}

// NewHabitCompleteJob creates a new habit complete job
func NewHabitCompleteJob(
	logger *logger.ZapLogger,
	repo repository.HabitRepository,
	broadcaster *notifService.Broadcaster,
	habitID, userID int,
	request *dto.LogHabitRequest,
) *HabitCompleteJob {
	return &HabitCompleteJob{
		BaseJob:     jobs.NewBaseJob("habit_complete", "", 30*time.Second, nil),
		logger:      logger,
		repo:        repo,
		broadcaster: broadcaster,
		habitID:     habitID,
		userID:      userID,
		request:     request,
	}
}

func (j *HabitCompleteJob) Execute(ctx context.Context) error {
	j.logger.Info("Habit complete job started", map[string]interface{}{
		"habit_id": j.habitID,
		"user_id":  j.userID,
		"lock_key": j.LockKey(),
		"action":   "HABIT_COMPLETE_JOB_STARTED",
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
			"action":   "HABIT_COMPLETE_JOB_FAILED",
		})
		return err
	}

	if habit == nil {
		err := errors.New("habit not found")
		j.logger.Error("Habit not found in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_COMPLETE_JOB_NOT_FOUND",
		})
		return err
	}

	if habit.UserID != j.userID {
		err := errors.New("unauthorized")
		j.logger.Error("Unauthorized habit complete in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_COMPLETE_JOB_UNAUTHORIZED",
		})
		return err
	}

	// Check if habit is already completed today
	alreadyCompleted, err := j.repo.HasLogForToday(dbCtx, j.habitID)
	if err != nil {
		j.logger.Error("Failed to check habit log in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_COMPLETE_JOB_CHECK_FAILED",
		})
		return err
	}

	if alreadyCompleted {
		err := errors.New("habit already completed today")
		j.logger.Error("Habit already completed today in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_COMPLETE_JOB_ALREADY_COMPLETED",
		})
		return err
	}

	// Prepare count
	count := j.request.Count
	if count < habit.TargetCount {
		count = habit.TargetCount
	}

	// Log habit for today
	today := time.Now().Truncate(24 * time.Hour)
	if err := j.repo.LogHabit(dbCtx, j.habitID, today, count, j.request.Notes); err != nil {
		j.logger.Error("Failed to log habit in job", err, map[string]interface{}{
			"habit_id": j.habitID,
			"user_id":  j.userID,
			"action":   "HABIT_COMPLETE_JOB_LOG_FAILED",
		})
		return err
	}

	// If count >= target, increment streak
	oldStreak := habit.CurrentStreak
	if count >= habit.TargetCount {
		habit.IncrementStreak()
		if err := j.repo.Update(dbCtx, habit); err != nil {
			j.logger.Error("Failed to update habit streak in job", err, map[string]interface{}{
				"habit_id": j.habitID,
				"user_id":  j.userID,
				"action":   "HABIT_COMPLETE_JOB_UPDATE_FAILED",
			})
			return err
		}
	}

	// Broadcast WebSocket message
	if j.broadcaster != nil {
		j.broadcaster.Publish(j.userID, "habit.completed", map[string]interface{}{
			"habit_id": j.habitID,
			"title":    habit.Name,
			"streak":   habit.CurrentStreak,
		})

		// If streak increased, notify
		if habit.CurrentStreak > oldStreak {
			j.broadcaster.Publish(j.userID, "habit.streak_increased", map[string]interface{}{
				"habit_id": j.habitID,
				"streak":   habit.CurrentStreak,
			})

			// Check for milestone (every 10 days)
			if habit.CurrentStreak%10 == 0 && habit.CurrentStreak > 0 {
				j.broadcaster.Publish(j.userID, "habit.milestone", map[string]interface{}{
					"habit_id":  j.habitID,
					"title":     habit.Name,
					"streak":    habit.CurrentStreak,
					"milestone": habit.CurrentStreak,
				})
			}
		}
	}

	j.logger.Info("Habit complete job completed", map[string]interface{}{
		"habit_id": j.habitID,
		"user_id":  j.userID,
		"streak":   habit.CurrentStreak,
		"action":   "HABIT_COMPLETE_JOB_COMPLETED",
	})

	return nil
}
