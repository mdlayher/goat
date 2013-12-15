package goat

func DbManager() {
	// channels
	SqlRequestChan := make(chan Request)
	MapRequestChan := make(chan Request, 100)

	// launch databases
	if Static.Config.Map {
		go new(MapDb).HandleDb(Static.RequestChan)
		Static.LogChan <- "MapDb instance launched"
	}
	if Static.Config.Sql {
		go new(SqlDb).HandleDb(Static.RequestChan)
		Static.LogChan <- "SqlDb instance launched"
	}

	if Static.Config.Map && Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				if hold.Data == nil {
					MapRequestChan <- hold
				} else {
					MapRequestChan <- hold
					SqlRequestChan <- hold
				}
			case hold := <-Static.PersistentChan:
				SqlRequestChan <- hold
			}
		}
	} else if Static.Config.Map {
		for {
			select {
			case hold := <-Static.RequestChan:
				MapRequestChan <- hold
			}
		}
	} else if Static.Config.Sql {
		for {
			select {
			case hold := <-Static.RequestChan:
				SqlRequestChan <- hold
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
	HandleDb(RequestChan chan Request)
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id      string
	Db      map[string]map[string]interface{}
	Workers map[string]MapWorker
}

// Handle data MapDb requests
func (db MapDb) HandleDb(RequestChan chan Request) {
	db.Db = make(map[string]map[string]interface{})
	for {
		select {
		case hold := <-RequestChan:
			l := len(hold.Id)
			s := Static.Config.CacheSize
			key := hold.Id[l-s : l-1]
			switch {
			// This logic needs to be refactored so that the number of shards can be
			// determined via the config file, which will define the number of digits
			// used from the ID for the shard name
			case hold.Data == nil:
				_, ok := db.Db[key]
				if !ok {
					db.Db[key] = make(map[string]interface{})
				}
				go new(MapWorker).Read(hold, db.Db[key])
			case hold.Data != nil:
				_, ok := db.Db[key]
				if !ok {
					db.Db[key] = make(map[string]interface{})
				}
				go new(MapWorker).Write(hold, db.Db[key])
			}
		}
	}
}

// SqlDb is a Sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb(RequestChan chan Request) {
}
