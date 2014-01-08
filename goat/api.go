package goat

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(urlArr)

	// Verify API method set
	if len(urlArr) < 3 {
		http.Error(w, string(APIErrorResponse("No API call")), 404)
		return
	}

	method := urlArr[2]
	fmt.Println(method)

	// Return error response
	http.Error(w, string(APIErrorResponse("Undefined API call")), 404)
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
