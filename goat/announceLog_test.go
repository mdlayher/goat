package goat

import (
	"net/url"
	"testing"
)

func TestAnnounceLog(t *testing.T) {
	config := loadConfig()
	static.Config = config

	query := url.Values{}
	query.Set("info_hash", "deadbeef")
	query.Set("ip", "127.0.0.1")
	query.Set("port", "5000")
	query.Set("uploaded", "0")
	query.Set("downloaded", "0")
	query.Set("left", "0")

	announce := new(announceLog).FromValues(query)
	if announce.InfoHash != "6465616462656566" {
		t.Fatalf("InfoHash, expected \"6465616462656566\", got %s", announce.InfoHash)
	}

	if announce.IP != "127.0.0.1" {
		t.Fatalf("IP, expected \"127.0.0.1\", got %s", announce.IP)
	}

	if announce.Port != 5000 {
		t.Fatalf("Port, expected 5000, got %d", announce.Port)
	}

	if announce.Uploaded != 0 {
		t.Fatalf("Uploaded, expected 0, got %d", announce.Uploaded)
	}

	if announce.Downloaded != 0 {
		t.Fatalf("Downloaded, expected 0, got %d", announce.Downloaded)
	}

	if announce.Left != 0 {
		t.Fatalf("Left, expected 0, got %d", announce.Left)
	}

	if !announce.Save() {
		t.Fatalf("Failed to save announceLog")
	}

	announce2 := announce.Load("6465616462656566", "info_hash")
	if announce.InfoHash != announce2.InfoHash {
		t.Fatalf("announce.InfoHash, expected %s, got %s", announce.InfoHash, announce2.InfoHash)
	}

	if !announce2.Delete() {
		t.Fatalf("Failed to delete announceLog")
	}
}
