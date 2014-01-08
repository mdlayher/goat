package goat

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// Any error responses from the API
type APIError struct {
	Error string `json:"error"`
}

// API router entry point
func APIRouter(w http.ResponseWriter, r *http.Request) {
	// Split request path
	urlArr := strings.Split(r.URL.Path, "/")

	// Verify API method set
	if len(urlArr) < 3 {
		http.Error(w, string(APIErrorResponse("No API call")), 404)
		return
	}

	// API response chan
	apiChan := make(chan []byte)

	// Choose API method
	switch urlArr[2] {
	// Server status
	case "status":
		go GetStatusJSON(apiChan)
	// Return error response
	default:
		http.Error(w, string(APIErrorResponse("Undefined API call")), 404)
		return
	}

	w.Write(<-apiChan)
	close(apiChan)
	return
}

// Return an API error
func APIErrorResponse(msg string) []byte {
	res := APIError{
		msg,
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return out
}
