package goat

import (
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
	if err := db.SaveWhitelistRecord(w); err != nil {
		log.Println(err.Error())
		return false
	}
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
	if w, err = db.LoadWhitelistRecord(id, col); err != nil {
		log.Println(err.Error())
		return whitelistRecord{}
	}
	return w
}
