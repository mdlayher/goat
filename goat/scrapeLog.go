package goat

import (
	"encoding/hex"
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
		Static.LogChan <- err.Error()
		return false
	}
	if err := db.SaveScrapeLog(s); nil != err {
		Static.LogChan <- err.Error()
		return false
	}
	return true
}

// Load ScrapeLog from storage
func (s ScrapeLog) Load(id interface{}, col string) ScrapeLog {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return s
	}
	s, err = db.LoadScrapeLog(id, col)
	if err != nil {
		Static.LogChan <- err.Error()
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
