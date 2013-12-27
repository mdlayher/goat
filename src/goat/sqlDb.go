package goat

// SqlDb is a Sql based database
type SqlDb struct {
}

// SqlDb write
func (db SqlDb) Write(req Request) {
	switch req.Data.(type) {
	case AnnounceLog:
	case FileRecord:
	case FileUserRecord:
	default:
	}
}

// SqlDb Read
func (db SqlDb) Read(req Request) {
	switch req.Data.(type) {
	case AnnounceLog:
	case FileRecord:
	case FileUserRecord:
	default:
	}
}

// Shutdown SqlDb
func (db SqlDb) Shutdown() {
}
