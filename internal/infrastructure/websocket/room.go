package websocket

import (
	"sync"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
)

// Room represents a user-specific room for broadcasting
type Room struct {
	userID  int
	clients map[*Client]bool
	mu      sync.RWMutex
}

// NewRoom creates a new room for a user
func NewRoom(userID int) *Room {
	return &Room{
		userID:  userID,
		clients: make(map[*Client]bool),
	}
}

// AddClient adds a client to the room
func (r *Room) AddClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[client] = true
}

// RemoveClient removes a client from the room
func (r *Room) RemoveClient(client *Client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, client)
}

// IsEmpty returns true if the room has no clients
func (r *Room) IsEmpty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients) == 0
}

// ClientCount returns the number of clients in the room
func (r *Room) ClientCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

// Broadcast sends a message to all clients in the room
func (r *Room) Broadcast(data []byte, logger *logger.ZapLogger) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for client := range r.clients {
		select {
		case client.send <- data:
			// Message queued successfully
		default:
			// Client buffer full, will be cleaned up by hub
			logger.Error("Client send buffer full, dropping message", nil, map[string]interface{}{
				"user_id": r.userID,
				"action":  "WS_BUFFER_FULL",
			})
		}
	}
}

// GetClients returns all clients in the room (for iteration)
func (r *Room) GetClients() []*Client {
	r.mu.RLock()
	defer r.mu.RUnlock()

	clients := make([]*Client, 0, len(r.clients))
	for client := range r.clients {
		clients = append(clients, client)
	}
	return clients
}
