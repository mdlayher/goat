package goat

import (
	"encoding/hex"
	"time"
)

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id        string
	Busy      bool
	MapStor   map[string]interface{}
	MapLookup map[string][]*interface{}
}

// addes new maps recursivly depending on how deep the structure needs to be
func addMap(m map[string]interface{}, size int) {

	for i := 0; i < 16; i++ {
		c := hex.EncodeToString([]byte(string(i)))[1:]
		m[c] = make(map[string]interface{})
		if size > 0 {
			addMap(m[c].(map[string]interface{}), size-1)
		}
	}
}
func (db MapDb) init() {
	if db.MapStor == nil {
		size := 2
		db.MapStor = make(map[string]interface{})
		addMap(db.MapStor, size)
	}
	if db.MapLookup == nil {
		db.MapLookup = make(map[string][]*interface{})
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
