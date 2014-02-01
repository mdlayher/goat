package data

import (
	"log"
	"net/url"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestAnnounceLog verifies that announce log creation, save, load, and delete work properly
func TestAnnounceLog(t *testing.T) {
	log.Println("TestAnnounceLog()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate fake announce query
	query := url.Values{}
	query.Set("info_hash", "deadbeef")
	query.Set("ip", "127.0.0.1")
	query.Set("port", "5000")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("left", "0")

	// Generate struct from query
	announce := new(AnnounceLog).FromValues(query)

	// Verify proper hex encode of info hash
	if announce.InfoHash != "6465616462656566" {
		t.Fatalf("InfoHash, expected \"6465616462656566\", got %s", announce.InfoHash)
	}

	// Verify same IP
	if announce.IP != "127.0.0.1" {
		t.Fatalf("IP, expected \"127.0.0.1\", got %s", announce.IP)
	}

	// Verify proper port integer conversion
	if announce.Port != 5000 {
		t.Fatalf("Port, expected 5000, got %d", announce.Port)
	}

	// Verify proper uploaded integer conversion
	if announce.Uploaded != 0 {
		t.Fatalf("Uploaded, expected 0, got %d", announce.Uploaded)
	}

	// Verify proper downloaded integer conversion
	if announce.Downloaded != 0 {
		t.Fatalf("Downloaded, expected 0, got %d", announce.Downloaded)
	}

	// Verify proper left integer conversion
	if announce.Left != 0 {
		t.Fatalf("Left, expected 0, got %d", announce.Left)
	}

	// Verify announce can be saved
	if !announce.Save() {
		t.Fatalf("Failed to save AnnounceLog")
	}

	// Verify announce can be loaded using hex info hash
	announce2 := announce.Load("6465616462656566", "info_hash")
	if announce2 == (AnnounceLog{}) {
		t.Fatal("Failed to load AnnounceLog")
	}

	// Verify announce is the same as previous one
	if announce.InfoHash != announce2.InfoHash {
		t.Fatalf("announce.InfoHash, expected %s, got %s", announce.InfoHash, announce2.InfoHash)
	}

	// Verify announce can be deleted
	if !announce2.Delete() {
		t.Fatalf("Failed to delete AnnounceLog")
	}
}
