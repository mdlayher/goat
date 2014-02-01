package data

import (
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestAPIKey verifies that APIKey creation, save, load, and delete work properly
func TestAPIKey(t *testing.T) {
	log.Println("TestAPIKey()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock APIKey
	key := APIKey{
		UserID: 1,
		Key:    "deadbeef",
		Salt:   "test",
	}

	// Verify key can be saved
	if !key.Save() {
		t.Fatalf("Failed to save keyLog")
	}

	// Verify key can be loaded using matching key
	key2 := key.Load(key.Key, "key")
	if key2 == (APIKey{}) {
		t.Fatal("Failed to load APIKey")
	}

	// Verify key is the same as previous one
	if key.Salt != key2.Salt {
		t.Fatalf("key.Salt, expected %s, got %s", key.Salt, key2.Salt)
	}

	// Verify key can be deleted
	if !key2.Delete() {
		t.Fatalf("Failed to delete APIKey")
	}
}
