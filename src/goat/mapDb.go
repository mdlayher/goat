package goat

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id        string
	Busy      bool
	MapStor   map[string]map[string]interface{}
	MapLookup map[string]*interface{}
}

func (db MapDb) init() {
	if db.MapStor == nil {
		db.MapStor = make(map[string]map[string]interface{})
	}
	if db.MapLookup == nil {
		db.MapLookup = make(map[string]*interface{})
	}
}

//MapDb write
func (db MapDb) Write(req Request) {
	switch req.Data.(type) {
	case AnnounceLog:
	case FileRecord:
	case FileUserRecord:
	default:
	}
}
func (db MapDb) Read(req Request) {
	switch req.Data.(type) {
	case AnnounceLog:
	case FileRecord:
	case FileUserRecord:
	default:
	}
}

// Shutdown MapDb
func (db MapDb) Shutdown() {
	// Wait until map is no longer busy
	for db.Busy {
		time.Sleep(500 * time.Millisecond)
	}

	Static.LogChan <- "stopping MapDb"
}
