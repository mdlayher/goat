package goat

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"testing"
)

// TestBasicAPIAuthenticator verifies that basicAPIAuthenticator.Auth works properly
func TestBasicAPIAuthenticator(t *testing.T) {
	log.Println("TestBasicAPIAuthenticator()")

	// Load config
	config := loadConfig()
	static.Config = config

	// Generate mock user
	user := userRecord{
		Username:     "test",
		Passkey:      "abcdef0123456789",
		TorrentLimit: 10,
	}

	// Save mock user
	if !user.Save() {
		t.Fatalf("Failed to save mock user")
	}

	// Load mock user to fetch ID
	user = user.Load(user.Username, "username")
	if user == (userRecord{}) {
		t.Fatalf("Failed to load mock user")
	}

	// Generate an API key and salt
	pass := "deadbeef"
	salt := "test"

	sha := sha1.New()
	sha.Write([]byte(pass + salt))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Generate mock API key
	key := apiKey{
		UserID: user.ID,
		Key:    hash,
		Salt:   salt,
	}

	// Save mock API key
	if !key.Save() {
		t.Fatalf("Failed to save mock apiKey")
	}

	// Load mock apiKey to fetch ID
	key = key.Load(key.Key, "key")
	if key == (apiKey{}) {
		t.Fatalf("Failed to load mock apiKey")
	}

	// Generate mock HTTP request
	r := http.Request{}
	headers := map[string][]string{
		"Authorization": []string{"Basic " + base64.URLEncoding.EncodeToString([]byte(user.Username+":"+pass))},
	}
	r.Header = headers

	// Perform authentication request
	auth := new(basicAPIAuthenticator).Auth(&r)
	if !auth {
		t.Fatalf("Failed to authenticate using basicAPIAuthenticator")
	}

	// Delete mock user
	if !user.Delete() {
		t.Fatalf("Failed to delete mock user")
	}

	// Delete mock API key
	if !key.Delete() {
		t.Fatalf("Failed to delete mock apiKey")
	}
}
