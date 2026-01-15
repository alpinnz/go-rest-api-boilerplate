# WebSocket Implementation

Real-time bidirectional communication server integrated with this boilerplate.

> For auth/token setup and how to obtain a JWT, follow the root `README.md` and Swagger docs (`/docs`).

## Features

- Real-time bidirectional communication
- Hub pattern (central connection manager)
- JWT authentication required
- Broadcast to all clients, or target a specific user
- Ping/pong health checks
- Per-user connection tracking
- Connection statistics endpoint

## Architecture

```
internal/
├── websocket/
│   ├── hub.go          # Central connection manager
│   ├── client.go       # WebSocket client handler
│   └── types.go        # Message type definitions
├── delivery/http/handler/
│   └── websocket_handler.go  # HTTP -> WebSocket upgrade handler
└── usecase/
    └── websocket_usecase.go  # Business logic layer
```

## Endpoints

### Connect

`GET /api/v1/ws`

- Authentication: required (`Authorization: Bearer <token>`)
- Upgrades to a WebSocket connection

### Stats

`GET /api/v1/ws/stats`

- Authentication: required

## Message Types

Messages follow the shape:

```json
{
  "type": "notification",
  "data": {},
  "from": "system",
  "to": "",
  "timestamp": "2026-01-06T10:00:00Z"
}
```

## Client Examples

### JavaScript (Browser)

```javascript
const token = 'your_jwt_token';

// Note: many browser WebSocket clients can't set custom headers.
// If you need headers, use a Node.js client or pass the token via query param (if you support it).
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
};
```

### Go Client

```text
package main

import (
	"net/http"
	"github.com/gorilla/websocket"
)

func main() {
	headers := http.Header{}
	headers.Add("Authorization", "Bearer your_jwt_token")

	_, _, _ = websocket.DefaultDialer.Dial("ws://localhost:8080/api/v1/ws", headers)
}
```

## Testing

Optional CLI testing (requires external tooling):
- `wscat` or Postman can be used to open a WebSocket connection

## Configuration

### Origin Policy

Adjust `CheckOrigin` in `websocket_handler.go` for production.

### Message Size / Timeouts

See constants in `client.go` for:
- max message size
- write timeout
- ping/pong intervals

## Monitoring

Get connection statistics:
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:8080/api/v1/ws/stats
```

## Security Considerations

- Require authentication (already enforced)
- Configure origin policy (CheckOrigin) for production
- Validate/sanitize incoming messages before processing
- Use WSS in production

## Troubleshooting

### Connection Failed
- Ensure token is valid and included in `Authorization` header
- Ensure the endpoint is reachable

### Messages Not Received
- Ensure the client reads from the connection
- Check ping/pong settings and connection lifetime
