package data

import (
	"log"
)

// APIKey represents a user's API key
type APIKey struct {
	ID     int
	UserID int `db:"user_id"`
	Key    string
	Salt   string
}

// Delete APIKey from storage
func (a APIKey) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete APIKey
	if err = db.DeleteAPIKey(a.Key, "key"); err != nil {
		log.Println(err.Error())
		return false
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Save APIKey to storage
func (a APIKey) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save APIKey
	if err := db.SaveAPIKey(a); err != nil {
		log.Println(err.Error())
		return false
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load APIKey from storage
func (a APIKey) Load(id interface{}, col string) (key APIKey) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Load APIKey
	if key, err = db.LoadAPIKey(id, col); err != nil {
		log.Println(err.Error())
	}

	// Close database connection
	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}
