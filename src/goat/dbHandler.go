package goat

func DbManager(RequestChan chan Request) {
	// channels
	SqlRequestChan := make(chan Request)
	MapRequestChan := make(chan Request, 100)

	// launch databases
	if Static.Config.Map {
		go new(MapDb).HandleDb(RequestChan)
		Static.LogChan <- "MapDb instance launched"
	}
	if Static.Config.Sql {
		go new(SqlDb).HandleDb(RequestChan)
		Static.LogChan <- "SqlDb instance launched"
	}

	if Static.Config.Map && Static.Config.Sql {
		for {
			select {
			case hold := <-RequestChan:
				if hold.MapOnly {
					MapRequestChan <- hold
				} else {
					MapRequestChan <- hold
					SqlRequestChan <- hold
				}
			}
		}
	} else if Static.Config.Map {
		for {
			select {
			case hold := <-RequestChan:
				MapRequestChan <- hold
			}
		}
	} else if Static.Config.Sql {
		for {
			select {
			case hold := <-RequestChan:
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
	Read         bool
	Write        bool
	DbOnly       bool
	MapOnly      bool
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
	for {
		select {
		case hold := <-RequestChan:
			switch {
			// This logic needs to be refactored so that the number of shards can be
			// determined via the config file, which will define the number of digits
			// used from the ID for the shard name
			case hold.Read:
				go new(MapWorker).Read(hold, db.Db[hold.Id[(len(hold.Id)-3):(len(hold.Id)-1)]])
			case hold.Write:
				go new(MapWorker).Write(hold, db.Db[hold.Id[len(hold.Id)-3:len(hold.Id)-1]])
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
