package data

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/mdlayher/goat/goat/common"
)

// APIKey represents a user's API key
type APIKey struct {
	ID     int
	UserID int `db:"user_id"`
	Pubkey string
	Secret string
	Expire int64
}

// JSONAPIKey represents output APIKey JSON for API
type JSONAPIKey struct {
	UserID int    `json:"userId"`
	Pubkey string `json:"pubkey"`
	Secret string `json:"secret"`
	Expire int64  `json:"expire"`
}

// ToJSON converts an APIKey to a JSONAPIKey struct
func (a APIKey) ToJSON() (JSONAPIKey, error) {
	j := JSONAPIKey{}
	j.UserID = a.UserID
	j.Pubkey = a.Pubkey
	j.Secret = a.Secret
	j.Expire = a.Expire

	return j, nil
}

// Create a new APIKey
func (a *APIKey) Create(userID int) error {
	a.UserID = userID

	// Generate API pubkey using a random SHA1 hash
	sha := sha1.New()
	if _, err := sha.Write([]byte(common.RandString())); err != nil {
		return err
	}
	a.Pubkey = fmt.Sprintf("%x", sha.Sum(nil))

	// Generate API secret using a random SHA1 hash
	sha2 := sha1.New()
	if _, err := sha2.Write([]byte(common.RandString())); err != nil {
		return err
	}
	a.Secret = fmt.Sprintf("%x", sha2.Sum(nil))

	// Set key to expire one week from now
	a.Expire = time.Now().Add(7 * 24 * time.Hour).Unix()

	return nil
}

// Delete APIKey from storage
func (a APIKey) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete APIKey
	if err = db.DeleteAPIKey(a.Pubkey, "pubkey"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load APIKey from storage
func (a APIKey) Load(id interface{}, col string) (APIKey, error) {
	a = APIKey{}

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return a, err
	}

	// Load APIKey
	a, err = db.LoadAPIKey(id, col)
	if err != nil {
		return a, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return a, err
	}

	return a, nil
}

// Save APIKey to storage
func (a APIKey) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save APIKey
	if err := db.SaveAPIKey(a); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
