package goat

import (
	"crypto/sha1"
	"fmt"
)

// Struct representing an announce, to be logged to storage
type AnnounceLog struct {
	Id         int
	InfoHash   string `db:"info_hash"`
	Passkey    string
	Ip         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	Client     string
	Time       int64
}

// Generate a SHA1 hash of the form: announce_log_InfoHash
func (log AnnounceLog) Hash() string {
	hash := sha1.New()
	hash.Write([]byte("announce_log" + log.InfoHash))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Save AnnounceLog to storage
func (a AnnounceLog) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Store announce log
	query := "INSERT INTO announce_log " +
		"(`info_hash`, `passkey`, `ip`, `port`, `uploaded`, `downloaded`, `left`, `event`, `client`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.Passkey, a.Ip, a.Port, a.Uploaded, a.Downloaded, a.Left, a.Event, a.Client)
	tx.Commit()

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(id interface{}, col string) AnnounceLog {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return a
	}

	// Fetch announce log into struct
	a = AnnounceLog{}
	db.Get(&a, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)

	return a
}
