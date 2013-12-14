package goat

func DbManager(RequestChan chan Request) {

	SqlRequestChan := make(chan Request)
	MapRequesChan := make(chan Request)

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
				MapRequesChan <- hold
				SqlRequestChan <- hold
			}
		}
	} else if Static.Config.Map {
		for {
			select {
			case hold := <-RequestChan:
				MapRequesChan <- hold

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
	Read, Write  bool
	Data         interface{}
	ResponseChan chan Response
}

// Holds information for response from database
type Response struct {
	Id, Db string
	Data   interface{}
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	HandleDb(RequestChan chan Request)
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id string
	Db map[string]map[string]interface{}
}

// Handle data MapDb requests
func (m MapDb) HandleDb(RequestChan chan Request) {

}

// SqlDb is a sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb(RequestChan chan Request) {

}
