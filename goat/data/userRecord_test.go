package data

import (
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestUserRecord verifies that user creation, save, load, and delete work properly
func TestUserRecord(t *testing.T) {
	log.Println("TestUserRecord()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Create a user
	user := new(UserRecord)
	if err := user.Create("test", 100); err != nil {
		t.Fatalf("Failed to create UserRecord")
	}

	// Verify proper passkey length
	if len(user.Passkey) != 40 {
		t.Fatalf("user.Passkey is %d characters, expected 40", len(user.Passkey))
	}

	// Verify user can be saved
	if err := user.Save(); err != nil {
		t.Fatalf("Failed to save UserRecord: %s", err.Error())
	}

	// Verify user can be loaded using username
	user2, err := user.Load("test", "username")
	if user2 == (UserRecord{}) || err != nil {
		t.Fatal("Failed to load UserRecord: %s", err.Error())
	}

	// Verify user is the same as previous one
	if user.Passkey != user2.Passkey {
		t.Fatalf("user.Passkey, expected %s, got %s", user.Passkey, user2.Passkey)
	}

	// Verify user can be deleted
	if err := user2.Delete(); err != nil {
		t.Fatalf("Failed to delete UserRecord: %s", err.Error())
	}
}
