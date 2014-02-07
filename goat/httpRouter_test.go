package goat

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// Table driven tests to iterate over and test the main HTTP router
var httpTests = []struct {
	url string
}{
	{"/test"},
	{"/announce"},
	{"/announce?info_hash=deadbeef&ip=127.0.0.1&port=abc&uploaded=0&downloaded=0&left=10"},
	{"/announce?info_hash=deadbeef&ip=127.0.0.1&port=5000&uploaded=0&downloaded=0&left=10"},
	{"/announce?info_hash=deadbeef&ip=127.0.0.1&port=5000&uploaded=0&downloaded=0&left=10&compact=1"},
	{"/scrape"},
	{"/scrape?info_hash=deadbeef"},
	{"/scrape?info_hash=deadbeef&info_hash=beefdead"},
}

// TestHTTPRouter verifies that the main HTTP router is working properly
func TestHTTPRouter(t *testing.T) {
	log.Println("TestHTTPRouter()")

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

	// Generate mock data.FileRecord
	file2 := data.FileRecord{
		InfoHash: "6265656664656164",
		Verified: true,
	}

	// Save mock file
	if !file2.Save() {
		t.Fatalf("Failed to save mock file")
	}

	// Iterate all HTTP tests
	for _, test := range httpTests {
		// Generate mock HTTP request
		r, err := http.NewRequest("GET", "http://localhost:8080"+test.url, nil)
		r.Header.Set("User-Agent", "goat_test")
		if err != nil {
			t.Fatalf("Failed to create HTTP request")
		}

		// Capture HTTP writer response with recorder
		w := httptest.NewRecorder()

		// Invoke HTTP router
		parseHTTP(w, r)
		log.Println("TEST:", test.url)
		log.Println(w.Body.String())
	}

	// Delete mock file2
	if !file.Delete() || !file2.Delete() {
		t.Fatalf("Failed to delete mock file")
	}
}
