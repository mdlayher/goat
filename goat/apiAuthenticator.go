package goat

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// apiAuthenticator interface which defines methods required to implement an authentication method
type apiAuthenticator interface {
	Auth(*http.Request) bool
}

// basicAPIAuthenticator uses the HTTP Basic authentication scheme
type basicAPIAuthenticator struct {
}

// Auth handles validation of HTTP Basic authentication
func (a basicAPIAuthenticator) Auth(r *http.Request) bool {
	// Retrieve Authorization header
	auth := r.Header.Get("Authorization")

	// No header provided
	if auth == "" {
		return false
	}

	// Ensure format is valid
	basic := strings.Split(auth, " ")
	if basic[0] != "Basic" {
		return false
	}

	// Decode base64'd user:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Split into username/password
	credentials := strings.Split(string(buf), ":")

	// Load user by username, verify user exists
	user := new(userRecord).Load(credentials[0], "username")
	if user == (userRecord{}) {
		return false
	}

	// Hash input password
	// TODO: use something better than SHA1, but still fast
	sha := sha1.New()
	sha.Write([]byte(credentials[1]))
	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Attempt to load user's API key from hash
	key := new(apiKey).Load(hash, "key")
	if key == (apiKey{}) {
		return false
	}

	// Authentication succeeded
	return true
}

// hmacAPIAuthenticator uses the HMAC-SHA1 authentication scheme
type hmacAPIAuthenticator struct {
}

// Auth handles validation of HMAC-SHA1 authentication
func (a hmacAPIAuthenticator) Auth(r *http.Request) bool {
	return true
}
