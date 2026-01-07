package websocket

import (
	"net/http"
	"os"
	"strconv"

	"github.com/M1ralai/go-modular-monolith-template/internal/infrastructure/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type Handler struct {
	hub    *Hub
	logger *logger.ZapLogger
}

func NewHandler(hub *Hub, logger *logger.ZapLogger) *Handler {
	return &Handler{hub: hub, logger: logger}
}

func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Token required", http.StatusUnauthorized)
		return
	}

	userID, err := h.validateToken(token)
	if err != nil {
		h.logger.Error("WebSocket auth failed", err, map[string]interface{}{
			"action": "WS_AUTH_FAILED",
		})
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("WebSocket upgrade failed", err, map[string]interface{}{
			"user_id": userID,
			"action":  "WS_UPGRADE_FAILED",
		})
		return
	}

	client := NewClient(h.hub, conn, userID)
	h.hub.register <- client

	go client.WritePump()
	go client.ReadPump()

	// Send welcome message
	h.hub.PublishToUser(userID, NewMessage(TypeConnected, userID, map[string]interface{}{
		"message": "Connected to WebSocket",
	}))
}

func (h *Handler) validateToken(tokenString string) (int, error) {
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIDFloat, ok := claims["user_id"].(float64); ok {
			return int(userIDFloat), nil
		}
		if userIDStr, ok := claims["user_id"].(string); ok {
			return strconv.Atoi(userIDStr)
		}
	}

	return 0, jwt.ErrSignatureInvalid
}
