package goat

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle incoming HTTP connections and serve
func (h HttpConnHandler) Handle(l net.Listener, httpDoneChan chan bool) {
	// Create shutdown function
	go func(l net.Listener, httpDoneChan chan bool) {
		// Wait for done signal
		<-Static.ShutdownChan

		// Close listener
		l.Close()
		httpDoneChan <- true
	}(l, httpDoneChan)

	// Set up HTTP routes for handling functions
	http.HandleFunc("/", parseHttp)

	// Serve HTTP requests
	http.Serve(l, nil)
}

// Parse incoming HTTP connections before making tracker calls
func parseHttp(w http.ResponseWriter, r *http.Request) {
	// Count incoming connections
	atomic.AddInt64(&Static.Http.Current, 1)
	atomic.AddInt64(&Static.Http.Total, 1)

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
	w.Header().Add("Server", fmt.Sprintf("%s/%s", APP, VERSION))

	// Verify that torrent client is advertising its User-Agent, so we can use a whitelist
	if _, ok := r.Header["User-Agent"]; !ok {
		w.Write(HttpTrackerError("Your client is not identifying itself"))
		return
	}

	client := r.Header["User-Agent"][0]

	// If configured, verify that torrent client is on whitelist
	if Static.Config.Whitelist {
		whitelist := new(WhitelistRecord).Load(client, "client")
		if whitelist == (WhitelistRecord{}) || !whitelist.Approved {
			w.Write(HttpTrackerError("Your client is not whitelisted"))

			// Insert unknown clients into list for later approval
			if whitelist == (WhitelistRecord{}) {
				whitelist.Client = client
				whitelist.Approved = false

				Static.LogChan <- fmt.Sprintf("whitelist: detected new client '%s', awaiting manual approval", client)

				go whitelist.Save()
			}

			return
		}
	}

	// Put client in query map
	query["client"] = client

	// Check if server is configured for passkey announce
	if Static.Config.Passkey && passkey == "" {
		w.Write(HttpTrackerError("No passkey found in announce URL"))
		return
	}

	// Validate passkey if needed
	user := new(UserRecord).Load(passkey, "passkey")
	if Static.Config.Passkey && user == (UserRecord{}) {
		w.Write(HttpTrackerError("Invalid passkey"))
		return
	}

	// Put passkey in query map
	query["passkey"] = user.Passkey

	// Mark client as HTTP
	query["udp"] = "0"

	// Create channel to return bencoded response to client on
	resChan := make(chan []byte)

	// Handle tracker functions via different URLs
	switch url {
	// Tracker announce
	case "announce":
		// Validate required parameter input
		required := []string{"info_hash", "ip", "port", "uploaded", "downloaded", "left"}
		// Validate required integer input
		reqInt := []string{"port", "uploaded", "downloaded", "left"}

		// Check for required parameters
		for _, r := range required {
			if _, ok := query[r]; !ok {
				w.Write(HttpTrackerError("Missing required parameter: " + r))
				close(resChan)
				return
			}
		}

		// Check for all valid integers
		for _, r := range reqInt {
			if _, ok := query[r]; ok {
				_, err := strconv.Atoi(query[r])
				if err != nil {
					w.Write(HttpTrackerError("Invalid integer parameter: " + r))
					close(resChan)
					return
				}
			}
		}

		// Only allow compact announce
		if _, ok := query["compact"]; !ok || query["compact"] != "1" {
			w.Write(HttpTrackerError("Your client does not support compact announce"))
			close(resChan)
			return
		}

		// Perform tracker announce
		go TrackerAnnounce(user, query, nil, resChan)
	// Tracker status
	case "status":
		w.Header().Add("Content-Type", "application/json")
		go GetStatusJson(resChan)
	// Any undefined handlers
	default:
		w.Write(HttpTrackerError("Malformed announce"))
		close(resChan)
		return
	}

	// Wait for response, and send it when ready
	w.Write(<-resChan)
	close(resChan)
	return
}
