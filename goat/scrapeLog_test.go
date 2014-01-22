package goat

import (
	"log"
	"net/url"
	"testing"
)

// TestScrapeLog verifies that scrape log creation, save, load, and delete work properly
func TestScrapeLog(t *testing.T) {
	log.Println("TestScrapeLog()")

	// Load config
	config := loadConfig()
	static.Config = config

	// Generate fake scrape query
	query := url.Values{}
	query.Set("info_hash", "deadbeef")
	query.Set("ip", "127.0.0.1")

	// Generate struct from query
	scrape := new(scrapeLog).FromValues(query)

	// Verify proper hex encode of info hash
	if scrape.InfoHash != "6465616462656566" {
		t.Fatalf("InfoHash, expected \"6465616462656566\", got %s", scrape.InfoHash)
	}

	// Verify same IP
	if scrape.IP != "127.0.0.1" {
		t.Fatalf("IP, expected \"127.0.0.1\", got %s", scrape.IP)
	}

	// Verify scrape can be saved
	if !scrape.Save() {
		t.Fatalf("Failed to save scrapeLog")
	}

	// Verify scrape can be loaded using hex info hash
	scrape2 := scrape.Load("6465616462656566", "info_hash")
	if scrape2 == (scrapeLog{}) {
		t.Fatal("Failed to load scrapeLog")
	}

	// Verify scrape is the same as previous one
	if scrape.InfoHash != scrape2.InfoHash {
		t.Fatalf("scrape.InfoHash, expected %s, got %s", scrape.InfoHash, scrape2.InfoHash)
	}

	// Verify scrape can be deleted
	if !scrape2.Delete() {
		t.Fatalf("Failed to delete scrapeLog")
	}
}
