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

	// Store current passkey URL
	var passkey string = ""
	url := r.URL.Path

	// Check for passkey present in URL, form: /ABDEF/announce
	urlArr := strings.Split(url, "/")
	url = urlArr[1]
	if len(urlArr) == 3 {
		passkey = urlArr[1]
		url = urlArr[2]
	}

	// Add header to identify goat
	w.Header().Add("Server", APP+"-git")

	// Check if server is configured for passkey announce
	if Static.Config.Passkey && passkey == "" {
		go TrackerError(resChan, "No passkey found in announce URL")

		// Wait for response, and send it when ready
		w.Write(<-resChan)
		close(resChan)
		return
	}

	// Handle tracker functions via different URLs
	switch url {
	// Tracker announce
	case "announce":
		// Validate required parameter input
		required := []string{"info_hash", "peer_id", "ip", "port", "uploaded", "downloaded", "left"}
		for _, r := range required {
			if _, ok := query[r]; !ok {
				go TrackerError(resChan, "Missing required parameter: "+r)

				// Wait for response, and send it when ready
				w.Write(<-resChan)
				close(resChan)
				return
			}
		}

		// Perform tracker announce
		go TrackerAnnounce(passkey, query, resChan)
	// Tracker status
	case "status":
		go GetServerStatus(resChan)
	// Any undefined handlers
	default:
		go TrackerError(resChan, "Malformed announce")
	}

	// Wait for response, and send it when ready
	w.Write(<-resChan)
	close(resChan)
	return
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l net.Listener) {
}
