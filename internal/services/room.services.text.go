package services

import "testing"

func TestFindRoomsByName(t *testing.T) {

	rooms, err := FindRoomsByName("gen")

	if err != nil {
		t.Fatalf("error finding rooms: %v", err)
	}

	if rooms == nil {
		t.Fatal("expected rooms result")
	}
}
