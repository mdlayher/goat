package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strings"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/mdlayher/goat/goat/data"
)

// APIAuthenticator interface which defines methods required to implement an authentication method
type APIAuthenticator interface {
	Auth(*http.Request) (bool, error)
	Session() (data.UserRecord, error)
}

// basicCredentials returns HTTP Basic authentication credentials from a header
func basicCredentials(header string) (string, string, error) {
	// No header provided
	if header == "" {
		return "", "", errors.New("basicCredentials: empty header provided")
	}

	// Ensure format is valid
	basic := strings.Split(header, " ")
	if basic[0] != "Basic" {
		return "", "", errors.New("basicCredentials: invalid header format")
	}

	// Decode base64'd user:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		return "", "", errors.New("basicCredentials: invalid header format")
	}

	// Split into username/password
	credentials := strings.Split(string(buf), ":")
	return credentials[0], credentials[1], nil
}

// BasicAuthenticator uses the HTTP Basic with bcrypt authentication scheme
type BasicAuthenticator struct {
	session data.UserRecord
}

// Auth handles validation of HTTP Basic with bcrypt authentication, used for user login
func (a *BasicAuthenticator) Auth(r *http.Request) (bool, error) {
	username, password, err := basicCredentials(r.Header.Get("Authorization"))
	if err != nil {
		return false, err
	}

	// Load user by username
	user, err := new(data.UserRecord).Load(username, "username")
	if err != nil || user == (data.UserRecord{}) {
		return false, err
	}

	// Compare input password with bcrypt password, checking for errors
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return false, err
	}

	// Store user for session
	a.session = user
	return true, nil
}

// Session attempts to return the user whose session was authenticated via this authenticator
func (a BasicAuthenticator) Session() (data.UserRecord, error) {
	if a.session == (data.UserRecord{}) {
		return data.UserRecord{}, errors.New("session: no session found")
	}

	return a.session, nil
}

// HMACAuthenticator uses the HMAC-SHA1 authentication scheme, used for API authentication
type HMACAuthenticator struct {
	session data.UserRecord
}

// Auth handles validation of HMAC-SHA1 authentication
func (a *HMACAuthenticator) Auth(r *http.Request) (bool, error) {
	return true, nil
}

// Session attempts to return the user whose session was authenticated via this authenticator
func (a HMACAuthenticator) Session() (data.UserRecord, error) {
	if a.session == (data.UserRecord{}) {
		return data.UserRecord{}, errors.New("session: no session found")
	}

	return a.session, nil
}
