package data

// WhitelistRecord represents a whitelist entry
type WhitelistRecord struct {
	ID       int
	Client   string
	Approved bool
}

// Delete WhitelistRecord from storage
func (w WhitelistRecord) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete WhitelistRecord
	if err = db.DeleteWhitelistRecord(w.Client, "client"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load WhitelistRecord from storage
func (w WhitelistRecord) Load(id interface{}, col string) (WhitelistRecord, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return WhitelistRecord{}, err
	}

	// Load WhitelistRecord using specified column
	if w, err = db.LoadWhitelistRecord(id, col); err != nil {
		return WhitelistRecord{}, err
	}

	if err := db.Close(); err != nil {
		return WhitelistRecord{}, err
	}

	return w, nil
}

// Save WhitelistRecord to storage
func (w WhitelistRecord) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save WhitelistRecord
	if err := db.SaveWhitelistRecord(w); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
