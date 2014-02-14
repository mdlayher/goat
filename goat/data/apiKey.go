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
	Key    string
	Expire int64
}

// JSONAPIKey represents output APIKey JSON for API
type JSONAPIKey struct {
	UserID int    `json:"userId"`
	Key    string `json:"key"`
}

// ToJSON converts an APIKey to a JSONAPIKey struct
func (a APIKey) ToJSON() (JSONAPIKey, error) {
	j := JSONAPIKey{}
	j.UserID = a.UserID
	j.Key = a.Key

	return j, nil
}

// Create a new APIKey
func (a *APIKey) Create(userID int) error {
	a.UserID = userID

	// Generate API key using a random SHA1 hash
	sha := sha1.New()
	if _, err := sha.Write([]byte(common.RandString())); err != nil {
		return err
	}
	a.Key = fmt.Sprintf("%x", sha.Sum(nil))

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
	if err = db.DeleteAPIKey(a.Key, "key"); err != nil {
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
