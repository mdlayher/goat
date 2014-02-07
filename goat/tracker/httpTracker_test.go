package tracker

import (
	"bytes"
	"log"
	"net/url"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"

	// Import bencode library
	bencode "code.google.com/p/bencode-go"
)

// TestHTTPAnnounce verifies that the HTTP tracker announce output format is correct
func TestHTTPAnnounce(t *testing.T) {
	log.Println("TestHTTPAnnounce()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "6465616462656566303030303030303030303030",
		Verified: true,
	}

	// Save mock file
	if !file.Save() {
		t.Fatalf("Failed to save mock file")
	}

	// Generate fake announce query
	query := url.Values{}
	query.Set("info_hash", "deadbeef")
	query.Set("ip", "127.0.0.1")
	query.Set("port", "5000")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("left", "0")

	// Create a HTTP tracker, trigger an announce
	tracker := HTTPTracker{}
	res := tracker.Announce(query, file)
	log.Println(string(res))

	// Unmarshal response
	announce := AnnounceResponse{}
	if err := bencode.Unmarshal(bytes.NewReader(res), &announce); err != nil {
		t.Fatalf("Failed to unmarshal bencode announce response")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}

// TestHTTPTrackerError verifies that the HTTP tracker error format is correct
func TestHTTPTrackerError(t *testing.T) {
	log.Println("TestHTTPTrackerError()")

	// Create a HTTP tracker, trigger an error
	tracker := HTTPTracker{}
	res := tracker.Error("Testing")
	log.Println(string(res))

	// Unmarshal response
	errRes := errorResponse{}
	if err := bencode.Unmarshal(bytes.NewReader(res), &errRes); err != nil {
		t.Fatalf("Failed to unmarshal bencode error response")
	}
}

// TestHTTPTrackerScrape verifies that the HTTP tracker scrape format is correct
func TestHTTPTrackerScrape(t *testing.T) {
	log.Println("TestHTTPTrackerScrape()")

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "6465616462656566303030303030303030303030",
		Verified: true,
	}

	// Save mock file
	if !file.Save() {
		t.Fatalf("Failed to save mock file")
	}

	// Store file in slice
	files := make([]data.FileRecord, 0)
	files = append(files[:], file)

	// Create a HTTP tracker, trigger a scrape
	tracker := HTTPTracker{}
	res := tracker.Scrape(files)
	log.Println(string(res))

	// Unmarshal response
	scrape := scrapeResponse{}
	if err := bencode.Unmarshal(bytes.NewReader(res), &scrape); err != nil {
		t.Fatalf("Failed to unmarshal bencode scrape response")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}
