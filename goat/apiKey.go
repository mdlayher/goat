package goat

import (
	"log"
)

// apiKey represents a user's API key
type apiKey struct {
	ID     int
	UserID int `db:"user_id"`
	Key    string
	Salt   string
}

// Delete apiKey from storage
func (a apiKey) Delete() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete apiKey
	if err = db.DeleteAPIKey(a.ID, "id"); err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

// Save apiKey to storage
func (a apiKey) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save apiKey
	if err := db.SaveAPIKey(a); err != nil {
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

	// Load apiKey
	if key, err = db.LoadAPIKey(id, col); err != nil {
		log.Println(err.Error())
	}

	return
}
