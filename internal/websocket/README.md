# WebSocket Implementation

Real-time bidirectional communication server integrated with this boilerplate.

## Features

- Real-time bidirectional communication
- Automatic connection management (Hub pattern)
- JWT authentication required
- Message broadcasting to all clients
- Targeted message sending to specific users
- Ping/Pong health checks (54 second interval)
- Graceful connection handling
- Thread-safe operations with sync.RWMutex
- Structured message types
- Detailed connection statistics
- Per-user connection tracking

## Architecture

```
internal/
├── websocket/
│   ├── hub.go          # Central connection manager
│   ├── client.go       # WebSocket client handler
│   └── types.go        # Message type definitions
├── delivery/http/handler/
│   └── websocket_handler.go  # HTTP to WebSocket upgrade handler
└── usecase/
    └── websocket_usecase.go   # Business logic layer
```

## Endpoints

### Connect to WebSocket
```
GET /api/v1/ws
```
- **Authentication**: Required (Bearer Token)
- **Protocol**: WebSocket
- **Response**: Upgrades to WebSocket connection

### Get WebSocket Statistics
```
GET /api/v1/ws/stats
```
- **Authentication**: Required (Bearer Token)
- **Response**: 
```json
{
  "code": null,
  "message": "Success",
  "data": {
    "total_connections": 25,
    "unique_users": 15,
    "broadcast_queue_length": 0,
    "register_queue_length": 0,
    "unregister_queue_length": 0,
    "top_users": [
      {
        "user_id": "user-123",
        "connections": 3
      },
      {
        "user_id": "user-456",
        "connections": 2
      }
    ]
  }
}
```

**Response Fields:**
- `total_connections`: Total active WebSocket connections
- `unique_users`: Number of unique users connected
- `broadcast_queue_length`: Number of messages in broadcast queue
- `register_queue_length`: Number of clients in registration queue
- `unregister_queue_length`: Number of clients in unregistration queue
- `top_users`: Top 10 users by connection count (useful for monitoring multi-device users)

## Message Types

### Notification Message
```json
{
  "type": "notification",
  "data": {
    "title": "Welcome",
    "message": "You are now connected",
    "level": "info"
  },
  "from": "system",
  "to": "",
  "timestamp": "2026-01-06T10:00:00Z"
}
```

### System Message
```json
{
  "type": "system",
  "data": {
    "event": "user_joined",
    "payload": {
      "user_id": "123",
      "username": "john_doe"
    }
  },
  "from": "system",
  "to": "",
  "timestamp": "2026-01-06T10:00:00Z"
}
```

### Chat Message
```json
{
  "type": "chat",
  "data": {
    "message": "Hello, World!"
  },
  "from": "user_123",
  "to": "user_456",
  "timestamp": "2026-01-06T10:00:00Z"
}
```

## Client Connection Example

### JavaScript (Browser)
```javascript
const token = 'your_jwt_token';
const ws = new WebSocket('ws://localhost:8080/api/v1/ws', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
  
  switch(message.type) {
    case 'notification':
      showNotification(message.data);
      break;
    case 'system':
      handleSystemEvent(message.data);
      break;
    case 'chat':
      displayMessage(message.data);
      break;
  }
};

ws.onerror = (error) => {
  console.error('WebSocket error:', error);
};

ws.onclose = () => {
  console.log('Disconnected from WebSocket');
};

// Send message
ws.send(JSON.stringify({
  type: 'chat',
  data: { message: 'Hello!' }
}));
```

### Go Client
```go
package main

import (
    "log"
    "github.com/gorilla/websocket"
)

func main() {
    headers := http.Header{}
    headers.Add("Authorization", "Bearer your_jwt_token")
    
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/api/v1/ws", headers)
    if err != nil {
        log.Fatal("dial:", err)
    }
    defer conn.Close()

    // Read messages
    go func() {
        for {
            _, message, err := conn.ReadMessage()
            if err != nil {
                log.Println("read:", err)
                return
            }
            log.Printf("received: %s", message)
        }
    }()

    // Send message
    msg := map[string]interface{}{
        "type": "chat",
        "data": map[string]string{"message": "Hello!"},
    }
    
    if err := conn.WriteJSON(msg); err != nil {
        log.Println("write:", err)
    }
    
    select {}
}
```

## Usage in Application

### Broadcasting to All Clients
```go
// In your handler or usecase
wsUseCase.BroadcastNotification(
    "System Alert",
    "Server maintenance in 10 minutes",
    "warning",
)
```

### Send Message to Specific User
```go
// Send notification to specific user
wsUseCase.SendNotificationToUser(
    userID,
    "New Message",
    "You have a new message from John",
    "info",
)
```

### Broadcast System Event
```go
// Broadcast system event
wsUseCase.BroadcastSystemEvent("user_online", map[string]interface{}{
    "user_id": "123",
    "username": "john_doe",
})
```

### Example: Send WebSocket Notification on User Creation
```go
// In user_usecase.go
func (uc *UserUseCase) Create(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {
    // ... create user logic ...
    
    // Send WebSocket notification to admins
    wsUseCase.BroadcastSystemEvent("user_created", map[string]interface{}{
        "user_id": user.ID,
        "email": user.Email,
    })
    
    return user, nil
}
```

## Testing WebSocket

### Using wscat (CLI Tool)
```bash
# Install wscat
npm install -g wscat

# Connect with authentication
wscat -c "ws://localhost:8080/api/v1/ws" -H "Authorization: Bearer YOUR_TOKEN"
```

### Using Postman
1. Create new WebSocket request
2. URL: `ws://localhost:8080/api/v1/ws`
3. Headers: Add `Authorization: Bearer YOUR_TOKEN`
4. Click Connect
5. Send JSON messages

## Configuration

### CORS Settings
Update `CheckOrigin` in `websocket_handler.go`:
```go
var upgrader = gorilla.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // Allow specific origins
        origin := r.Header.Get("Origin")
        return origin == "https://yourdomain.com"
    },
}
```

### Message Size Limit
Update in `client.go`:
```go
const maxMessageSize = 512 * 1024 // 512KB
```

### Timeouts
```go
const (
    writeWait  = 10 * time.Second  // Write timeout
    pongWait   = 60 * time.Second  // Read timeout
    pingPeriod = 54 * time.Second  // Ping interval
)
```

## Best Practices

1. **Always authenticate**: WebSocket endpoint requires valid JWT token
2. **Handle reconnection**: Implement reconnection logic in client
3. **Message validation**: Validate message structure before processing
4. **Error handling**: Handle connection errors gracefully
5. **Resource cleanup**: Ensure connections are properly closed
6. **Rate limiting**: Consider implementing rate limiting for messages
7. **Message size**: Keep messages small for better performance

## Monitoring

Get connection statistics:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/v1/ws/stats
```

## Security Considerations

- Authentication required via JWT
- Token validation on connection
- Configure CORS properly for production
- Implement rate limiting for message sending
- Validate and sanitize all incoming messages
- Use WSS (WebSocket Secure) in production

## Complete Usage Examples

### 1. Broadcasting Notification to All Connected Clients
```go
func NotifyAllUsers(wsUseCase *usecase.WebSocketUseCase) {
    wsUseCase.BroadcastNotification(
        "System Announcement",
        "Server maintenance scheduled for tonight at 10 PM",
        "warning",
    )
}
```

### 2. Send Notification to Specific User
```go
func NotifyUser(wsUseCase *usecase.WebSocketUseCase, userID string) {
    wsUseCase.SendNotificationToUser(
        userID,
        "Welcome!",
        "Thank you for joining our platform",
        "success",
    )
}
```

### 3. Broadcast System Event (e.g., User Online Status)
```go
func BroadcastUserOnline(wsUseCase *usecase.WebSocketUseCase, userID, username string) {
    wsUseCase.BroadcastSystemEvent("user_online", map[string]interface{}{
        "user_id": userID,
        "username": username,
        "timestamp": time.Now(),
    })
}
```

### 4. Integration: User Creation with WebSocket Notification
```go
type UserUseCase struct {
    userRepo      domain.UserRepository
    roleRepo      domain.RoleRepository
    wsUseCase     *WebSocketUseCase  // Add this field
    contextTimeout time.Duration
}

func (uc *UserUseCase) Create(ctx context.Context, req *dto.CreateUserRequest) (*entity.User, error) {
    // ... existing user creation logic ...

    // Send notification to all admins about new user
    go uc.wsUseCase.BroadcastSystemEvent("user_created", map[string]interface{}{
        "user_id": user.ID,
        "email": user.Email,
        "created_at": user.CreatedAt,
    })

    return user, nil
}
```

### 5. Integration: Notify User on Profile Update
```go
func (uc *UserUseCase) Update(ctx context.Context, id string, req *dto.UpdateUserRequest) error {
    // ... existing update logic ...

    // Notify the user about the update
    go uc.wsUseCase.SendNotificationToUser(
        id,
        "Profile Updated",
        "Your profile has been successfully updated",
        "success",
    )

    return nil
}
```

### 6. Real-time Chat Implementation
```go
func SendChatMessage(hub *Hub, fromUserID, toUserID, message string) error {
    msg := NewMessage(
        MessageTypeChat,
        map[string]interface{}{
            "message": message,
            "timestamp": time.Now(),
        },
        fromUserID,
        toUserID,
    )

    return hub.BroadcastToUserJSON(toUserID, msg)
}
```

### 7. Get WebSocket Statistics
```go
func GetConnectionStats(wsUseCase *usecase.WebSocketUseCase) int {
    return wsUseCase.GetConnectedClientsCount()
}
```

### 8. Custom Message Broadcasting
```go
func BroadcastCustomMessage(wsUseCase *usecase.WebSocketUseCase, data interface{}) error {
    return wsUseCase.BroadcastMessage(data)
}
```

### 9. Notify Admins About Critical Events
```go
func NotifyAdmins(wsUseCase *usecase.WebSocketUseCase, adminIDs []string, title, message string) {
    for _, adminID := range adminIDs {
        go wsUseCase.SendNotificationToUser(
            adminID,
            title,
            message,
            "error",
        )
    }
}
```

### 10. Broadcast Product Update to All Users
```go
func BroadcastProductUpdate(wsUseCase *usecase.WebSocketUseCase, productID, productName string, price float64) {
    wsUseCase.BroadcastSystemEvent("product_updated", map[string]interface{}{
        "product_id": productID,
        "product_name": productName,
        "price": price,
        "updated_at": time.Now(),
    })
}
```

## Troubleshooting

### Connection Failed
- Verify JWT token is valid
- Check if token is passed in Authorization header
- Ensure WebSocket endpoint is accessible

### Messages Not Received
- Verify client is properly listening for messages
- Check if connection is still alive (ping/pong)
- Verify message format is correct

### High Memory Usage
- Check for connection leaks
- Ensure connections are properly closed
- Monitor client count via stats endpoint

## Future Enhancements

- [ ] Add room/channel support
- [ ] Implement message persistence
- [ ] Add typing indicators
- [ ] Support binary messages
- [ ] Add compression support
- [ ] Implement backpressure handling
- [ ] Add metrics and monitoring
- [ ] Support horizontal scaling with Redis pub/sub

