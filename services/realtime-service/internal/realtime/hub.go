package realtime

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[uuid.UUID]map[*Client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]map[*Client]struct{}),
	}
}

func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[client.UserID] == nil {
		h.clients[client.UserID] = make(map[*Client]struct{})
	}

	h.clients[client.UserID][client] = struct{}{}
}

func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections, ok := h.clients[client.UserID]
	if !ok {
		return
	}

	delete(connections, client)

	if len(connections) == 0 {
		delete(h.clients, client.UserID)
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, message Message) {
	payload, err := json.Marshal(message)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients[userID] {
		client.Send(payload)
	}
}

func (h *Hub) Stats() (users int, connections int) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users = len(h.clients)

	for _, conns := range h.clients {
		connections += len(conns)
	}

	return users, connections
}

func (h *Hub) CloseAll() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, connections := range h.clients {
		for client := range connections {
			client.Close()
		}
	}

	h.clients = make(map[uuid.UUID]map[*Client]struct{})
}
