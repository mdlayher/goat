package api

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
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
	config := common.LoadConfig()
	common.Static.Config = config

	// Generate mock user
	user := data.UserRecord{
		Username:     "test",
		Passkey:      "abcdef0123456789",
		TorrentLimit: 10,
	}

	// Save mock user
	if err := user.Save(); err != nil {
		t.Fatalf("Failed to save mock user: %s", err.Error())
	}

	// Load mock user to fetch ID
	user, err := user.Load(user.Username, "username")
	if user == (data.UserRecord{}) || err != nil {
		t.Fatalf("Failed to load mock user: %s", err.Error())
	}

	// Generate an API key and salt
	pass := "deadbeef"
	salt := "test"

	sha := sha1.New()
	sha.Write([]byte(pass + salt))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Generate mock API key
	key := data.APIKey{
		UserID: user.ID,
		Key:    hash,
		Salt:   salt,
	}

	// Save mock API key
	if err := key.Save(); err != nil {
		t.Fatalf("Failed to save mock data.APIKey: %s", err.Error())
	}

	// Load mock data.APIKey to fetch ID
	key, err = key.Load(key.Key, "key")
	if err != nil || key == (data.APIKey{}) {
		t.Fatalf("Failed to load mock data.APIKey: %s", err.Error())
	}

	// Generate mock HTTP request
	r := http.Request{}
	headers := map[string][]string{
		"Authorization": {"Basic " + base64.URLEncoding.EncodeToString([]byte(user.Username+":"+pass))},
	}
	r.Header = headers

	// Perform authentication request
	auth := new(BasicAuthenticator).Auth(&r)
	if !auth {
		t.Fatalf("Failed to authenticate using BasicAuthenticator")
	}

	// Delete mock user
	if err := user.Delete(); err != nil {
		t.Fatalf("Failed to delete mock user: %s", err.Error())
	}

	// Delete mock API key
	if err := key.Delete(); err != nil {
		t.Fatalf("Failed to delete mock data.APIKey: %s", err.Error())
	}
}
