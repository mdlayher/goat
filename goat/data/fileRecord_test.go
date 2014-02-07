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
	config, err := common.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	common.Static.Config = config

	// Generate mock FileRecord
	file := FileRecord{
		InfoHash: "deadbeef",
		Verified: true,
	}

	// Save mock file
	if err := file.Save(); err != nil {
		t.Fatalf("Failed to save mock file: %s", err.Error())
	}

	// Load mock file to fetch ID
	file, err = file.Load(file.InfoHash, "info_hash")
	if file == (FileRecord{}) || err != nil {
		t.Fatalf("Failed to load mock file: %s", err.Error())
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
	if err := file.Delete(); err != nil {
		t.Fatalf("Failed to delete mock file: %s", err.Error())
	}
}
