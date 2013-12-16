package goat

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"time"
)

// Connect to MySQL database
func DbConnect() (*sqlx.DB, error) {
	return sqlx.Connect("mysql", fmt.Sprintf("%s:%s@/%s", "goat", "goat", "goat"))
}

// Read a struct with specified ID from storage
func DbRead(id string) (Response, bool) {
	// Generate request
	var req Request
	req.Id = id

	// Generate channel for response
	queryResChan := make(chan Response)
	req.ResponseChan = queryResChan

	// Perform request
	Static.RequestChan <- req

	// Read and return response
	res := <-queryResChan
	return res, true
}

// Write a struct with specified ID to storage
func DbWrite(id string, data interface{}) (Response, bool) {
	// Generate request
	var req Request
	req.Id = id
	req.Data = data

	// Generate channel for response
	queryResChan := make(chan Response)
	req.ResponseChan = queryResChan

	// Perform request
	Static.RequestChan <- req

	// Read and return response
	res := <-queryResChan
	return res, true
}

func DbManager(dbDoneChan chan bool) {
	// Storage handler instances
	mapDb := new(MapDb)
	sqlDb := new(SqlDb)

	// channels
	sqlRequestChan := make(chan Request)
	mapRequestChan := make(chan Request, 100)

	// Shutdown function
	go func(dbDoneChan chan bool, mapDb *MapDb, sqlDb *SqlDb) {
		// Wait for shutdown
		<-Static.ShutdownChan

		if Static.Config.Map {
			mapDb.Shutdown()
		}
		if Static.Config.Sql {
			sqlDb.Shutdown()
		}

		dbDoneChan <- true
	}(dbDoneChan, mapDb, sqlDb)

	// launch databases
	if Static.Config.Map {
		go mapDb.HandleDb(mapRequestChan)
		Static.LogChan <- "MapDb instance launched"
	}
	if Static.Config.Sql {
		go sqlDb.HandleDb(sqlRequestChan)
		Static.LogChan <- "SqlDb instance launched"
	}

	if Static.Config.Map && Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				if hold.Data == nil {
					mapRequestChan <- hold
				} else {
					mapRequestChan <- hold
					sqlRequestChan <- hold
				}
			case hold := <-Static.PersistentChan:
				sqlRequestChan <- hold
			}
		}
	} else if Static.Config.Map {
		for {
			select {
			case hold := <-Static.RequestChan:
				mapRequestChan <- hold
			}
		}
	} else if Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				sqlRequestChan <- hold
			}
		}
	} else {
		Static.LogChan <- "No database in use."
	}
}

// Holds information for request from database
type Request struct {
	Id           string
	Data         interface{}
	ResponseChan chan Response
}

// Holds information for response from database
type Response struct {
	Id   string
	Db   string
	Data interface{}
}
type WriteResponse struct {
	Complete bool
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	HandleDb(chan Request)
	Shutdown()
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id      string
	Db      map[string]map[string]interface{}
	Workers map[string]MapWorker
	Busy    bool
}

// Handle data MapDb requests
func (db MapDb) HandleDb(mapChan chan Request) {
	// Initialize map database
	db.Db = make(map[string]map[string]interface{})

	// Loop until shutdown
	for {
		db.Busy = false
		select {
		case hold := <-mapChan:

			if len(hold.Id) < Static.Config.CacheSize {
				// Mark database as busy
				db.Busy = true

				l := len(hold.Id)
				s := Static.Config.CacheSize
				key := hold.Id[l-s : l-1]
				switch {
				// This logic needs to be refactored so that the number of shards can be
				// determined via the config file, which will define the number of digits
				// used from the ID for the shard name
				case hold.Data == nil:
					// check if key exits
					_, ok := db.Db[key]
					if !ok {
						// if key does not exist, make a new map with the given key
						db.Db[key] = make(map[string]interface{})
					}
					go new(MapWorker).Read(hold, db.Db[key])
				case hold.Data != nil:
					// check if key exits
					_, ok := db.Db[key]
					if !ok {
						// if key does not exist, make a new map with the given key
						db.Db[key] = make(map[string]interface{})
					}
					go new(MapWorker).Write(hold, db.Db[key])
				}
			} else {
				var err ErrorRes
				err.ErrLocation = "HandleDb"
				err.Time = time.Now().UnixNano()
				var res Response
				res.Id = hold.Id
				res.Data = err
				Static.ErrChan <- res
			}

		case <-Static.ShutdownChan:
			db.Busy = false
			break
		}
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

// SqlDb is a Sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb(sqlChan chan Request) {
}

// Shutdown SqlDb
func (db SqlDb) Shutdown() {
}
