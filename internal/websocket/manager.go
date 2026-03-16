package websocket

import "sync"

// Manager tracks which clients are subscribed to which rooms.
// All methods are safe for concurrent use.
type Manager struct {
	mu    sync.RWMutex
	Rooms map[int]map[*Client]bool
}

func NewManager() *Manager {
	return &Manager{
		Rooms: make(map[int]map[*Client]bool),
	}
}

func (m *Manager) JoinRoom(roomID int, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.Rooms[roomID]; !ok {
		m.Rooms[roomID] = make(map[*Client]bool)
	}
	m.Rooms[roomID][client] = true
	client.RoomID = roomID
}

func (m *Manager) LeaveRoom(roomID int, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if clients, ok := m.Rooms[roomID]; ok {
		delete(clients, client)
		if len(clients) == 0 {
			delete(m.Rooms, roomID)
		}
	}
}

func (m *Manager) Broadcast(roomID int, message []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.Rooms[roomID] {
		client.Send <- message
	}
}
