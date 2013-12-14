package goat

type IoWorker interface {
	Read(request Request)
	Write(request Request)
}

type MapWorker struct {
}

type SqlWorker struct {
}

// map read and write

//takes a pointer to a shard map and a request and replies with a response struct
func (m MapWorker) Read(shard map[string]interface{}, request Request) {
	var response Response
	response.Data = <-shard[request.Id]
	response.Id = request.Id
	response.Db = request.Id[len(request.Id)-3 : len(request.Id)-1]
	close(request.ResponseChan)
}

// takes a pointer to a map shard and writes the data from a request and
// sends back a conformation that the write has taken place.
func (m MapWorker) Write(shard map[string]interface{}, request Request) {
	var comp WriteResponse
	var response Response
	shard[request.Id] = request.Data
	comp.Complete = true
	response.Data = comp
	response.Db = request.Id[len(request.Id)-3 : len(request.Id)-1]
	response.Id = request.Id
	request.ResponseChan <- response

}

// sql read and write
func (s SqlWorker) Read(request Request) {

}
func (s SqlWorker) Write(Request Request) {

}
