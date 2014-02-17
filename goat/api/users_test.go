package api

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// TestGetUsersJSON verifies that /api/users returns proper JSON output
func TestGetUsersJSON(t *testing.T) {
	log.Println("TestGetUsersJSON()")

	// Load config
	config, err := common.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	common.Static.Config = config

	// Generate mock data.UserRecord
	mockUser := new(data.UserRecord)
	if err := mockUser.Create("test", "test", 100); err != nil {
		t.Fatalf("Failed to create mock user: %s", err.Error())
	}

	// Save mock user
	if err := mockUser.Save(); err != nil {
		t.Fatalf("Failed to save mock user: %s", err.Error())
	}

	// Load mock user to fetch ID
	user, err := mockUser.Load(mockUser.Username, "username")
	if user == (data.UserRecord{}) || err != nil {
		t.Fatalf("Failed to load mock user: %s", err.Error())
	}

	// Request output JSON from API for this user
	res, err := getUsersJSON(user.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve users JSON: %s", err.Error())
	}

	// Unmarshal output JSON
	var user2 data.UserRecord
	err = json.Unmarshal(res, &user2)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON for single user: %s", err.Error())
	}

	// Verify objects are the same
	if user.ID != user2.ID {
		t.Fatalf("ID, expected %d, got %d", user.ID, user2.ID)
	}

	// Request output JSON from API for all users
	res, err = getUsersJSON(-1)
	if err != nil {
		t.Fatalf("Failed to retrieve all users JSON: %s", err.Error())
	}

	// Unmarshal all output JSON
	var allUsers []data.JSONUserRecord
	err = json.Unmarshal(res, &allUsers)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON for all users: %s", err.Error())
	}

	// Verify known user is in result set
	found := false
	for _, f := range allUsers {
		if f.ID == user.ID {
			found = true
		}
	}

	if !found {
		t.Fatalf("Expected user not found in all users result set")
	}

	// Delete mock user
	if err := user.Delete(); err != nil {
		t.Fatalf("Failed to delete mock user: %s", err.Error())
	}
}
