package goat

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// SqlDb is a Sql based database
type SqlDb struct {
}

//MapDb write
func (db SqlDb) Write(req Request) {
	switch req.Data.(type) {
	case AnnounceLog:
	case FileRecord:
	case FileUserRecord:
	default:
	}
}
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
