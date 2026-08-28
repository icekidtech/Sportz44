package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

// MatchEventsChannel is the Redis pub/sub channel for live match updates.
const MatchEventsChannel = "match-events"

// Client represents a single connected WebSocket client.
type Client struct {
	conn *websocket.Conn
	send chan []byte
}

// Hub manages WebSocket clients and fans out messages from Redis pub/sub.
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	rdb     *redis.Client
}

// NewHub creates a Hub and starts a goroutine that subscribes to the
// match-events Redis channel and broadcasts to all connected clients.
func NewHub(rdb *redis.Client) *Hub {
	h := &Hub{
		clients: make(map[*Client]bool),
		rdb:     rdb,
	}
	go h.listen()
	return h
}

// listen subscribes to the Redis match-events channel and broadcasts.
func (h *Hub) listen() {
	ctx := context.Background()
	pubsub := h.rdb.Subscribe(ctx, MatchEventsChannel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		h.broadcast([]byte(msg.Payload))
	}
}

// broadcast sends a message to every connected client.
func (h *Hub) broadcast(data []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients {
		select {
		case c.send <- data:
		default:
			// Client is slow; drop the message rather than block.
		}
	}
}

// Register adds a new client and starts its read/write pumps.
func (h *Hub) Register(conn *websocket.Conn) {
	c := &Client{conn: conn, send: make(chan []byte, 256)}
	h.mu.Lock()
	h.clients[c] = true
	h.mu.Unlock()

	go h.writePump(c)
	go h.readPump(c)
}

// writePump writes queued messages to the client.
func (h *Hub) writePump(c *Client) {
	defer func() {
		_ = c.conn.Close()
		h.unregister(c)
	}()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

// readPump reads (and discards) client messages to detect disconnects.
func (h *Hub) readPump(c *Client) {
	defer func() {
		_ = c.conn.Close()
		h.unregister(c)
	}()
	c.conn.SetReadLimit(512)
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			return
		}
	}
}

// unregister removes a client from the hub.
func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	if _, ok := h.clients[c]; ok {
		delete(h.clients, c)
		close(c.send)
	}
	h.mu.Unlock()
}

// PublishMatchUpdate publishes a live match update to Redis so the hub and
// other consumers can receive it.
func (h *Hub) PublishMatchUpdate(ctx context.Context, update interface{}) error {
	b, err := json.Marshal(update)
	if err != nil {
		return err
	}
	return h.rdb.Publish(ctx, MatchEventsChannel, b).Err()
}
