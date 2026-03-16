package websocket_test

import (
	"testing"

	ws "github.com/Deepsayan-Das/chatter_GO/internal/websocket"
)

// ── Manager ───────────────────────────────────────────────────────────────────

func newTestClient(id int) *ws.Client {
	return &ws.Client{
		ID:   id,
		Send: make(chan []byte, 8), // buffered so Broadcast doesn't block in tests
	}
}

func TestManager_JoinRoom_CreatesRoom(t *testing.T) {
	m := ws.NewManager()
	c := newTestClient(1)

	m.JoinRoom(10, c)

	if c.RoomID != 10 {
		t.Errorf("client.RoomID: got %d, want 10", c.RoomID)
	}
}

func TestManager_JoinRoom_MultipleClients(t *testing.T) {
	m := ws.NewManager()
	c1 := newTestClient(1)
	c2 := newTestClient(2)

	m.JoinRoom(10, c1)
	m.JoinRoom(10, c2)

	m.Broadcast(10, []byte("hello"))

	for _, c := range []*ws.Client{c1, c2} {
		select {
		case msg := <-c.Send:
			if string(msg) != "hello" {
				t.Errorf("client %d: got %q, want \"hello\"", c.ID, msg)
			}
		default:
			t.Errorf("client %d did not receive broadcast message", c.ID)
		}
	}
}

func TestManager_LeaveRoom_RemovesClient(t *testing.T) {
	m := ws.NewManager()
	c := newTestClient(1)

	m.JoinRoom(5, c)
	m.LeaveRoom(5, c)

	// After leaving, a broadcast should not reach the client
	m.Broadcast(5, []byte("should not arrive"))
	select {
	case msg := <-c.Send:
		t.Errorf("client received message after leaving room: %q", msg)
	default:
		// expected — channel should have no message
	}
}

func TestManager_LeaveRoom_PrunesEmptyRoom(t *testing.T) {
	m := ws.NewManager()
	c := newTestClient(1)

	m.JoinRoom(7, c)
	m.LeaveRoom(7, c)

	// Broadcast to a pruned room should be a no-op (not panic)
	m.Broadcast(7, []byte("noop"))
}

func TestManager_Broadcast_UnknownRoom(t *testing.T) {
	m := ws.NewManager()
	// Must not panic when broadcasting to a room with no members
	m.Broadcast(999, []byte("noop"))
}

// ── Hub ───────────────────────────────────────────────────────────────────────

func TestNewHub_Initialized(t *testing.T) {
	h := ws.NewHub()
	if h.Clients == nil {
		t.Error("Hub.Clients should not be nil")
	}
	if h.UserClients == nil {
		t.Error("Hub.UserClients should not be nil")
	}
	if h.Manager == nil {
		t.Error("Hub.Manager should not be nil")
	}
}

func TestHub_SendToUser_UnknownUser(t *testing.T) {
	h := ws.NewHub()
	// SendToUser for an offline user should return false without panicking
	if h.SendToUser(9999, []byte("hello")) {
		t.Error("expected false for a user not registered in the hub")
	}
}
