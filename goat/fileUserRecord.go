package goat

import (
	"database/sql"
	"log"
)

// FileUserRecord represents a file tracked by tracker
type FileUserRecord struct {
	FileID     int    `db:"file_id" json:"fileId"`
	UserID     int    `db:"user_id" json:"userId"`
	IP         string `json:"ip"`
	Active     bool   `json:"active"`
	Completed  bool   `json:"completed"`
	Announced  int    `json:"announced"`
	Uploaded   int64  `json:"uploaded"`
	Downloaded int64  `json:"downloaded"`
	Left       int64  `json:"left"`
	Time       int64  `json:"time"`
}

// FileUserRecordRepository is used to contain methods to load multiple FileRecord structs
type FileUserRecordRepository struct {
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
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return FileUserRecord{}
	}

	return f
}

// Select loads selected FileUserRecord structs from storage
func (f FileUserRecordRepository) Select(id interface{}, col string) []FileUserRecord {
	fileUsers := make([]FileUserRecord, 0)

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return fileUsers
	}

	// Load all files
	rows, err := db.Queryx("SELECT * FROM files_users WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return fileUsers
	}

	// Iterate all rows and build array
	fileUser := FileUserRecord{}
	for rows.Next() {
		// Scan row results
		rows.StructScan(&fileUser)
		fileUsers = append(fileUsers[:], fileUser)
	}

	return fileUsers
}
