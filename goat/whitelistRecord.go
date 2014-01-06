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
	if err := db.SaveWhitelistRecord(w); nil != err {
		Static.LogChan <- err.Error()
		return false
	}
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
	w, err = db.LoadWhitelistRecord(id, col)
	if err != nil {
		Static.LogChan <- err.Error()
		return WhitelistRecord{}
	}
	return w
}
