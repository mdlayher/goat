package goat

import (
	"time"
)

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
func (m MapWorker) Read(request Request, shard map[string]interface{}) {
	var response Response
	_, ok := shard[request.Id]
	if ok {
		response.Data = shard[request.Id]
		response.Id = request.Id
		response.Db = "map"
		request.ResponseChan <- response
	} else {
		responseChan := make(chan Response)
		var dbReadRequest Request
		dbReadRequest.ResponseChan = responseChan
		dbReadRequest.Id = request.Id
		Static.RequestChan <- dbReadRequest
		select {
		case response := <-responseChan:
			var mapWriteRequest Request
			mapWriteRequest.Id = request.Id
			mapWriteRequest.Data = response.Data
			Static.RequestChan <- mapWriteRequest
		case <-time.After(100 * time.Millisecond):
			close(responseChan)
		}
	}
	<-request.ResponseChan
}

// takes a pointer to a map shard and writes the data from a request and
// sends back a conformation that the write has taken place.
func (m MapWorker) Write(request Request, shard map[string]interface{}) {
	var comp WriteResponse
	var response Response
	shard[request.Id] = request.Data
	comp.Complete = true
	response.Data = comp
	response.Db = "map"
	response.Id = request.Id
	request.ResponseChan <- response
}

// sql read and write
func (s SqlWorker) Read(request Request) {
}
func (s SqlWorker) Write(Request Request) {
}
