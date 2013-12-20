package goat

// Struct representing a file tracked by tracker
type FileUserRecord struct {
	FileId     int `db:"file_id"`
	UserId     int `db:"user_id"`
	Ip         string
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
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
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
	tx.Execl(query, f.FileId, f.UserId, f.Ip, f.Active, f.Completed, f.Announced, f.Uploaded, f.Downloaded, f.Left)
	tx.Commit()

	return true
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileId int, userId int, ip string) FileUserRecord {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return f
	}

	// Fetch announce log into struct
	f = FileUserRecord{}
	db.Get(&f, "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?", fileId, userId, ip)
	return f
}
