package goat

import (
	"log"
	"testing"
	"time"
)

// TestFileUserRecord verifies that fileUserRecord creation, methods, save, load, and delete work properly
func TestFileUserRecord(t *testing.T) {
	log.Println("TestFileUserRecord()")

	// Generate mock fileUserRecord
	fileUser := fileUserRecord{
		FileID:     1,
		UserID:     1,
		IP:         "127.0.0.1",
		Active:     true,
		Completed:  true,
		Announced:  10,
		Uploaded:   5000,
		Downloaded: 5000,
		Left:       0,
		Time:       time.Now().Unix(),
	}

	// Save mock fileUser
	if !fileUser.Save() {
		t.Fatalf("Failed to save mock fileUser")
	}

	// Load mock fileUser
	fileUser = fileUser.Load(fileUser.FileID, fileUser.UserID, fileUser.IP)
	if fileUser == (fileUserRecord{}) {
		t.Fatalf("Failed to load mock fileUser")
	}

	// Delete mock fileUser
	if !fileUser.Delete() {
		t.Fatalf("Failed to delete mock fileUser")
	}
}
