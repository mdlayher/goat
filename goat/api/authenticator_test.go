package api

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// TestAuthenticator verifies that the Basic and HMAC authenticators work properly
func TestAuthenticator(t *testing.T) {
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

	// Load user to get ID
	user2, err := user.Load("test", "username")
	if err != nil || (user2 == data.UserRecord{}) {
		t.Fatalf("Failed to load mock user: %s", err.Error())
	}

	// Generate mock HTTP request
	r, err := http.NewRequest("POST", "http://localhost:8080/api/login", nil)
	if err != nil {
		t.Fatalf("Failed to generate HTTP request: %s", err.Error())
	}

	headers := map[string][]string{
		"Authorization": {"Basic " + base64.URLEncoding.EncodeToString([]byte("test:test"))},
	}
	r.Header = headers

	// Capture HTTP response with recorder
	w := httptest.NewRecorder()

	// Perform HTTP Basic authentication
	var apiAuth APIAuthenticator
	apiAuth = new(BasicAuthenticator)

	// Attempt Basic authentication
	clientErr, serverErr := apiAuth.Auth(r)
	if clientErr != nil {
		t.Fatalf("Failed to authenticate: client: %s", clientErr.Error())
	}
	if serverErr != nil {
		t.Fatalf("Failed to authenticate: server: %s", serverErr.Error())
	}

	// Invoke API router
	Router(w, r, user2)

	// Read HTTP response body
	body, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read HTTP response body: %s", err.Error())
	}

	// Unmarshal result JSON
	var login data.JSONAPIKey
	if err := json.Unmarshal(body, &login); err != nil {
		t.Fatalf("Failed to unmarshal login JSON: %s", err.Error())
	}

	// Generate API signature
	nonce := "abcdef"
	method := "GET"
	resource := "/api/status"

	signature, err := apiSignature(login.UserID, nonce, method, resource, login.Secret)
	if err != nil {
		t.Fatalf("Failed to generate API signature: %s", err.Error())
	}
	log.Println("signature:", signature)

	// Generate mock HTTP request
	r, err = http.NewRequest(method, "http://localhost:8080"+resource, nil)
	if err != nil {
		t.Fatalf("Failed to generate HTTP request: %s", err.Error())
	}

	headers = map[string][]string{
		"Authorization": {"Basic " + base64.URLEncoding.EncodeToString([]byte(login.Pubkey+":"+nonce+"/"+signature))},
	}
	r.Header = headers

	// Capture HTTP response with recorder
	w = httptest.NewRecorder()

	// Attempt HMAC authentication
	apiAuth = new(HMACAuthenticator)
	clientErr, serverErr = apiAuth.Auth(r)
	if clientErr != nil {
		t.Fatalf("Failed to authenticate: client: %s", clientErr.Error())
	}
	if serverErr != nil {
		t.Fatalf("Failed to authenticate: server: %s", serverErr.Error())
	}

	// Invoke API router
	Router(w, r, user2)

	// Read HTTP response body
	body, err = ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fatalf("Failed to read HTTP response body: %s", err.Error())
	}
	log.Println(string(body))

	// Delete mock user
	if err := user2.Delete(); err != nil {
		t.Fatalf("Failed to delete mock user: %s", err.Error())
	}
}
