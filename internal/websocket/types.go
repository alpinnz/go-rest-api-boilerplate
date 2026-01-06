package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	MessageTypeNotification MessageType = "notification"
	MessageTypeChat         MessageType = "chat"
	MessageTypeSystem       MessageType = "system"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
)

// Message represents a WebSocket message structure
type Message struct {
	Type      MessageType            `json:"type"`
	Data      interface{}            `json:"data"`
	From      string                 `json:"from,omitempty"`
	To        string                 `json:"to,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// NewMessage creates a new Message instance
func NewMessage(msgType MessageType, data interface{}, from, to string) *Message {
	return &Message{
		Type:      msgType,
		Data:      data,
		From:      from,
		To:        to,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
}

// NotificationData represents notification data
type NotificationData struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Level   string `json:"level"` // info, success, warning, error
}

// SystemData represents system message data
type SystemData struct {
	Event   string      `json:"event"`
	Payload interface{} `json:"payload,omitempty"`
}

// HubStats represents WebSocket hub statistics
type HubStats struct {
	TotalConnections   int            `json:"total_connections"`
	UniqueUsers        int            `json:"unique_users"`
	UserConnections    map[string]int `json:"user_connections"`
	BroadcastQueueLen  int            `json:"broadcast_queue_length"`
	RegisterQueueLen   int            `json:"register_queue_length"`
	UnregisterQueueLen int            `json:"unregister_queue_length"`
}
