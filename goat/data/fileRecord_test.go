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

	// Verify completed is functional
	if _, err := file.Completed(); err != nil {
		t.Fatalf("Failed to fetch file completed")
	}

	// Verify seeders is functional
	if _, err := file.Seeders(); err != nil {
		t.Fatalf("Failed to fetch file seeders")
	}

	// Verify leechers is functional
	if _, err := file.Leechers(); err != nil {
		t.Fatalf("Failed to fetch file leechers")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}
