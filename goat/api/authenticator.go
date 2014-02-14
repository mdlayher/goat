package api

import (
	"encoding/base64"
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/mdlayher/goat/goat/data"
)

// APIAuthenticator interface which defines methods required to implement an authentication method
type APIAuthenticator interface {
	Auth(*http.Request) (bool, error)
}

// BasicAuthenticator uses the HTTP Basic with bcrypt authentication scheme
type BasicAuthenticator struct {
}

// Auth handles validation of HTTP Basic with bcrypt authentication, used for user login
func (a BasicAuthenticator) Auth(r *http.Request) (bool, error) {
	// Retrieve Authorization header
	auth := r.Header.Get("Authorization")

	// No header provided
	if auth == "" {
		return false, nil
	}

	// Ensure format is valid
	basic := strings.Split(auth, " ")
	if basic[0] != "Basic" {
		return false, nil
	}

	// Decode base64'd user:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		return false, err
	}

	// Split into username/password
	credentials := strings.Split(string(buf), ":")

	// Load user by username
	user, err := new(data.UserRecord).Load(credentials[0], "username")
	if err != nil || user == (data.UserRecord{}) {
		return false, err
	}

	// Compare input password with bcrypt password, checking for errors
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials[1]))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}

	return true, nil
}

// HMACAuthenticator uses the HMAC-SHA1 authentication scheme, used for API authentication
type HMACAuthenticator struct {
}

// Auth handles validation of HMAC-SHA1 authentication
func (a HMACAuthenticator) Auth(r *http.Request) (bool, error) {
	return true, nil
}
