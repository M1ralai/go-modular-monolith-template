package notification

const (
	// Task events
	EventTaskCreated   = "task_created"
	EventTaskUpdated   = "task_updated"
	EventTaskCompleted = "task_completed"
	EventTaskDeleted   = "task_deleted"

	// Habit events
	EventHabitReminder   = "habit_reminder"
	EventHabitCompleted  = "habit_completed"
	EventHabitSkipped    = "habit_skipped"
	EventStreakIncreased = "streak_increased"
	EventStreakMilestone = "streak_milestone"
	EventStreakBroken    = "streak_broken"

	// Calendar events
	EventSyncStarted   = "sync_started"
	EventSyncProgress  = "sync_progress"
	EventSyncCompleted = "sync_completed"
	EventSyncFailed    = "sync_failed"

	// Event module events
	EventEventCreated = "event_created"
	EventEventUpdated = "event_updated"
	EventEventDeleted = "event_deleted"

	// Schedule events
	EventConflictDetected = "conflict_detected"
	EventEventsGenerated  = "events_generated"

	// Goal events
	EventGoalCreated      = "goal_created"
	EventGoalCompleted    = "goal_completed"
	EventMilestoneReached = "milestone_reached"

	// System events
	EventConnected = "connected"
	EventError     = "error"
)

func GetMilestoneName(streak int) string {
	switch streak {
	case 10:
		return "Getting Started 🌱"
	case 30:
		return "Building Momentum 🔥"
	case 50:
		return "Halfway Hero 🏅"
	case 100:
		return "Century Club 💯"
	default:
		return "Milestone Reached 🎯"
	}
}
