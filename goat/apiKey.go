package goat

import (
	"database/sql"
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

	// Insert or update an API key
	query := "INSERT INTO api_keys " +
		"(`user_id`, `key`) " +
		"VALUES (?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"`key`=values(`key`);"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, a.UserID, a.Key)
	tx.Commit()

	return true
}

// Load apiKey from storage
func (a apiKey) Load(id interface{}, col string) apiKey {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return a
	}

	// Fetch apiKey into struct
	a = apiKey{}
	err = db.Get(&a, "SELECT * FROM api_keys WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return apiKey{}
	}

	return a
}
