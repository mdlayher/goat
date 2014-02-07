package data

import (
	"log"
	"testing"
	"time"

	"github.com/mdlayher/goat/goat/common"
)

// TestFileUserRecord verifies that FileUserRecord creation, methods, save, load, and delete work properly
func TestFileUserRecord(t *testing.T) {
	log.Println("TestFileUserRecord()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock FileUserRecord
	fileUser := FileUserRecord{
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
	if err := fileUser.Save(); err != nil {
		t.Fatalf("Failed to save mock fileUser: %s", err.Error())
	}

	// Load mock fileUser
	fileUser, err := fileUser.Load(fileUser.FileID, fileUser.UserID, fileUser.IP)
	if fileUser == (FileUserRecord{}) || err != nil {
		t.Fatalf("Failed to load mock fileUser: %s", err.Error())
	}

	// Delete mock fileUser
	if err := fileUser.Delete(); err != nil {
		t.Fatalf("Failed to delete mock fileUser: %s", err.Error())
	}
}
