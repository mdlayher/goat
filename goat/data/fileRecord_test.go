package data

import (
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestFileRecord verifies that FileRecord creation, methods, save, load, and delete work properly
func TestFileRecord(t *testing.T) {
	log.Println("TestFileRecord()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock FileRecord
	file := FileRecord{
		InfoHash: "deadbeef",
		Verified: true,
	}

	// Save mock file
	if !file.Save() {
		t.Fatalf("Failed to save mock file")
	}

	// Load mock file to fetch ID
	file = file.Load(file.InfoHash, "info_hash")
	if file == (FileRecord{}) {
		t.Fatalf("Failed to load mock file")
	}

	// Verify positive number of completed on file
	if file.Completed() < 0 {
		t.Fatalf("Failed to fetch file completed")
	}

	// Verify positive number of seeders on file
	if file.Seeders() < 0 {
		t.Fatalf("Failed to fetch file seeders")
	}

	// Verify positive number of leechers on file
	if file.Leechers() < 0 {
		t.Fatalf("Failed to fetch file leechers")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}
