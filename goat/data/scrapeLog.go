package data

import (
	"database/sql"
	"encoding/hex"
	"errors"
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
func (s ScrapeLog) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete ScrapeLog
	if err = db.DeleteScrapeLog(s.ID, "id"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Save ScrapeLog to storage
func (s ScrapeLog) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save ScrapeLog
	if err := db.SaveScrapeLog(s); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load ScrapeLog from storage
func (s ScrapeLog) Load(id interface{}, col string) (ScrapeLog, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return ScrapeLog{}, err
	}

	// Load ScrapeLog
	s, err = db.LoadScrapeLog(id, col)
	if err != nil && err != sql.ErrNoRows {
		return ScrapeLog{}, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return ScrapeLog{}, err
	}

	return s, nil
}

// FromValues generates a ScrapeLog struct from a url.Values map
func (s *ScrapeLog) FromValues(query url.Values) error {
	// Required parameters

	// info_hash (20 characters, 40 characters after hex encode)
	s.InfoHash = hex.EncodeToString([]byte(query.Get("info_hash")))
	if len(s.InfoHash) != 40 {
		return errors.New("info_hash must be exactly 20 characters")
	}

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

	return nil
}
