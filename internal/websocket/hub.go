package websocket

import (
	"encoding/json"
	"sync"

	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 256),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Info().
				Str("user_id", client.userID).
				Msg("WebSocket client registered")

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Info().
					Str("user_id", client.userID).
					Msg("WebSocket client unregistered")
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage broadcasts a message to all connected clients
func (h *Hub) BroadcastMessage(message []byte) {
	h.broadcast <- message
}

// Register registers a new client with the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client.userID == userID {
			select {
			case client.send <- message:
			default:
				log.Warn().
					Str("user_id", userID).
					Msg("Failed to send message to user")
			}
		}
	}
}

// BroadcastJSON broadcasts a JSON message to all connected clients
func (h *Hub) BroadcastJSON(v interface{}) error {
	message, err := json.Marshal(v)
	if err != nil {
		return err
	}
	h.BroadcastMessage(message)
	return nil
}

// BroadcastToUserJSON sends a JSON message to a specific user
func (h *Hub) BroadcastToUserJSON(userID string, v interface{}) error {
	message, err := json.Marshal(v)
	if err != nil {
		return err
	}
	h.BroadcastToUser(userID, message)
	return nil
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetUserClients returns all clients for a specific user
func (h *Hub) GetUserClients(userID string) []*Client {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := []*Client{}
	for client := range h.clients {
		if client.userID == userID {
			clients = append(clients, client)
		}
	}
	return clients
}

// GetStats returns detailed statistics about the hub
func (h *Hub) GetStats() HubStats {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userConnections := make(map[string]int)
	for client := range h.clients {
		userConnections[client.userID]++
	}

	return HubStats{
		TotalConnections:   len(h.clients),
		UniqueUsers:        len(userConnections),
		UserConnections:    userConnections,
		BroadcastQueueLen:  len(h.broadcast),
		RegisterQueueLen:   len(h.register),
		UnregisterQueueLen: len(h.unregister),
	}
}
