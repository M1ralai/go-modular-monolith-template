package notification

const (
	// Task events
	EventTaskCreated   = "task.created"
	EventTaskUpdated   = "task.updated"
	EventTaskCompleted = "task.completed"
	EventTaskDeleted   = "task.deleted"

	// Habit events
	EventHabitCreated   = "habit.created"
	EventHabitUpdated   = "habit.updated"
	EventHabitDeleted   = "habit.deleted"
	EventHabitReminder   = "habit.reminder"
	EventHabitCompleted  = "habit.completed"
	EventHabitSkipped    = "habit.skipped"
	EventStreakIncreased = "streak.increased"
	EventStreakMilestone = "habit.milestone"
	EventStreakBroken    = "streak.broken"

	// Course events
	EventCourseCreated     = "course.created"
	EventCourseUpdated     = "course.updated"
	EventCourseDeleted     = "course.deleted"
	EventComponentCreated  = "component.created"
	EventComponentUpdated  = "component.updated"
	EventComponentDeleted  = "component.deleted"
	EventComponentGraded   = "component.graded"
	EventScheduleCreated   = "schedule.created"
	EventScheduleUpdated   = "schedule.updated"
	EventScheduleDeleted   = "schedule.deleted"

	// Calendar events
	EventSyncStarted   = "sync.started"
	EventSyncProgress  = "sync.progress"
	EventSyncCompleted = "sync.completed"
	EventSyncFailed    = "sync.failed"

	// Event module events
	EventEventCreated = "event.created"
	EventEventUpdated = "event.updated"
	EventEventDeleted = "event.deleted"

	// Schedule events
	EventConflictDetected = "conflict.detected"
	EventEventsGenerated  = "events.generated"

	// Goal events
	EventGoalCreated      = "goal.created"
	EventGoalUpdated      = "goal.updated"
	EventGoalDeleted      = "goal.deleted"
	EventGoalCompleted    = "goal.completed"
	EventMilestoneReached = "milestone.reached"

	// Note events
	EventNoteCreated = "note.created"
	EventNoteUpdated = "note.updated"
	EventNoteDeleted = "note.deleted"

	// LifeArea events
	EventLifeAreaCreated = "lifearea.created"
	EventLifeAreaUpdated = "lifearea.updated"
	EventLifeAreaDeleted = "lifearea.deleted"

	// People events
	EventPersonCreated = "person.created"
	EventPersonUpdated = "person.updated"
	EventPersonDeleted = "person.deleted"

	// Journal events
	EventJournalCreated = "journal.created"
	EventJournalUpdated = "journal.updated"
	EventJournalDeleted = "journal.deleted"

	// Finance events
	EventTransactionCreated = "transaction.created"
	EventTransactionUpdated = "transaction.updated"
	EventTransactionDeleted = "transaction.deleted"

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
