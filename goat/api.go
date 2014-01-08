package goat

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// APIError represents an error response from the API
type APIError struct {
	Error string `json:"error"`
}

// APIRouter handles the routing of HTTP API requests
func APIRouter(w http.ResponseWriter, r *http.Request) {
	// Split request path
	urlArr := strings.Split(r.URL.Path, "/")

	// Verify API method set
	if len(urlArr) < 3 {
		http.Error(w, string(APIErrorResponse("No API call")), 404)
		return
	}

	// Check for an ID
	ID := -1
	if len(urlArr) == 4 {
		i, err := strconv.Atoi(urlArr[3])
		if err == nil {
			ID = i
		} else {
			log.Println(err.Error())
			ID = -1
		}
	}

	// API response chan
	apiChan := make(chan []byte)

	// Choose API method
	switch urlArr[2] {
	// Files on tracker
	case "files":
		go GetFilesJSON(ID, apiChan)
	// Server status
	case "status":
		go GetStatusJSON(apiChan)
	// Return error response
	default:
		http.Error(w, string(APIErrorResponse("Undefined API call")), 404)
		close(apiChan)
		return
	}

	w.Write(<-apiChan)
	close(apiChan)
	return
}

// APIErrorResponse return an APIError as JSON
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
