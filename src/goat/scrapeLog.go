package goat

// Struct representing a scrapelog, to be logged to storage
type ScrapeLog struct {
	Id       int
	InfoHash string `db:"info_hash"`
	Passkey  string
	Ip       string
	Time     int64
}

// Save ScrapeLog to storage
func (s ScrapeLog) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Store scrape log
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.Ip)
	tx.Commit()

	return true
}

// Load ScrapeLog from storage
func (s ScrapeLog) Load(id interface{}, col string) ScrapeLog {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return s
	}

	// Fetch announce log into struct
	s = ScrapeLog{}
	db.Get(&s, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)

	return s
}
