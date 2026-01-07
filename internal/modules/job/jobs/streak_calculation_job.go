package jobimpl

import (
	"context"
	"time"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/jobs"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
)

// StreakCalculationJob calculates habit streaks daily
type StreakCalculationJob struct {
	jobs.BaseJob
	logger       *logger.ZapLogger
	eventEmitter jobs.JobEventEmitter
}

// NewStreakCalculationJob creates a new streak calculation job
func NewStreakCalculationJob(logger *logger.ZapLogger, emitter jobs.JobEventEmitter) *StreakCalculationJob {
	return &StreakCalculationJob{
		BaseJob:      jobs.NewBaseJob("streak_calculation", "0 0 * * *", 10*time.Minute, nil),
		logger:       logger,
		eventEmitter: emitter,
	}
}

func (j *StreakCalculationJob) Execute(ctx context.Context) error {
	j.logger.Info("Streak calculation job started", map[string]interface{}{
		"job":    j.Name(),
		"action": "STREAK_CALC_STARTED",
	})

	if j.eventEmitter != nil {
		j.eventEmitter.EmitJobStarted(ctx, j.Name())
	}

	// TODO: Implement streak calculation logic
	// 1. Get all active habits
	// 2. Check if habit was completed yesterday
	// 3. Update streak counts
	// 4. Detect broken streaks and emit events

	j.logger.Info("Streak calculation job completed", map[string]interface{}{
		"job":    j.Name(),
		"action": "STREAK_CALC_COMPLETED",
	})

	if j.eventEmitter != nil {
		j.eventEmitter.EmitJobCompleted(ctx, j.Name(), map[string]interface{}{
			"habits_processed": 0,
			"broken_streaks":   0,
		})
	}

	return nil
}
