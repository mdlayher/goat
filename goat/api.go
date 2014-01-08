package goat

import (
	"encoding/json"
	"log"
	"net/http"
)

// Any error responses from the API
type APIError struct {
	Error string `json:"error"`
}

// API router entry point
func APIRouter(r *http.Request, resChan chan []byte) {
	res := APIError{
		"Hello World!",
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		resChan <- nil
	}

	// Return response
	resChan <- out
}
