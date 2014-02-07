package data

import (
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestWhitelistRecord verifies that WhitelistRecord creation, methods, save, load, and delete work properly
func TestWhitelistRecord(t *testing.T) {
	log.Println("TestWhitelistRecord()")

	// Load config
	config, err := common.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	common.Static.Config = config

	// Generate mock WhitelistRecord
	whitelist := WhitelistRecord{
		Client: "goat_test",
		Approved: true,
	}

	// Save mock whitelist
	if err := whitelist.Save(); err != nil {
		t.Fatalf("Failed to save mock whitelist: %s", err.Error())
	}

	// Load mock whitelist to fetch ID
	whitelist, err = whitelist.Load(whitelist.Client, "client")
	if whitelist == (WhitelistRecord{}) || err != nil {
		t.Fatalf("Failed to load mock whitelist: %s", err.Error())
	}

	// Delete mock whitelist
	if err := whitelist.Delete(); err != nil {
		t.Fatalf("Failed to delete mock whitelist: %s", err.Error())
	}
}
