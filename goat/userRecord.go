package goat

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
		Static.LogChan <- err.Error()
		return false
	}
	if err := db.SaveUserRecord(u); err != nil {
		Static.LogChan <- err.Error()
		return false
	}
	return true
}

// Load UserRecord from storage
func (u UserRecord) Load(id interface{}, col string) UserRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return u
	}
	u, err = db.LoadUserRecord(id, col)
	if err != nil {
		Static.LogChan <- err.Error()
		return UserRecord{}
	}
	return u
}

// Uploaded loads this user's total upload
func (u UserRecord) Uploaded() int64 {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return -1
	}
	uploaded, err := db.GetUserUploaded(u.ID)
	if err != nil {
		Static.LogChan <- err.Error()
		return -1
	}
	return uploaded
}

// Downloaded loads this user's total download
func (u UserRecord) Downloaded() int64 {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return 0
	}
	downloaded, err := db.GetUserDownloaded(u.ID)
	if err != nil {
		Static.LogChan <- err.Error()
		return -1
	}
	return downloaded
}
