package data

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

// Delete FileUserRecord from storage
func (f FileUserRecord) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete FileUserRecord
	if err = db.DeleteFileUserRecord(f.FileID, f.UserID, f.IP); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileID int, userID int, ip string) (FileUserRecord, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return FileUserRecord{}, err
	}

	// Load FileUserRecord using file ID, user ID, IP triple
	f, err = db.LoadFileUserRecord(fileID, userID, ip)
	if err != nil {
		return FileUserRecord{}, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return FileUserRecord{}, err
	}

	return f, nil
}

// Save FileUserRecord to storage
func (f FileUserRecord) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save FileUserRecord
	if err := db.SaveFileUserRecord(f); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Select loads selected FileUserRecord structs from storage
func (f FileUserRecordRepository) Select(id interface{}, col string) ([]FileUserRecord, error) {
	fileUsers := make([]FileUserRecord, 0)

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return fileUsers, err
	}

	// Load FileUserRecords matching specified conditions
	fileUsers, err = db.LoadFileUserRepository(id, col)
	if err != nil {
		return fileUsers, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return fileUsers, err
	}

	return fileUsers, nil
}
