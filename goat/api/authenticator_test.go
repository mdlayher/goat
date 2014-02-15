package api

import (
	"encoding/base64"
	"log"
	"net/http"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// TestBasicAuthenticator verifies that BasicAuthenticator.Auth works properly
func TestBasicAuthenticator(t *testing.T) {
	log.Println("TestBasicAuthenticator()")

	// Load config
	config, err := common.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	common.Static.Config = config

	// Generate mock user
	user := new(data.UserRecord)
	if err := user.Create("test", "test", 10); err != nil {
		t.Fatalf("Failed to create mock user: %s", err.Error())
	}

	// Save mock user
	if err := user.Save(); err != nil {
		t.Fatalf("Failed to save mock user: %s", err.Error())
	}

	// Generate mock HTTP request
	r := http.Request{}
	headers := map[string][]string{
		"Authorization": {"Basic " + base64.URLEncoding.EncodeToString([]byte("test:test"))},
	}
	r.Header = headers

	// Perform authentication request
	clientErr, serverErr := new(BasicAuthenticator).Auth(&r)
	if clientErr != nil {
		t.Fatalf("Failed to authenticate using BasicAuthenticator: client: %s", clientErr.Error())
	}

	if serverErr != nil {
		t.Fatalf("Failed to authenticate using BasicAuthenticator: server: %s", serverErr.Error())
	}

	// Delete mock user
	if err := user.Delete(); err != nil {
		t.Fatalf("Failed to delete mock user: %s", err.Error())
	}
}
