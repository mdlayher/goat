package goat

import (
	"encoding/hex"
	"time"
)

// Struct representing a scrapelog, to be logged to storage
type ScrapeLog struct {
	Id       int
	InfoHash string `db:"info_hash"`
	Passkey  string
	Ip       string
	Time     int64
}

// Save ScrapeLog to storage
func (s ScrapeLog) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Store scrape log
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.Ip)
	tx.Commit()

	return true
}

// Load ScrapeLog from storage
func (s ScrapeLog) Load(id interface{}, col string) ScrapeLog {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return s
	}

	// Fetch announce log into struct
	s = ScrapeLog{}
	db.Get(&s, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)

	return s
}

// Generate a ScrapeLog struct from a query map
func (s ScrapeLog) FromMap(query map[string]string) ScrapeLog {
	s = ScrapeLog{}

	// Required parameters

	// info_hash
	s.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// passkey
	s.Passkey = query["passkey"]

	// ip
	s.Ip = query["ip"]

	// Current UNIX timestamp
	s.Time = time.Now().Unix()

	// Return the created scrape
	return s
}
