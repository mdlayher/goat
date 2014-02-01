package data

import (
	"database/sql"
	"encoding/hex"
	"log"
	"net/url"
	"time"
)

// ScrapeLog represents a scrapelog, to be logged to storage
type ScrapeLog struct {
	ID       int
	InfoHash string `db:"info_hash"`
	Passkey  string
	IP       string
	Time     int64
	UDP      bool
}

// Delete ScrapeLog from storage
func (s ScrapeLog) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete ScrapeLog
	if err = db.DeleteScrapeLog(s.ID, "id"); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Save ScrapeLog to storage
func (s ScrapeLog) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save ScrapeLog
	if err := db.SaveScrapeLog(s); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

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

	// Load ScrapeLog
	s, err = db.LoadScrapeLog(id, col)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return ScrapeLog{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return s
}

// FromValues generates a ScrapeLog struct from a url.Values map
func (s ScrapeLog) FromValues(query url.Values) ScrapeLog {
	s = ScrapeLog{}

	// Required parameters

	// info_hash
	s.InfoHash = hex.EncodeToString([]byte(query.Get("info_hash")))

	// passkey
	s.Passkey = query.Get("passkey")

	// ip
	s.IP = query.Get("ip")

	// Current UNIX timestamp
	s.Time = time.Now().Unix()

	// udp
	if query.Get("udp") == "1" {
		s.UDP = true
	} else {
		s.UDP = false
	}

	// Return the created scrape
	return s
}
