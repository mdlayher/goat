package goat

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

// Handle incoming HTTP connections and serve
func handleHTTP(l net.Listener, httpDoneChan chan bool) {
	// Create shutdown function
	go func(l net.Listener, httpDoneChan chan bool) {
		// Wait for done signal
		static.ShutdownChan <- <-static.ShutdownChan

		// Close listener
		l.Close()
		httpDoneChan <- true
	}(l, httpDoneChan)

	// Log API configuration
	if static.Config.API {
		log.Println("API functionality enabled")
	}

	// Set up HTTP routes for handling functions
	http.HandleFunc("/", parseHTTP)

	// Serve HTTP requests
	http.Serve(l, nil)
}

// Parse incoming HTTP connections before making tracker calls
func parseHTTP(w http.ResponseWriter, r *http.Request) {
	// Count incoming connections
	atomic.AddInt64(&static.HTTP.Current, 1)
	atomic.AddInt64(&static.HTTP.Total, 1)

	// Add header to identify goat
	w.Header().Add("Server", fmt.Sprintf("%s/%s", App, Version))

	// Parse querystring into a Values map
	query := r.URL.Query()

	// Check if IP was previously set
	if query.Get("ip") == "" {
		// If no IP set, detect and store it in query map
		query.Set("ip", strings.Split(r.RemoteAddr, ":")[0])
	}

	// Store current URL path
	url := r.URL.Path

	// Split URL into segments
	urlArr := strings.Split(url, "/")

	// If configured, Detect if client is making an API call
	url = urlArr[1]
	if url == "api" {
		// Output JSON
		w.Header().Add("Content-Type", "application/json")

		// API enabled
		if static.Config.API {
			// Handle API calls, output JSON
			apiRouter(w, r)
			return
		} else {
			http.Error(w, string(apiErrorResponse("API is currently disabled")), 503)
			return
		}
	}

	// Detect if passkey present in URL
	var passkey string
	if len(urlArr) == 3 {
		passkey = urlArr[1]
		url = urlArr[2]
	}

	// Verify that torrent client is advertising its User-Agent, so we can use a whitelist
	if _, ok := r.Header["User-Agent"]; !ok {
		w.Write(httpTrackerError("Your client is not identifying itself"))
		return
	}

	client := r.Header["User-Agent"][0]

	// If configured, verify that torrent client is on whitelist
	if static.Config.Whitelist {
		whitelist := new(whitelistRecord).Load(client, "client")
		if whitelist == (whitelistRecord{}) || !whitelist.Approved {
			w.Write(httpTrackerError("Your client is not whitelisted"))

			// Block things like browsers and web crawlers, because they will just clutter up the table
			if strings.Contains(client, "Mozilla") || strings.Contains(client, "Opera") {
				return
			}

			// Insert unknown clients into list for later approval
			if whitelist == (whitelistRecord{}) {
				whitelist.Client = client
				whitelist.Approved = false

				log.Printf("whitelist: detected new client '%s', awaiting manual approval", client)

				go whitelist.Save()
			}

			return
		}
	}

	// Put client in query map
	query.Set("client", client)

	// Check if server is configured for passkey announce
	if static.Config.Passkey && passkey == "" {
		w.Write(httpTrackerError("No passkey found in announce URL"))
		return
	}

	// Validate passkey if needed
	user := new(userRecord).Load(passkey, "passkey")
	if static.Config.Passkey && user == (userRecord{}) {
		w.Write(httpTrackerError("Invalid passkey"))
		return
	}

	// Put passkey in query map
	query.Set("passkey", user.Passkey)

	// Mark client as HTTP
	query.Set("udp", "0")

	// Get user's total number of active torrents
	seeding := user.Seeding()
	leeching := user.Leeching()
	if seeding == -1 || leeching == -1 {
		w.Write(httpTrackerError("Failed to calculate active torrents"))
		return
	}

	// Verify that client has not exceeded this user's torrent limit
	activeSum := seeding + leeching
	if user.TorrentLimit < activeSum {
		w.Write(httpTrackerError(fmt.Sprintf("Exceeded active torrent limit: %d > %d", activeSum, user.TorrentLimit)))
		return
	}

	// Create channel to return response to client
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
			if query.Get(r) == "" {
				w.Write(httpTrackerError("Missing required parameter: " + r))
				close(resChan)
				return
			}
		}

		// Check for all valid integers
		for _, r := range reqInt {
			if query.Get(r) != "" {
				_, err := strconv.Atoi(query.Get(r))
				if err != nil {
					w.Write(httpTrackerError("Invalid integer parameter: " + r))
					close(resChan)
					return
				}
			}
		}

		// Only allow compact announce
		if query.Get("compact") == "" || query.Get("compact") != "1" {
			w.Write(httpTrackerError("Your client does not support compact announce"))
			close(resChan)
			return
		}

		// Perform tracker announce
		go trackerAnnounce(user, query, nil, resChan)
	// Tracker scrape
	case "scrape":
		go trackerScrape(user, query, resChan)
	// Any undefined handlers
	default:
		w.Write(httpTrackerError("Malformed announce"))
		close(resChan)
		return
	}

	// Wait for response, and send it when ready
	w.Write(<-resChan)
	close(resChan)
	return
}
