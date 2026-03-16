package websocket

import "testing"

func TestHubCreation(t *testing.T) {

	hub := NewHub()

	if hub == nil {
		t.Fatal("hub should not be nil")
	}

	if hub.Clients == nil {
		t.Fatal("clients map not initialized")
	}
}
