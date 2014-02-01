package tracker

import (
	"bytes"
	"log"
	"net/url"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
	"github.com/mdlayher/goat/goat/data/udp"
)

// TestUDPAnnounce verifies that the UDP tracker announce output format is correct
func TestUDPAnnounce(t *testing.T) {
	log.Println("TestUDPAnnounce()")

	// Load config
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "6465616462656566",
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
	query.Set("numwant", "50")

	// Create a UDP tracker, trigger an announce
	tracker := UDPTracker{TransID: uint32(1234)}
	res := tracker.Announce(query, file)

	// Decode response
	announce, err := new(udp.AnnounceResponse).FromBytes(res)
	if err != nil {
		t.Fatalf("Failed to decode UDP announce response")
	}
	log.Println(announce)

	// Verify correct action
	if announce.Action != 1 {
		t.Fatalf("Incorrect UDP action, expected 1")
	}

	// Encode response, verify same as before
	announceBuf, err := announce.ToBytes()
	if err != nil {
		t.Fatalf("Failed to encode UDP announce response")
	}

	if !bytes.Equal(res, announceBuf) {
		t.Fatalf("Byte slices are not identical")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}

// TestUDPTrackerError verifies that the UDP tracker error format is correct
func TestUDPTrackerError(t *testing.T) {
	log.Println("TestUDPTrackerError()")

	// Create a UDP tracker, trigger an error
	tracker := UDPTracker{TransID: uint32(1234)}
	res := tracker.Error("Testing")

	// Decode response
	errRes, err := new(udp.ErrorResponse).FromBytes(res)
	if err != nil {
		t.Fatalf("Failed to decode UDP error response")
	}
	log.Println(errRes)

	// Verify correct action
	if errRes.Action != 3 {
		t.Fatalf("Incorrect UDP action, expected 3")
	}

	// Verify correct error
	if errRes.Error != "Testing" {
		t.Fatalf("Incorrect UDP error, expected 'Testing'")
	}

	// Encode response, verify same as before
	errResBuf, err := errRes.ToBytes()
	if err != nil {
		t.Fatalf("Failed to encode UDP error response")
	}

	if !bytes.Equal(res, errResBuf) {
		t.Fatalf("Byte slices are not identical")
	}
}

// TestUDPTrackerScrape verifies that the UDP tracker scrape format is correct
func TestUDPTrackerScrape(t *testing.T) {
	log.Println("TestUDPTrackerScrape()")

	// Generate mock data.FileRecord
	file := data.FileRecord{
		InfoHash: "6465616462656566",
		Verified: true,
	}

	// Save mock file
	if !file.Save() {
		t.Fatalf("Failed to save mock file")
	}

	// Store file in slice
	files := make([]data.FileRecord, 0)
	files = append(files[:], file)

	// Create a UDP tracker, trigger a scrape
	tracker := UDPTracker{TransID: uint32(1234)}
	res := tracker.Scrape(files)

	// Decode response
	scrape, err := new(udp.ScrapeResponse).FromBytes(res)
	if err != nil {
		t.Fatalf("Failed to decode UDP scrape response")
	}
	log.Println(scrape)

	// Verify correct action
	if scrape.Action != 2 {
		t.Fatalf("Incorrect UDP action, expected 2")
	}

	// Encode response, verify same as before
	scrapeBuf, err := scrape.ToBytes()
	if err != nil {
		t.Fatalf("Failed to encode UDP scrape response")
	}

	if !bytes.Equal(res, scrapeBuf) {
		t.Fatalf("Byte slices are not identical")
	}

	// Delete mock file
	if !file.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}
