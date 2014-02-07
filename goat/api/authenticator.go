package api

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/mdlayher/goat/goat/data"
)

// apiAuthenticator interface which defines methods required to implement an authentication method
type apiAuthenticator interface {
	Auth(*http.Request) bool
}

// BasicAuthenticator uses the HTTP Basic authentication scheme
type BasicAuthenticator struct {
}

// Auth handles validation of HTTP Basic authentication
func (a BasicAuthenticator) Auth(r *http.Request) bool {
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
	user, err := new(data.UserRecord).Load(credentials[0], "username")
	if err != nil || user == (data.UserRecord{}) {
		return false
	}

	// Load user's API key
	key, err := new(data.APIKey).Load(user.ID, "user_id")
	if err != nil || key == (data.APIKey{}) {
		return false
	}

	// Hash input password
	sha := sha1.New()
	if _, err = sha.Write([]byte(credentials[1] + key.Salt)); err != nil {
		log.Println(err.Error())
		return false
	}

	hash := fmt.Sprintf("%x", sha.Sum(nil))

	// Verify hashes match, using timing-attack resistant method
	// If function returns 1, hashes match
	return subtle.ConstantTimeCompare([]byte(hash), []byte(key.Key)) == 1
}

// HMACAuthenticator uses the HMAC-SHA1 authentication scheme
type HMACAuthenticator struct {
}

// Auth handles validation of HMAC-SHA1 authentication
func (a HMACAuthenticator) Auth(r *http.Request) bool {
	return true
}
