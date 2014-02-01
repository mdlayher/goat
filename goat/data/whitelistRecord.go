package data

import (
	"log"
)

// WhitelistRecord represents a whitelist entry
type WhitelistRecord struct {
	ID       int
	Client   string
	Approved bool
}

// Save WhitelistRecord to storage
func (w WhitelistRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save WhitelistRecord
	if err := db.SaveWhitelistRecord(w); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load WhitelistRecord from storage
func (w WhitelistRecord) Load(id interface{}, col string) WhitelistRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return w
	}

	// Load WhitelistRecord using specified column
	if w, err = db.LoadWhitelistRecord(id, col); err != nil {
		log.Println(err.Error())
		return WhitelistRecord{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return w
}
