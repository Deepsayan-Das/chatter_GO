// hub.go - central connection registry
package websocket

import "sync"

// Hub maintains the set of active clients and routes messages.
type Hub struct {
	mu          sync.RWMutex
	Clients     map[*Client]bool
	UserClients map[int]*Client // user ID -> client (for DM routing)
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan []byte
	Manager     *Manager
}

func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		UserClients: make(map[int]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan []byte),
		Manager:     NewManager(),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.UserClients[client.ID] = client
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				delete(h.UserClients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				client.Send <- message
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser delivers a message directly to a specific connected user.
// Returns false if the target user is not currently online.
func (h *Hub) SendToUser(userID int, message []byte) bool {
	h.mu.RLock()
	client, ok := h.UserClients[userID]
	h.mu.RUnlock()
	if !ok {
		return false
	}
	client.Send <- message
	return true
}
