package goat

import (
	"database/sql"
	"log"
)

// whitelistRecord represents a whitelist entry
type whitelistRecord struct {
	ID       int
	Client   string
	Approved bool
}

// Save whitelistRecord to storage
func (w whitelistRecord) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Store whitelist record
	// NOTE: Not using INSERT IGNORE because it ignores all errors
	// Thanks: http://stackoverflow.com/questions/2366813/on-duplicate-key-ignore
	query := "INSERT INTO whitelist " +
		"(`client`, `approved`) " +
		"VALUES (?, ?) " +
		"ON DUPLICATE KEY UPDATE `client`=`client`;"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, w.Client, w.Approved)
	tx.Commit()

	return true
}

// Load whitelistRecord from storage
func (w whitelistRecord) Load(id interface{}, col string) whitelistRecord {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return w
	}

	// Fetch record into struct
	w = whitelistRecord{}
	err = db.Get(&w, "SELECT * FROM whitelist WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return whitelistRecord{}
	}

	return w
}
