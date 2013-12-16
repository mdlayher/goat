package goat

import (
	"crypto/sha1"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Generates a unique hash for a given struct with its ID, to be used as ID for storage
type Hashable interface {
	Hash() string
}

// Defines methods to save, load, and delete data for an object
type Persistent interface {
	Save() bool
	Load(interface{}) interface{}
	Delete() bool
}

type ErrorRes struct {
	ErrLocation string
	Time        int64
	Error       string
}

// Struct representing an announce, to be logged to storage
type AnnounceLog struct {
	Id         int
	InfoHash   string `db:"info_hash"`
	PeerId     string `db:"peer_id"`
	Ip         string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Event      string
	Time       int64
}

// Generate a SHA1 hash of the form: announce_log_InfoHash
func (log AnnounceLog) Hash() string {
	hash := sha1.New()
	hash.Write([]byte("announce_log"+log.InfoHash))
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

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl("INSERT INTO announce_log (`info_hash`, `peer_id`, `ip`, `port`, `uploaded`, `downloaded`, `left`, `event`, `time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());",
		a.InfoHash, a.PeerId, a.Ip, a.Port, a.Uploaded, a.Downloaded, a.Left, a.Event)
	tx.Commit()

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(id interface{}) AnnounceLog {
	// Open database connection
	db, _ := DbConnect()

	// Create database transaction
	res := AnnounceLog{}
	db.Get(&res, "SELECT * FROM announce_log WHERE id=?", id)

	return res
}

// Struct representing a file tracked by tracker
type FileRecord struct {
	InfoHash   string
	Leechers   int
	Seeders    int
	Completed  int
	Time       int64
	CreateTime int64
}

// Struct representing a scrapelog, to be logged to storage
type ScrapeLog struct {
	InfoHash string
	UserId   string
	Ip       string
	Time     int64
}
