package jobs

import (
	"context"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/websocket"
)

// WebSocketJobEmitter implements JobEventEmitter using WebSocket
type WebSocketJobEmitter struct {
	hub    *websocket.Hub
	logger *logger.ZapLogger
}

// NewWebSocketJobEmitter creates a new WebSocket job event emitter
func NewWebSocketJobEmitter(hub *websocket.Hub, logger *logger.ZapLogger) *WebSocketJobEmitter {
	return &WebSocketJobEmitter{
		hub:    hub,
		logger: logger,
	}
}

func (e *WebSocketJobEmitter) EmitJobStarted(ctx context.Context, jobName string) {
	e.logger.Info("Emitting job started event", map[string]interface{}{
		"job":    jobName,
		"action": "JOB_EVENT_STARTED",
	})

	msg := websocket.NewMessage(websocket.TypeJobStarted, 0, map[string]interface{}{
		"job_name": jobName,
		"status":   "running",
	})

	e.hub.BroadcastToAll(msg)
}

func (e *WebSocketJobEmitter) EmitJobProgress(ctx context.Context, jobName string, progress float64, message string) {
	e.logger.Info("Emitting job progress event", map[string]interface{}{
		"job":      jobName,
		"progress": progress,
		"action":   "JOB_EVENT_PROGRESS",
	})

	msg := websocket.NewMessage(websocket.TypeJobProgress, 0, map[string]interface{}{
		"job_name": jobName,
		"progress": progress,
		"message":  message,
	})

	e.hub.BroadcastToAll(msg)
}

func (e *WebSocketJobEmitter) EmitJobCompleted(ctx context.Context, jobName string, result interface{}) {
	e.logger.Info("Emitting job completed event", map[string]interface{}{
		"job":    jobName,
		"action": "JOB_EVENT_COMPLETED",
	})

	msg := websocket.NewMessage(websocket.TypeJobCompleted, 0, map[string]interface{}{
		"job_name": jobName,
		"status":   "completed",
		"result":   result,
	})

	e.hub.BroadcastToAll(msg)
}

func (e *WebSocketJobEmitter) EmitJobFailed(ctx context.Context, jobName string, err error) {
	e.logger.Info("Emitting job failed event", map[string]interface{}{
		"job":    jobName,
		"error":  err.Error(),
		"action": "JOB_EVENT_FAILED",
	})

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	msg := websocket.NewMessage(websocket.TypeJobFailed, 0, map[string]interface{}{
		"job_name": jobName,
		"status":   "failed",
		"error":    errMsg,
	})

	e.hub.BroadcastToAll(msg)
}
