package goat

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/mdlayher/goat/goat/api"
	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
	"github.com/mdlayher/goat/goat/tracker"
)

// Handle incoming HTTP connections and serve
func handleHTTP(l net.Listener, sendChan chan bool, recvChan chan bool) {
	// Create shutdown function
	go func(l net.Listener, sendChan chan bool, recvChan chan bool) {
		// Wait for done signal
		<-sendChan

		// Close listener
		if err := l.Close(); err != nil {
			log.Println(err.Error())
		}

		log.Println("HTTP(S) listener stopped")
		recvChan <- true
	}(l, sendChan, recvChan)

	// Log API configuration
	if common.Static.Config.API {
		log.Println("API functionality enabled")
	}

	// Serve HTTP requests
	if err := http.Serve(l, nil); err != nil {
		// Ignore connection closing error, caused by stopping listener
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Println("Could not serve HTTP(S), exiting now")
			panic(err)
		}
	}
}

// Parse incoming HTTP connections before making tracker calls
func parseHTTP(w http.ResponseWriter, r *http.Request) {
	// Create a tracker to handle this client
	httpTracker := tracker.HTTPTracker{}

	// Add header to identify goat
	w.Header().Add("Server", fmt.Sprintf("%s/%s", App, Version))

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
		if common.Static.Config.API {
			// Count incoming connections
			atomic.AddInt64(&common.Static.API.Minute, 1)
			atomic.AddInt64(&common.Static.API.HalfHour, 1)
			atomic.AddInt64(&common.Static.API.Hour, 1)
			atomic.AddInt64(&common.Static.API.Total, 1)

			// API authentication
			auth := new(api.BasicAuthenticator).Auth(r)
			if !auth {
				http.Error(w, api.ErrorResponse("Authentication failed"), 401)
				return
			}

			// Handle API calls, output JSON
			api.Router(w, r)
			return
		}

		http.Error(w, api.ErrorResponse("API is currently disabled"), 503)
		return
	}

	// Count incoming connections
	atomic.AddInt64(&common.Static.HTTP.Minute, 1)
	atomic.AddInt64(&common.Static.HTTP.HalfHour, 1)
	atomic.AddInt64(&common.Static.HTTP.Hour, 1)
	atomic.AddInt64(&common.Static.HTTP.Total, 1)

	// Check for maintenance mode
	if common.Static.Maintenance {
		// Return tracker error with maintenance message
		if _, err := w.Write(httpTracker.Error("Maintenance: " + common.Static.StatusMessage)); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Detect if passkey present in URL
	var passkey string
	if len(urlArr) == 3 {
		passkey = urlArr[1]
		url = urlArr[2]
	}

	// Make sure URL is valid torrent function
	if url != "announce" && url != "scrape" {
		if _, err := w.Write(httpTracker.Error("Malformed announce")); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Verify that torrent client is advertising its User-Agent, so we can use a whitelist
	if r.Header.Get("User-Agent") == "" {
		if _, err := w.Write(httpTracker.Error("Your client is not identifying itself")); err != nil {
			log.Println(err.Error())
		}

		return
	}

	client := r.Header.Get("User-Agent")

	// If configured, verify that torrent client is on whitelist
	if common.Static.Config.Whitelist {
		whitelist, err := new(data.WhitelistRecord).Load(client, "client")
		if err != nil {
			log.Println(err.Error())
		}

		if whitelist == (data.WhitelistRecord{}) || !whitelist.Approved {
			if _, err := w.Write(httpTracker.Error("Your client is not whitelisted")); err != nil {
				log.Println(err.Error())
			}

			// Block things like browsers and web crawlers, because they will just clutter up the table
			if strings.Contains(client, "Mozilla") || strings.Contains(client, "Opera") {
				return
			}

			// Insert unknown clients into list for later approval
			if whitelist == (data.WhitelistRecord{}) {
				whitelist.Client = client
				whitelist.Approved = false

				log.Printf("whitelist: detected new client '%s', awaiting manual approval", client)

				// Save whitelist record asynchronously
				go func(whitelist data.WhitelistRecord) {
					if err := whitelist.Save(); err != nil {
						log.Println(err.Error())
					}
				}(whitelist)
			}

			return
		}
	}

	// Parse querystring into a Values map
	query := r.URL.Query()

	// Check if IP was previously set
	if query.Get("ip") == "" {
		// If no IP set, detect and store it in query map
		query.Set("ip", strings.Split(r.RemoteAddr, ":")[0])
	}

	// Put client in query map
	query.Set("client", client)

	// Check if server is configured for passkey announce
	if common.Static.Config.Passkey && passkey == "" {
		if _, err := w.Write(httpTracker.Error("No passkey found in announce URL")); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Validate passkey if needed
	user, err := new(data.UserRecord).Load(passkey, "passkey")
	if err != nil || (common.Static.Config.Passkey && user == (data.UserRecord{})) {
		if err != nil {
			log.Println(err.Error())
		}

		if _, err := w.Write(httpTracker.Error("Invalid passkey")); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Put passkey in query map
	query.Set("passkey", user.Passkey)

	// Mark client as HTTP
	query.Set("udp", "0")

	// Get user's number of active torrents
	seeding, err := user.Seeding()
	leeching, err2 := user.Leeching()
	if err != nil || err2 != nil {
		if _, err := w.Write(httpTracker.Error("Failed to calculate active torrents")); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Verify that client has not exceeded this user's torrent limit
	activeSum := seeding + leeching
	if user.TorrentLimit < activeSum {
		msg := fmt.Sprintf("Exceeded active torrent limit: %d > %d", activeSum, user.TorrentLimit)
		if _, err := w.Write(httpTracker.Error(msg)); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Tracker announce
	if url == "announce" {
		// Validate required parameter input
		required := []string{"info_hash", "ip", "port", "uploaded", "downloaded", "left"}
		// Validate required integer input
		reqInt := []string{"port", "uploaded", "downloaded", "left"}

		// Check for required parameters
		for _, r := range required {
			if query.Get(r) == "" {
				if _, err := w.Write(httpTracker.Error("Missing required parameter: " + r)); err != nil {
					log.Println(err.Error())
				}

				return
			}
		}

		// Check for all valid integers
		for _, r := range reqInt {
			if query.Get(r) != "" {
				_, err := strconv.Atoi(query.Get(r))
				if err != nil {
					if _, err := w.Write(httpTracker.Error("Invalid integer parameter: " + r)); err != nil {
						log.Println(err.Error())
					}

					return
				}
			}
		}

		// Only allow compact announce
		if query.Get("compact") == "" || query.Get("compact") != "1" {
			if _, err := w.Write(httpTracker.Error("Your client does not support compact announce")); err != nil {
				log.Println(err.Error())
			}

			return
		}

		// NOTE: currently, we do not bother using gzip to compress the tracker announce response
		// This is done for two reasons:
		// 1) Clients may or may not support gzip in the first place
		// 2) gzip may actually make announce response larger, as per testing in What.CD's ocelot

		// Perform tracker announce
		if _, err := w.Write(tracker.Announce(httpTracker, user, query)); err != nil {
			log.Println(err.Error())
		}

		return
	}

	// Tracker scrape
	if url == "scrape" {
		// Check for required parameter info_hash
		if query.Get("info_hash") == "" {
			if _, err := w.Write(httpTracker.Error("Missing required parameter: info_hash")); err != nil {
				log.Println(err.Error())
			}

			return
		}

		if _, err := w.Write(tracker.Scrape(httpTracker, query)); err != nil {
			log.Println(err.Error())
		}

		return
	}

	return
}
