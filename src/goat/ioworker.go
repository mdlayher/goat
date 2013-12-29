package goat

import (
	"errors"
)

type IoWorker struct {
	db *MapDb
}

// looks up where in the dbStor the requested data is stored
func (worker IoWorker) Lookup(id string) (interface{}, error) {
	if val, ok := worker.db.MapLookup[id]; ok {
		return val, nil
	} else {
		return nil, errors.New("Key does not exist in this map")
	}

}

// replaces the pointer in the lookup with a new pointer to the new data
func (worker IoWorker) UpdateLookup(id string, update interface{}) {
	worker.db.MapLookup[id] = &update
}

// append a pointer on the end of the lookup
func (worker IoWorker) AppendLookup(id string, data *interface{}) {
	hold := worker.db.MapLookup[id].(Relation)
	hold.Index = append(hold.Index, data)
	worker.db.MapLookup[id] = hold

}

// remove a pointer from lookup
func (worker IoWorker) RemoveLookup(id string) {
	worker.db.MapLookup[id] = nil
}

// replace the data for a given id where ever it might be(which is identified by the Lookup function)
func (worker IoWorker) UpdateStor(id string, data interface{}, location map[string]interface{}) {
	location[id] = data
}
