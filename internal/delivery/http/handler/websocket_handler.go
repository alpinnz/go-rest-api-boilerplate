package handler

import (
	"net/http"
	"sort"

	"github.com/alpinnz/go-rest-api-boilerplate/internal/delivery/http/dto"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/domain"
	"github.com/alpinnz/go-rest-api-boilerplate/internal/websocket"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/context"
	"github.com/alpinnz/go-rest-api-boilerplate/pkg/response"
	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// TODO: Configure this based on your CORS policy
		return true
	},
}

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub *websocket.Hub
}

// NewWebSocketHandler creates a new WebSocketHandler instance
func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// Connect godoc
// @Summary WebSocket connection endpoint
// @Description Establish a WebSocket connection for real-time communication
// @Tags websocket
// @Security BearerAuth
// @Success 101 {string} string "Switching Protocols"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /ws [get]
func (h *WebSocketHandler) Connect(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := context.GetUserID(c.Request.Context())
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, domain.NewAppError("auth.unauthorized"))
		return
	}

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, domain.NewAppError("websocket.upgrade_failed"))
		return
	}

	// Create new client
	client := websocket.NewClient(h.hub, conn, userID)
	h.hub.Register(client)

	// Start client goroutines
	go client.WritePump()
	go client.ReadPump()
}

// GetStats godoc
// @Summary Get WebSocket statistics
// @Description Get detailed WebSocket connection statistics including total connections, unique users, and top connected users
// @Tags websocket
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{data=dto.WebSocketStatsResponse}
// @Failure 401 {object} response.Response "Unauthorized"
// @Router /ws/stats [get]
func (h *WebSocketHandler) GetStats(c *gin.Context) {
	stats := h.hub.GetStats()

	// Convert user connections map to slice and sort by connection count
	topUsers := make([]dto.UserConnectionStats, 0, len(stats.UserConnections))
	for userID, count := range stats.UserConnections {
		topUsers = append(topUsers, dto.UserConnectionStats{
			UserID:      userID,
			Connections: count,
		})
	}

	// Sort by connection count descending
	sort.Slice(topUsers, func(i, j int) bool {
		return topUsers[i].Connections > topUsers[j].Connections
	})

	// Limit to top 10 users
	if len(topUsers) > 10 {
		topUsers = topUsers[:10]
	}

	responseStats := dto.WebSocketStatsResponse{
		TotalConnections:   stats.TotalConnections,
		UniqueUsers:        stats.UniqueUsers,
		BroadcastQueueLen:  stats.BroadcastQueueLen,
		RegisterQueueLen:   stats.RegisterQueueLen,
		UnregisterQueueLen: stats.UnregisterQueueLen,
		TopUsers:           topUsers,
	}

	response.Success(c, responseStats)
}
