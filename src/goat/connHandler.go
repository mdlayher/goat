package goat

import (
	"net"
	"net/http"
	"strings"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(l net.Listener)
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle incoming HTTP connections and serve
func (h HttpConnHandler) Handle(l net.Listener) {
	// Set up HTTP routes for handling functions
	http.HandleFunc("/", parseHttp)

	// Serve HTTP requests
	err := http.Serve(l, nil)
	if err != nil {
		Static.LogChan <- err.Error()
	}
}

// Parse incoming HTTP connections before making tracker calls
func parseHttp(w http.ResponseWriter, r *http.Request) {
	// Count incoming connections
	Static.Http.Current++
	Static.Http.Total++

	// Create channel to return bencoded response to client on
	resChan := make(chan []byte)

	// Parse querystring
	querystring := r.URL.Query()

	// Flatten arrays into single values
	query := map[string]string{}
	for k, v := range querystring {
		query[k] = v[0]
	}

	// Check if IP was previously set
	if _, ok := query["ip"]; !ok {
		// If no IP set, detect and store it in query map
		query["ip"] = strings.Split(r.RemoteAddr, ":")[0]
	}

	// Handle tracker functions via different URLs
	switch r.URL.Path {
	// Tracker announce
	case "/announce":
		go TrackerAnnounce(query, resChan)
	// Tracker status
	case "/status":
		go GetServerStatus(resChan)
	// Any undefined handlers
	default:
		go TrackerError(resChan, "Malformed announce")
	}

	// Add header to identify goat
	w.Header().Add("Server", APP+"-git")

	// Wait for response, and send it when ready
	w.Write(<-resChan)
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l net.Listener) {
}
