package api

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	// Fetch credentials from HTTP Basic auth
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
	// Check for Authorization header
	auth := r.Header.Get("Authorization")
	if auth == "" {
		// Check for X-Goat-Authorization header override
		auth = r.Header.Get("X-Goat-Authorization")
		if auth == "" {
			return false, nil
		}
	}

	// Fetch credentials from HTTP Basic auth
	pubkey, apiSignature, err := basicCredentials(auth)
	if err != nil {
		return false, err
	}

	// Load API key by pubkey
	key, err := new(data.APIKey).Load(pubkey, "pubkey")
	if err != nil || key == (data.APIKey{}) {
		return false, err
	}

	// Check if key is expired, delete it if it is
	if key.Expire <= time.Now().Unix() {
		go func(key data.APIKey) {
			if err := key.Delete(); err != nil {
				log.Println(err.Error())
			}
		}(key)

		return false, nil
	}

	// Generate API signature string
	signString := fmt.Sprintf("%d-%s-%s", key.UserID, r.Method, r.URL.Path)

	// Calculate HMAC-SHA1 signature from string, using API secret
	mac := hmac.New(sha1.New, []byte(key.Secret))
	if _, err := mac.Write([]byte(signString)); err != nil {
		return false, err
	}
	expected := fmt.Sprintf("%x", mac.Sum(nil))

	// Verify that HMAC signature is correct
	if !hmac.Equal([]byte(apiSignature), []byte(expected)) {
		return false, nil
	}

	// Update API key expiration time
	key.Expire = time.Now().Add(7 * 24 * time.Hour).Unix()
	go func(key data.APIKey) {
		if err := key.Save(); err != nil {
			log.Println(err.Error())
		}
	}(key)

	// Load user by user ID
	user, err := new(data.UserRecord).Load(key.UserID, "id")
	if err != nil || user == (data.UserRecord{}) {
		return false, err
	}

	// Store user for session
	a.session = user
	return true, nil
}

// Session attempts to return the user whose session was authenticated via this authenticator
func (a HMACAuthenticator) Session() (data.UserRecord, error) {
	if a.session == (data.UserRecord{}) {
		return data.UserRecord{}, errors.New("session: no session found")
	}

	return a.session, nil
}
