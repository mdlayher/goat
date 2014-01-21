package goat

import (
	"log"
	"testing"
)

// TestApiKey verifies that apiKey creation, save, load, and delete work properly
func TestApiKey(t *testing.T) {
	log.Println("TestApiKey()")

	// Load config
	config := loadConfig()
	static.Config = config

	// Generate mock apiKey
	key := apiKey{
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
	if key2 == (apiKey{}) {
		t.Fatal("Failed to load apiKey")
	}

	// Verify key is the same as previous one
	if key.Salt != key2.Salt {
		t.Fatalf("key.Salt, expected %s, got %s", key.Salt, key2.Salt)
	}

	// Verify key can be deleted
	if !key2.Delete() {
		t.Fatalf("Failed to delete apiKey")
	}
}
