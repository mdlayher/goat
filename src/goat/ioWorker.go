package goat

type IoWorker interface {
	Read(ResponseChan chan Response, request Request)
	Write(ResponseChan chan Response, request Request)
}
