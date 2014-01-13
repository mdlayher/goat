package goat

import (
	"log"
)

// userRecord represents a user on the tracker
type userRecord struct {
	ID           int
	Username     string
	Passkey      string
	TorrentLimit int `db:"torrent_limit"`
}

// Save userRecord to storage
func (u userRecord) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}
	if err := db.SaveUserRecord(u); err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

// Load userRecord from storage
func (u userRecord) Load(id interface{}, col string) userRecord {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return u
	}
	u, err = db.LoadUserRecord(id, col)
	if err != nil {
		log.Println(err.Error())
		return userRecord{}
	}
	return u
}

// Uploaded loads this user's total upload
func (u userRecord) Uploaded() int64 {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	uploaded, err := db.GetUserUploaded(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	return uploaded
}

// Downloaded loads this user's total download
func (u userRecord) Downloaded() int64 {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	downloaded, err := db.GetUserDownloaded(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	return downloaded
}

//

// Seeding counts the number of torrents this user is seeding
func (u userRecord) Seeding() int {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	seeding, err := db.GetUserSeeding(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	return seeding
}

// Leeching counts the number of torrents this user is leeching
func (u userRecord) Leeching() int {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	leeching, err := db.GetUserSeeding(u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}
	return leeching
}
