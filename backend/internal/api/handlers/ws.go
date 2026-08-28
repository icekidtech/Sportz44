package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WSHub is the interface the handler needs from the hub.
type WSHub interface {
	Register(conn *websocket.Conn)
}

// WSHandler handles WebSocket connections.
type WSHandler struct {
	hub            WSHub
	allowedOrigins map[string]bool
}

// NewWSHandler creates a new WSHandler. allowedOrigins should be the same
// list from the ALLOWED_ORIGINS env var; a wildcard "*" allows all origins.
func NewWSHandler(hub WSHub, allowedOrigins []string) *WSHandler {
	origins := make(map[string]bool, len(allowedOrigins))
	for _, o := range allowedOrigins {
		origins[strings.TrimSpace(o)] = true
	}
	return &WSHandler{hub: hub, allowedOrigins: origins}
}

// Handle upgrades the connection and registers it with the hub.
func (h *WSHandler) Handle(c *gin.Context) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if h.allowedOrigins["*"] {
				return true
			}
			origin := r.Header.Get("Origin")
			return origin != "" && h.allowedOrigins[origin]
		},
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	h.hub.Register(conn)
}
