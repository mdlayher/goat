package api

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// TestPostLogin verifies that /api/login returns proper JSON output
func TestPostLogin(t *testing.T) {
	log.Println("TestPostLogin()")

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

	// Perform login request for this user
	res, err := postLogin(user)
	if err != nil {
		t.Fatalf("Failed to retrieve login JSON: %s", err.Error())
	}
	log.Println(string(res))

	// Unmarshal output JSON
	var key data.JSONAPIKey
	err = json.Unmarshal(res, &key)
	if err != nil {
		t.Fatalf("Failed to unmarshal login JSON: %s", err.Error())
	}

	// Verify same user ID from API
	if user.ID != key.UserID {
		t.Fatalf("Mismatched user IDs, got %d, expected %d", key.UserID, user.ID)
	}

	// Delete mock user
	if err := user.Delete(); err != nil {
		t.Fatalf("Failed to delete mock user: %s", err.Error())
	}
}
