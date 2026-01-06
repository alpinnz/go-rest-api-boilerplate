package dto

// WebSocketStatsResponse represents WebSocket statistics response
type WebSocketStatsResponse struct {
	TotalConnections   int                   `json:"total_connections" example:"25"`
	UniqueUsers        int                   `json:"unique_users" example:"15"`
	BroadcastQueueLen  int                   `json:"broadcast_queue_length" example:"0"`
	RegisterQueueLen   int                   `json:"register_queue_length" example:"0"`
	UnregisterQueueLen int                   `json:"unregister_queue_length" example:"0"`
	TopUsers           []UserConnectionStats `json:"top_users"`
}

// UserConnectionStats represents connection statistics per user
type UserConnectionStats struct {
	UserID      string `json:"user_id" example:"user-123"`
	Connections int    `json:"connections" example:"2"`
}
