package goat

import (
	"log"
)

// apiKey represents a user's API key
type apiKey struct {
	ID     int
	UserID int `db:"user_id"`
	Key    string
}

// Save apiKey to storage
func (a apiKey) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if err := db.SaveApiKey(a); err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// Load apiKey from storage
func (a apiKey) Load(id interface{}, col string) (key apiKey) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if key, err = db.LoadApiKey(id, col); err != nil {
		log.Println(err.Error())
	}
	return
}
