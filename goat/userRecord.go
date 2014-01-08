package goat

import (
	"database/sql"
	"log"
)

// UserRecord represents a user on the tracker
type UserRecord struct {
	ID           int
	Username     string
	Passkey      string
	TorrentLimit int `db:"torrent_limit"`
}

// Save UserRecord to storage
func (u UserRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Insert or update a user record
	query := "INSERT INTO users " +
		"(`username`, `passkey`, `torrent_limit`) " +
		"VALUES (?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"`username`=values(`username`), `passkey`=values(`passkey`), `torrent_limit`=values(`torrent_limit`);"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, u.Username, u.Passkey, u.TorrentLimit)
	tx.Commit()

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

	// Fetch announce log into struct
	u = UserRecord{}
	err = db.Get(&u, "SELECT * FROM users WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return UserRecord{}
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

	// Anonymous Uploaded struct
	uploaded := struct {
		Uploaded int64
	}{
		0,
	}

	// Calculate sum of this user's upload via their file/user relationship records
	err = db.Get(&uploaded, "SELECT SUM(uploaded) AS uploaded FROM files_users WHERE user_id=?", u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return uploaded.Uploaded
}

// Downloaded loads this user's total download
func (u UserRecord) Downloaded() int64 {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Downloaded struct
	downloaded := struct {
		Downloaded int64
	}{
		0,
	}

	// Calculate sum of this user's download via their file/user relationship records
	err = db.Get(&downloaded, "SELECT SUM(downloaded) AS downloaded FROM files_users WHERE user_id=?", u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return downloaded.Downloaded
}

// Seeding counts the number of torrents this user is seeding
func (u UserRecord) Seeding() int {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Seeding struct
	seeding := struct {
		Seeding int
	}{
		0,
	}

	// Calculate sum of this user's seeding torrents via their file/user relationship records
	err = db.Get(&seeding, "SELECT COUNT(user_id) AS seeding FROM files_users WHERE user_id = ? AND active = 1 AND completed = 1 AND `left` = 0", u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return seeding.Seeding
}

// Leeching counts the number of torrents this user is leeching
func (u UserRecord) Leeching() int {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Leeching struct
	leeching := struct {
		Leeching int
	}{
		0,
	}

	// Calculate sum of this user's leeching torrents via their file/user relationship records
	err = db.Get(&leeching, "SELECT COUNT(user_id) AS leeching FROM files_users WHERE user_id = ? AND active = 1 AND completed = 0 AND `left` > 0", u.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return leeching.Leeching
}
