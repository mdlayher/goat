package data

// APIKey represents a user's API key
type APIKey struct {
	ID     int
	UserID int `db:"user_id"`
	Key    string
	Salt   string
}

// Delete APIKey from storage
func (a APIKey) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete APIKey
	if err = db.DeleteAPIKey(a.Key, "key"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load APIKey from storage
func (a APIKey) Load(id interface{}, col string) (APIKey, error) {
	a = APIKey{}

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return a, err
	}

	// Load APIKey
	a, err = db.LoadAPIKey(id, col)
	if err != nil {
		return a, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return a, err
	}

	return a, nil
}

// Save APIKey to storage
func (a APIKey) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save APIKey
	if err := db.SaveAPIKey(a); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}
