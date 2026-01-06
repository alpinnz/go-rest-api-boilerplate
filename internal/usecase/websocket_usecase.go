package usecase

import (
	"encoding/json"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/websocket"
	"github.com/rs/zerolog/log"
)

// WebSocketUseCase handles WebSocket business logic
type WebSocketUseCase struct {
	hub *websocket.Hub
}

// NewWebSocketUseCase creates a new WebSocketUseCase instance
func NewWebSocketUseCase(hub *websocket.Hub) *WebSocketUseCase {
	return &WebSocketUseCase{
		hub: hub,
	}
}

// BroadcastNotification broadcasts a notification to all connected clients
func (uc *WebSocketUseCase) BroadcastNotification(title, message, level string) error {
	msg := websocket.NewMessage(
		websocket.MessageTypeNotification,
		websocket.NotificationData{
			Title:   title,
			Message: message,
			Level:   level,
		},
		"system",
		"",
	)

	return uc.hub.BroadcastJSON(msg)
}

// SendNotificationToUser sends a notification to a specific user
func (uc *WebSocketUseCase) SendNotificationToUser(userID, title, message, level string) error {
	msg := websocket.NewMessage(
		websocket.MessageTypeNotification,
		websocket.NotificationData{
			Title:   title,
			Message: message,
			Level:   level,
		},
		"system",
		userID,
	)

	return uc.hub.BroadcastToUserJSON(userID, msg)
}

// BroadcastSystemEvent broadcasts a system event to all connected clients
func (uc *WebSocketUseCase) BroadcastSystemEvent(event string, payload interface{}) error {
	msg := websocket.NewMessage(
		websocket.MessageTypeSystem,
		websocket.SystemData{
			Event:   event,
			Payload: payload,
		},
		"system",
		"",
	)

	return uc.hub.BroadcastJSON(msg)
}

// SendSystemEventToUser sends a system event to a specific user
func (uc *WebSocketUseCase) SendSystemEventToUser(userID, event string, payload interface{}) error {
	msg := websocket.NewMessage(
		websocket.MessageTypeSystem,
		websocket.SystemData{
			Event:   event,
			Payload: payload,
		},
		"system",
		userID,
	)

	return uc.hub.BroadcastToUserJSON(userID, msg)
}

// BroadcastMessage broadcasts a raw message to all connected clients
func (uc *WebSocketUseCase) BroadcastMessage(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return err
	}

	uc.hub.BroadcastMessage(data)
	return nil
}

// SendMessageToUser sends a raw message to a specific user
func (uc *WebSocketUseCase) SendMessageToUser(userID string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return err
	}

	uc.hub.BroadcastToUser(userID, data)
	return nil
}

// GetConnectedClientsCount returns the number of connected clients
func (uc *WebSocketUseCase) GetConnectedClientsCount() int {
	return uc.hub.GetClientCount()
}
