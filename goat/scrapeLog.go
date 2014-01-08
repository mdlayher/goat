package goat

import (
	"database/sql"
	"encoding/hex"
	"log"
	"net/url"
	"time"
)

// scrapeLog represents a scrapelog, to be logged to storage
type scrapeLog struct {
	ID       int
	InfoHash string `db:"info_hash"`
	Passkey  string
	IP       string
	Time     int64
}

// Save scrapeLog to storage
func (s scrapeLog) Save() bool {
	// Open database connection
	db, err := dbConnect()
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

// Load scrapeLog from storage
func (s scrapeLog) Load(id interface{}, col string) scrapeLog {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return s
	}

	// Fetch announce log into struct
	s = scrapeLog{}
	err = db.Get(&s, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return scrapeLog{}
	}

	return s
}

// FromValues generates a scrapeLog struct from a url.Values map
func (s scrapeLog) FromValues(query url.Values) scrapeLog {
	s = scrapeLog{}

	// Required parameters

	// info_hash
	s.InfoHash = hex.EncodeToString([]byte(query.Get("info_hash")))

	// passkey
	s.Passkey = query.Get("passkey")

	// ip
	s.IP = query.Get("ip")

	// Current UNIX timestamp
	s.Time = time.Now().Unix()

	// Return the created scrape
	return s
}
