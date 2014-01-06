package goat

// WhitelistRecord represents a whitelist entry
type WhitelistRecord struct {
	ID       int
	Client   string
	Approved bool
}

// Save WhitelistRecord to storage
func (w WhitelistRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Store whitelist record
	// NOTE: Not using INSERT IGNORE because it ignores all errors
	// Thanks: http://stackoverflow.com/questions/2366813/on-duplicate-key-ignore
	query := "INSERT INTO whitelist " +
		"(`client`, `approved`) " +
		"VALUES (?, ?) " +
		"ON DUPLICATE KEY UPDATE `client`=`client`;"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, w.Client, w.Approved)
	tx.Commit()

	return true
}

// Load WhitelistRecord from storage
func (w WhitelistRecord) Load(id interface{}, col string) WhitelistRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return w
	}

	// Fetch record into struct
	w = WhitelistRecord{}
	err = db.Get(&w, "SELECT * FROM whitelist WHERE `"+col+"`=?", id)
	if err != nil {
		Static.LogChan <- err.Error()
		return WhitelistRecord{}
	}

	return w
}
