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

// usedNonces is a map containing the used nonce values recently
var usedNonces = map[string]int{}

// APIAuthenticator interface which defines methods required to implement an authentication method
type APIAuthenticator interface {
	Auth(*http.Request) (error, error)
	Session() (data.UserRecord, error)
}

// basicCredentials returns HTTP Basic authentication credentials from a header
func basicCredentials(header string) (string, string, error) {
	// No header provided
	if header == "" {
		return "", "", errors.New("empty HTTP Basic header")
	}

	// Ensure format is valid
	basic := strings.Split(header, " ")
	if basic[0] != "Basic" {
		return "", "", errors.New("invalid HTTP Basic header")
	}

	// Decode base64'd user:password pair
	buf, err := base64.URLEncoding.DecodeString(basic[1])
	if err != nil {
		return "", "", errors.New("invalid HTTP Basic header")
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
func (a *BasicAuthenticator) Auth(r *http.Request) (error, error) {
	// Fetch credentials from HTTP Basic auth
	username, password, err := basicCredentials(r.Header.Get("Authorization"))
	if err != nil {
		return err, nil
	}

	// Load user by username
	user, err := new(data.UserRecord).Load(username, "username")
	if err != nil || user == (data.UserRecord{}) {
		return errors.New("no such user"), err
	}

	// Compare input password with bcrypt password, checking for errors
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
		return errors.New("invalid password"), err
	}

	// Store user for session
	a.session = user
	return nil, nil
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
func (a *HMACAuthenticator) Auth(r *http.Request) (error, error) {
	// Check for Authorization header
	auth := r.Header.Get("Authorization")
	if auth == "" {
		// Check for X-Goat-Authorization header override
		auth = r.Header.Get("X-Goat-Authorization")
	}

	// Fetch credentials from HTTP Basic auth
	pubkey, credentials, err := basicCredentials(auth)
	if err != nil {
		return err, nil
	}

	// Split credentials into nonce and API signature
	pair := strings.Split(credentials, "/")
	if len(pair) < 2 {
		return errors.New("no nonce value"), nil
	}

	nonce := pair[0]
	apiSignature := pair[1]

	// Check if nonce previously used, to prevent replay attacks
	if _, ok := usedNonces[nonce]; ok {
		return errors.New("repeated API request"), nil
	}

	// Set nonce into map, set it to expire in 2 minutes to conserve memory
	usedNonces[nonce] = 1
	go func(nonce string) {
		<-time.After(2 * time.Minute)
		delete(usedNonces, nonce)
	}(nonce)

	// Load API key by pubkey
	key, err := new(data.APIKey).Load(pubkey, "pubkey")
	if err != nil || key == (data.APIKey{}) {
		return errors.New("no such public key"), err
	}

	// Check if key is expired, delete it if it is
	if key.Expire <= time.Now().Unix() {
		go func(key data.APIKey) {
			if err := key.Delete(); err != nil {
				log.Println(err.Error())
			}
		}(key)

		return errors.New("expired API key"), nil
	}

	// Generate API signature string
	signString := fmt.Sprintf("%d-%s-%s-%s", key.UserID, nonce, r.Method, r.URL.Path)

	// Calculate HMAC-SHA1 signature from string, using API secret
	mac := hmac.New(sha1.New, []byte(key.Secret))
	if _, err := mac.Write([]byte(signString)); err != nil {
		return nil, err
	}
	expected := fmt.Sprintf("%x", mac.Sum(nil))

	// Verify that HMAC signature is correct
	if !hmac.Equal([]byte(apiSignature), []byte(expected)) {
		return errors.New("invalid API signature"), nil
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
		return errors.New("no such user"), err
	}

	// Store user for session
	a.session = user
	return nil, nil
}

// Session attempts to return the user whose session was authenticated via this authenticator
func (a HMACAuthenticator) Session() (data.UserRecord, error) {
	if a.session == (data.UserRecord{}) {
		return data.UserRecord{}, errors.New("session: no session found")
	}

	return a.session, nil
}
