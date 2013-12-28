package goat

type IoWorker interface {
	Lookup(MapDb, Request)
	UpdateLookup(MapDb, Request)
}
type Worker struct {
}

// looks up where in the dbStor the requested data is stored
func (worker Worker) Lookup(db MapDb, id string) []*interface{} {
	return db.MapLookup[id]

}

// reqplaces the pointer in the lookup with a new pointer to the new data
func (worker Worker) UpdateLookup(db MapDb, id string, data interface{}) {
	db.MapLookup[id][0] = &data
}

func (worker Worker) AppendLookup(db MapDb, id string, data interface{}) {
	db.MapLookup[id] = append(db.MapLookup[id], &data)
}
