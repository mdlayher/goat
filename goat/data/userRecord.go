package data

import (
	"crypto/sha1"
	"fmt"
	"log"

	"github.com/mdlayher/goat/goat/common"
)

// UserRecord represents a user on the tracker
type UserRecord struct {
	ID           int
	Username     string
	Passkey      string
	TorrentLimit int `db:"torrent_limit"`
}

// Create a UserRecord, using defined parameters
func (u UserRecord) Create(username string, torrentLimit int) UserRecord {
	// Set username and torrent limit
	u.Username = username
	u.TorrentLimit = torrentLimit

	// Randomly generate a new passkey
	sha := sha1.New()
	if _, err := sha.Write([]byte(common.RandString())); err != nil {
		log.Println(err.Error())
		return UserRecord{}
	}

	u.Passkey = fmt.Sprintf("%x", sha.Sum(nil))
	return u
}

// Delete UserRecord from storage
func (u UserRecord) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete UserRecord
	if err = db.DeleteUserRecord(u.ID, "id"); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Save UserRecord to storage
func (u UserRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save UserRecord
	if err := db.SaveUserRecord(u); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load UserRecord from storage
func (u UserRecord) Load(id interface{}, col string) UserRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return u
	}

	// Load UserRecord by specified column
	u, err = db.LoadUserRecord(id, col)
	if err != nil {
		log.Println(err.Error())
		return UserRecord{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return u
}

// Uploaded loads this user's total upload
func (u UserRecord) Uploaded() int64 {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	// Retrieve total bytes user has uploaded
	uploaded, err := db.GetUserUploaded(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return uploaded
}

// Downloaded loads this user's total download
func (u UserRecord) Downloaded() int64 {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Retrieve total bytes user has downloaded
	downloaded, err := db.GetUserDownloaded(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return downloaded
}

// Seeding counts the number of torrents this user is seeding
func (u UserRecord) Seeding() int {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Retrieve total number of torrents user is actively seeding
	seeding, err := db.GetUserSeeding(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return seeding
}

// Leeching counts the number of torrents this user is leeching
func (u UserRecord) Leeching() int {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Retrieve total number of torrents user is actively leeching
	leeching, err := db.GetUserSeeding(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return leeching
}
