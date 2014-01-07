package goat

import (
	"log"
)

// FileUserRecord represents a file tracked by tracker
type FileUserRecord struct {
	FileID     int `db:"file_id"`
	UserID     int `db:"user_id"`
	IP         string
	Active     bool
	Completed  bool
	Announced  int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Time       int64
}

// Save FileUserRecord to storage
func (f FileUserRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Insert or update a file/user relationship record
	query := "INSERT INTO files_users " +
		"(`file_id`, `user_id`, `ip`, `active`, `completed`, `announced`, `uploaded`, `downloaded`, `left`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP()) " +
		"ON DUPLICATE KEY UPDATE " +
		"`active`=values(`active`), `completed`=values(`completed`), `announced`=values(`announced`), " +
		"`uploaded`=values(`uploaded`), `downloaded`=values(`downloaded`), `left`=values(`left`), " +
		"`time`=UNIX_TIMESTAMP();"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, f.FileID, f.UserID, f.IP, f.Active, f.Completed, f.Announced, f.Uploaded, f.Downloaded, f.Left)
	tx.Commit()

	return true
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileID int, userID int, ip string) FileUserRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return f
	}

	// Fetch announce log into struct
	f = FileUserRecord{}
	err = db.Get(&f, "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?", fileID, userID, ip)
	if err != nil {
		log.Println(err.Error())
		return FileUserRecord{}
	}

	return f
}
