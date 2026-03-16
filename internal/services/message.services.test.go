package services

import (
	"testing"
)

func TestCreateMessage(t *testing.T) {

	userID := 1
	roomID := 1
	content := "test message"

	messageID, err := CreateMessage(userID, roomID, content)

	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	if messageID == 0 {
		t.Fatal("message id should not be zero")
	}

}
