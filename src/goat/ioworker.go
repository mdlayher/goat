package goat

type IoWorker struct {
	db *MapDb
}

// looks up where in the dbStor the requested data is stored
func (worker IoWorker) Lookup(id string) []*interface{} {
	return worker.db.MapLookup[id]

}

// reqplaces the pointer in the lookup with a new pointer to the new data
func (worker IoWorker) UpdateLookup(id string, data interface{}) {
	worker.db.MapLookup[id][0] = &data
}

// append a pointer on the end of the lookup
func (worker IoWorker) AppendLookup(id string, data interface{}) {
	worker.db.MapLookup[id] = append(worker.db.MapLookup[id], &data)
}

func (worker IoWorker) RemoveLookup(id string) {
	worker.db.MapLookup[id] = nil
}
