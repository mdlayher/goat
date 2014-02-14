package api

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Error represents an error response from the API
type Error struct {
	Error string `json:"error"`
}

// Router handles the routing of HTTP API requests
func Router(w http.ResponseWriter, r *http.Request) {
	// API allows the following HTTP methods:
	//   - GET: read-only access to data
	//   - POST: create a new item via an API endpoint
	if r.Method != "GET" && r.Method != "POST" {
		http.Error(w, ErrorResponse("Method not allowed"), 405)
		return
	}

	// Log API calls
	log.Printf("API: [http %s] %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)

	// Split request path
	urlArr := strings.Split(r.URL.Path, "/")

	// Verify API method set
	if len(urlArr) < 3 {
		http.Error(w, ErrorResponse("No API call"), 404)
		return
	}

	// API method
	apiMethod := urlArr[2]

	// Response buffer
	res := make([]byte, 0)

	// HTTP GET
	if r.Method == "GET" {
		// Default value retrieves all records
		ID := -1

		// Check for an ID
		if len(urlArr) == 4 {
			i, err := strconv.Atoi(urlArr[3])
			if err != nil || i < 1 {
				http.Error(w, ErrorResponse("Invalid integer ID"), 400)
				return
			}

			ID = i
		}

		// Check for error
		var err error

		// Choose API method
		switch apiMethod {
		// Files on tracker
		case "files":
			res, err = getFilesJSON(ID)
		// Server status
		case "status":
			res, err = getStatusJSON()
		// Users registered to tracker
		case "users":
			res, err = getUsersJSON(ID)
		// Return error response
		default:
			http.Error(w, ErrorResponse("Undefined API call: GET /api/"+apiMethod), 404)
			return
		}

		// Check for server error
		if err != nil {
			log.Println(err.Error())
			http.Error(w, ErrorResponse("API failure: GET /api/"+apiMethod), 500)
			return
		}
	}

	// HTTP POST
	if r.Method == "POST" {
		// Attempt to read the request body
		body, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			http.Error(w, ErrorResponse("Malformed request body"), 400)
		}

		// Check for client string and server error
		// Note: client error is string to satisfy golint and for client error consistency
		var clientErr string
		var serverErr error

		// Choose API method
		switch apiMethod {
		// Users registered to tracker
		case "users":
			// Attempt to create user from JSON
			clientErr, serverErr = postUsersJSON(body)
		// Return error response
		default:
			http.Error(w, ErrorResponse("Undefined API call: POST /api/"+apiMethod), 404)
			return
		}

		// Check for client string error
		if clientErr != "" {
			http.Error(w, ErrorResponse(clientErr), 400)
			return
		}

		// Check for server error
		if serverErr != nil {
			log.Println(serverErr.Error())
			http.Error(w, ErrorResponse("API failure: POST /api/"+apiMethod), 500)
			return
		}

		// Return HTTP 204 on success
		http.Error(w, "", 204)
		return
	}

	// If requested, compress response using gzip
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.Header().Add("Content-Encoding", "gzip")

		// Write gzip'd response
		gz := gzip.NewWriter(w)
		if _, err := gz.Write(res); err != nil {
			log.Println(err.Error())
			return
		}

		if err := gz.Close(); err != nil {
			log.Println(err.Error())
		}

		return
	}

	if _, err := w.Write(res); err != nil {
		log.Println(err.Error())
	}

	return
}

// ErrorResponse returns an Error as JSON
func ErrorResponse(msg string) string {
	res := Error{
		msg,
	}

	out, err := json.Marshal(res)
	if err != nil {
		log.Println(err.Error())
		return `{"error":"` + msg + `"}`
	}

	return string(out)
}
