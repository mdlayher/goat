package goat

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
		Static.LogChan <- err.Error()
		return false
	}

	if err := db.SaveFileUserRecord(f); nil != err {
		Static.LogChan <- err.Error()
		return false
	}

	return true
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileID int, userID int, ip string) FileUserRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return f
	}

	f, err = db.LoadFileUserRecord(fileID, userID, ip)
	if err != nil {
		Static.LogChan <- err.Error()
	}
	return f
}
