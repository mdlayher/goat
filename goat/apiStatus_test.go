package goat

import (
	"encoding/json"
	"log"
	"testing"
)

// TestGetStatusJSON verifies that /api/status returns proper JSON output
func TestGetStatusJSON(t *testing.T) {
	log.Println("TestGetStatusJSON()")

	// Get status directly
	status := getServerStatus()

	// Request output JSON from API for this status
	var status2 serverStatus
	err := json.Unmarshal(getStatusJSON(), &status2)
	if err != nil {
		t.Fatalf("Failed to unmarshal result JSON")
	}

	// Verify objects are the same
	if status.PID != status2.PID {
		t.Fatalf("PID, expected %d, got %d", status.PID, status2.PID)
	}

	if status.Hostname != status2.Hostname {
		t.Fatalf("Hostname, expected %s, got %s", status.Hostname, status2.Hostname)
	}

	if status.Platform != status2.Platform {
		t.Fatalf("Platform, expected %s, got %s", status.Platform, status2.Platform)
	}

	if status.Architecture != status2.Architecture {
		t.Fatalf("Architecture, expected %s, got %s", status.Architecture, status2.Architecture)
	}

	if status.NumCPU != status2.NumCPU {
		t.Fatalf("NumCPU, expected %d, got %d", status.NumCPU, status2.NumCPU)
	}
}
