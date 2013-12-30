package goat

import (
	"time"
)

type IoWorker struct {
	myLookup MapStorage
	db       *MapDb
}

// looks up where in the dbStor the requested data is stored
func (worker IoWorker) Lookup(id string) interface{} {
	if val, ok := worker.myLookup[id]; ok {
		return val
	} else {
		var err ErrorRes
		err.ErrLocation = "Lookup"
		err.Error = "id not found in map"
		err.Time = time.Now().Unix()
		var res Response
		res.Data = err
		res.Id = id
		Static.ErrChan <- res
		return nil
	}

}

// replaces the pointer in the lookup with a new pointer to the new data
func (worker IoWorker) UpdateLookup(id string, update interface{}) {
	worker.db.MapLookup[id] = &update
}

// Append new pointer to the end of a relation or create a new one if relation is not found
func (worker IoWorker) AppendRelation(id string, data *interface{}) {

	if hold, ok := worker.myLookup[id].(Relation); ok {
		hold.Index = append(hold.Index, data)
		worker.myLookup[id] = hold
	} else {
		n := new(Relation)
		n.Index[0] = data
		worker.myLookup[id] = n
	}
}

// Removes a given item from the index of a relation
func (worker IoWorker) RemoveRelation(id string, data *interface{}) {
	if hold, ok := worker.myLookup[id].(Relation); ok {
		for i := 0; i < len(hold.Index); i++ {
			if hold.Index[i] == data {
				hold.Index = append(hold.Index[:i], hold.Index[i+1:]...)
			}
		}
	} else {
		var err ErrorRes
		err.ErrLocation = "RemoveRelation"
		err.Error = "id not found in map"
		err.Time = time.Now().Unix()
		var res Response
		res.Data = err
		res.Id = id
		Static.ErrChan <- res
	}
}

// remove a pointer from lookup
func (worker IoWorker) RemoveLookup(id string) {
	worker.myLookup[id] = nil
}

// replace the data for a given id where ever it might be(which is identified by the Lookup function)
func (worker IoWorker) UpdateStor(id string, data interface{}, location map[string]interface{}) {
	location[id] = data
}
