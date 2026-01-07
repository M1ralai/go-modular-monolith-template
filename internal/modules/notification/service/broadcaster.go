package notification

import (
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/websocket"
)

type Broadcaster struct {
	hub    *websocket.Hub
	logger *logger.ZapLogger
}

func NewBroadcaster(hub *websocket.Hub, logger *logger.ZapLogger) *Broadcaster {
	return &Broadcaster{hub: hub, logger: logger}
}

func (b *Broadcaster) Publish(userID int, eventType string, data map[string]interface{}) {
	b.logger.Info("Broadcasting notification", map[string]interface{}{
		"user_id": userID,
		"type":    eventType,
		"action":  "BROADCAST_NOTIFICATION",
	})

	message := websocket.NewMessage(eventType, userID, data)
	b.hub.PublishToUser(userID, message)
}

func (b *Broadcaster) TaskCreated(userID int, taskID int, title string) {
	b.Publish(userID, websocket.TypeTaskCreated, map[string]interface{}{
		"task_id": taskID,
		"title":   title,
	})
}

func (b *Broadcaster) TaskCompleted(userID int, taskID int, title string) {
	b.Publish(userID, websocket.TypeTaskCompleted, map[string]interface{}{
		"task_id": taskID,
		"title":   title,
	})
}

func (b *Broadcaster) HabitCompleted(userID int, habitID int, title string, streak int) {
	b.Publish(userID, websocket.TypeHabitCompleted, map[string]interface{}{
		"habit_id": habitID,
		"title":    title,
		"streak":   streak,
	})
}

func (b *Broadcaster) StreakMilestone(userID int, habitID int, title string, streak int) {
	b.Publish(userID, websocket.TypeHabitMilestone, map[string]interface{}{
		"habit_id": habitID,
		"title":    title,
		"streak":   streak,
	})
}

func (b *Broadcaster) SyncProgress(userID int, provider string, progress int) {
	b.Publish(userID, websocket.TypeJobProgress, map[string]interface{}{
		"provider": provider,
		"progress": progress,
	})
}

func (b *Broadcaster) SyncCompleted(userID int, provider string) {
	b.Publish(userID, websocket.TypeJobCompleted, map[string]interface{}{
		"provider": provider,
	})
}

func (b *Broadcaster) ConflictDetected(userID int, reason string, start, end string) {
	b.Publish(userID, websocket.TypeCalendarConflict, map[string]interface{}{
		"reason": reason,
		"start":  start,
		"end":    end,
	})
}
