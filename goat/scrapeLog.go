package goat

import (
	"database/sql"
	"encoding/hex"
	"log"
	"time"
)

// ScrapeLog represents a scrapelog, to be logged to storage
type ScrapeLog struct {
	ID       int
	InfoHash string `db:"info_hash"`
	Passkey  string
	IP       string
	Time     int64
}

// Save ScrapeLog to storage
func (s ScrapeLog) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Store scrape log
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.IP)
	tx.Commit()

	return true
}

// Load ScrapeLog from storage
func (s ScrapeLog) Load(id interface{}, col string) ScrapeLog {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return s
	}

	// Fetch announce log into struct
	s = ScrapeLog{}
	err = db.Get(&s, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return ScrapeLog{}
	}

	return s
}

// FromMap generates a ScrapeLog struct from a query map
func (s ScrapeLog) FromMap(query map[string]string) ScrapeLog {
	s = ScrapeLog{}

	// Required parameters

	// info_hash
	s.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// passkey
	s.Passkey = query["passkey"]

	// ip
	s.IP = query["ip"]

	// Current UNIX timestamp
	s.Time = time.Now().Unix()

	// Return the created scrape
	return s
}
