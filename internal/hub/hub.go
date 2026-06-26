package hub

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Conn      *websocket.Conn
	Role      string // "admin" | "visitor"
	SessionID string // untuk visitor chat
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
}

var H = &Hub{
	clients: make(map[*Client]bool),
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) BroadcastToAdmins(payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, _ := json.Marshal(payload)
	for c := range h.clients {
		if c.Role == "admin" {
			_ = c.Conn.WriteMessage(1, data)
		}
	}
}

func (h *Hub) BroadcastToAll(payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, _ := json.Marshal(payload)
	for c := range h.clients {
		_ = c.Conn.WriteMessage(1, data)
	}
}

func (h *Hub) BroadcastToSession(sessionID string, payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, _ := json.Marshal(payload)
	for c := range h.clients {
		// Kirim ke visitor yang punya session ini
		if c.SessionID == sessionID && c.Role == "visitor" {
			_ = c.Conn.WriteMessage(1, data)
		}
	}
}

func (h *Hub) BroadcastToAdminSession(sessionID string, payload any) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	data, _ := json.Marshal(payload)
	for c := range h.clients {
		if c.Role == "admin" && c.SessionID == sessionID {
			_ = c.Conn.WriteMessage(1, data)
		}
	}
}
