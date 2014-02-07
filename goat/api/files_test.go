package api

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// TestGetFilesJSON verifies that /api/files returns proper JSON output
func TestGetFilesJSON(t *testing.T) {
	log.Println("TestGetFilesJSON()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "deadbeef",
		Verified: true,
	}

	// Save mock file
	if err := file.Save(); err != nil {
		t.Fatalf("Failed to save mock file: %s", err.Error())
	}

	// Load mock file to fetch ID
	file, err := file.Load(file.InfoHash, "info_hash")
	if file == (data.FileRecord{}) || err != nil {
		t.Fatalf("Failed to load mock file: %s", err.Error())
	}

	// Request output JSON from API for this file
	res, err := getFilesJSON(file.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve files JSON: %s", err.Error())
	}

	// Unmarshal output JSON
	var file2 data.FileRecord
	err = json.Unmarshal(res, &file2)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON for single file: %s", err.Error())
	}

	// Verify objects are the same
	if file.ID != file2.ID {
		t.Fatalf("ID, expected %d, got %d", file.ID, file2.ID)
	}

	// Request output JSON from API for all files
	res, err = getFilesJSON(-1)
	if err != nil {
		t.Fatalf("Failed to retrieve all files JSON: %s", err.Error())
	}

	// Unmarshal all output JSON
	var allFiles []data.FileRecord
	err = json.Unmarshal(res, &allFiles)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON for all files: %s", err.Error())
	}

	// Verify known file is in result set
	found := false
	for _, f := range allFiles {
		if f.ID == file.ID {
			found = true
		}
	}

	if !found {
		t.Fatalf("Expected file not found in all files result set")
	}

	// Delete mock file
	if err := file.Delete(); err != nil {
		t.Fatalf("Failed to delete mock file: %s", err.Error())
	}
}
