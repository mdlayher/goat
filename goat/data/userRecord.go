package data

import (
	"crypto/sha1"
	"fmt"

	"github.com/mdlayher/goat/goat/common"
)

// UserRecord represents a user on the tracker
type UserRecord struct {
	ID           int
	Username     string
	Passkey      string
	TorrentLimit int `db:"torrent_limit"`
}

// UserRecordRepository is used to contain methods to load multiple UserRecord structs
type UserRecordRepository struct {
}

// JSONUserRecord represents output UserRecord JSON for API
type JSONUserRecord struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	TorrentLimit int    `json:"torrentLimit"`
}

// ToJSON converts a UserRecord to a JSONUserRecord struct
func (u UserRecord) ToJSON() (JSONUserRecord, error) {
	// Convert standard fields to JSON equivalent struct
	j := JSONUserRecord{}
	j.ID = u.ID
	j.Username = u.Username
	j.TorrentLimit = u.TorrentLimit

	return j, nil
}

// Create a UserRecord, using defined parameters
func (u *UserRecord) Create(username string, torrentLimit int) error {
	// Set username and torrent limit
	u.Username = username
	u.TorrentLimit = torrentLimit

	// Randomly generate a new passkey
	sha := sha1.New()
	if _, err := sha.Write([]byte(common.RandString())); err != nil {
		return err
	}

	u.Passkey = fmt.Sprintf("%x", sha.Sum(nil))
	return nil
}

// Delete UserRecord from storage
func (u UserRecord) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete UserRecord
	if err = db.DeleteUserRecord(u.ID, "id"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Save UserRecord to storage
func (u UserRecord) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save UserRecord
	if err := db.SaveUserRecord(u); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load UserRecord from storage
func (u UserRecord) Load(id interface{}, col string) (UserRecord, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return UserRecord{}, err
	}

	// Load UserRecord by specified column
	u, err = db.LoadUserRecord(id, col)
	if err != nil {
		return UserRecord{}, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return UserRecord{}, err
	}

	return u, nil
}

// Uploaded loads this user's total upload
func (u UserRecord) Uploaded() (int64, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve total bytes user has uploaded
	uploaded, err := db.GetUserUploaded(u.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return uploaded, nil
}

// Downloaded loads this user's total download
func (u UserRecord) Downloaded() (int64, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve total bytes user has downloaded
	downloaded, err := db.GetUserDownloaded(u.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return downloaded, nil
}

// Seeding counts the number of torrents this user is seeding
func (u UserRecord) Seeding() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve total number of torrents user is actively seeding
	seeding, err := db.GetUserSeeding(u.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return seeding, nil
}

// Leeching counts the number of torrents this user is leeching
func (u UserRecord) Leeching() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve total number of torrents user is actively leeching
	leeching, err := db.GetUserSeeding(u.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return leeching, nil
}

// All loads all UserRecord structs from storage
func (u UserRecordRepository) All() ([]UserRecord, error) {
	users := make([]UserRecord, 0)

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return users, err
	}

	// Retrieve all users
	users, err = db.GetAllUserRecords()
	if err != nil {
		return users, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return users, err
	}

	return users, nil
}
