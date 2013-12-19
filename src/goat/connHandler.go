package goat

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
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
		go TrackerAnnounce(user, query, resChan)
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

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l *net.UDPConn, udpDoneChan chan bool) {
	// Create shutdown function
	go func(l *net.UDPConn, udpDoneChan chan bool) {
		// Wait for done signal
		<-Static.ShutdownChan

		// Close listener
		l.Close()
		udpDoneChan <- true
	}(l, udpDoneChan)

	// Initial connection handshake
	initId := 4715956011469373440

	for {
		buf := make([]byte, 2048)
		rlen, addr, err := l.ReadFromUDP(buf)
		Static.LogChan <- fmt.Sprintf("len: %d, addr: %s, err: %s", rlen, addr, err)

		// Verify length is at least 16 bytes
		if rlen < 16 {
			Static.LogChan <- "Invalid length"
			continue
		}

		connId := binary.BigEndian.Uint64(buf[0:8])
		action := binary.BigEndian.Uint32(buf[8:12])
		transId := binary.BigEndian.Uint32(buf[12:16])

		Static.LogChan <- fmt.Sprintf("id: %d, action: %d, trans: %d", connId, action, transId)

		// Verify valid connection ID
		_ = initId
		/*
		if connId != initId {
			Static.LogChan <- "Invalid connection ID"
			continue
		}
		*/

		// Action switch
		switch action {
		// Connect
		case 0:
			Static.LogChan <- "connect()"
			res := bytes.NewBuffer(make([]byte, 0))

			// Action
			binary.Write(res, binary.BigEndian, uint32(0))
			// Transaction ID
			binary.Write(res, binary.BigEndian, uint32(transId))
			// Connection ID
			binary.Write(res, binary.BigEndian, uint64(1234))

			rlen, err := l.WriteToUDP(res.Bytes(), addr)
			if err != nil {
				Static.LogChan <- err.Error()
				continue
			}

			Static.LogChan <- fmt.Sprintf("udp: Wrote %d bytes, %s", rlen, hex.EncodeToString(res.Bytes()))
			continue
		// Announce
		case 1:
			Static.LogChan <- "announce()"

			query := map[string]string{}

			// Connection ID
			connId := binary.BigEndian.Uint64(buf[0:8])

			// Action
			action := binary.BigEndian.Uint32(buf[8:12])

			// Transaction ID
			transId := binary.BigEndian.Uint32(buf[12:16])

			// Info hash
			query["info_hash"] = string(buf[16:36])

			// Peer ID
			query["peer_id"] = string(buf[36:56])

			// Downloaded
			t, _ := strconv.ParseInt(hex.EncodeToString(buf[56:64]), 16, 64)
			query["downloaded"] = strconv.FormatInt(t, 10)
			Static.LogChan <- fmt.Sprintf("downloaded: %s", query["downloaded"])

			// Left
			t, _ = strconv.ParseInt(hex.EncodeToString(buf[64:72]), 16, 64)
			query["left"] = strconv.FormatInt(t, 10)
			Static.LogChan <- fmt.Sprintf("left: %s", query["left"])

			// Uploaded
			t, _ = strconv.ParseInt(hex.EncodeToString(buf[72:80]), 16, 64)
			query["uploaded"] = strconv.FormatInt(t, 10)
			Static.LogChan <- fmt.Sprintf("uploaded: %s", query["uploaded"])

			// Event
			t, _ = strconv.ParseInt(hex.EncodeToString(buf[80:84]), 16, 32)
			query["event"] = strconv.FormatInt(t, 10)

			// Convert event to actual string
			switch query["event"] {
			case "0":
				query["event"] = ""
			case "1":
				query["event"] = "completed"
			case "2":
				query["event"] = "started"
			case "3":
				query["event"] = "stopped"
			}

			Static.LogChan <- fmt.Sprintf("event: %s", query["event"])

			// IP address
			t, _ = strconv.ParseInt(hex.EncodeToString(buf[84:88]), 16, 32)
			query["ip"] = strconv.FormatInt(t, 10)

			// If no IP address set, use the UDP source
			if query["ip"] == "0" {
				query["ip"] = strings.Split(addr.String(), ":")[0]
			}

			Static.LogChan <- fmt.Sprintf("ip: %s", query["ip"])

			// Key
			query["key"] = hex.EncodeToString(buf[88:92])
			Static.LogChan <- fmt.Sprintf("key: %s", query["key"])

			// Numwant
			query["numwant"] = hex.EncodeToString(buf[92:96])

			// If numwant is hex max value, default to 50
			if query["numwant"] == "ffffffff" {
				query["numwant"] = "50"
			}

			Static.LogChan <- fmt.Sprintf("numwant: %s", query["numwant"])

			// Port
			t, _ = strconv.ParseInt(hex.EncodeToString(buf[96:98]), 16, 32)
			query["port"] = strconv.FormatInt(t, 10)
			Static.LogChan <- fmt.Sprintf("port: %s", query["port"])

			_, _, _ = connId, transId, action

			// Response buffer
			res := bytes.NewBuffer(make([]byte, 0))

			// Load matching file
			file := new(FileRecord).Load(query["info_hash"], "info_hash")

			// Action (1 for announce)
			binary.Write(res, binary.BigEndian, uint32(1))
			// Transaction ID
			binary.Write(res, binary.BigEndian, uint32(transId))
			// Interval
			binary.Write(res, binary.BigEndian, uint32(RandRange(Static.Config.Interval - 600, Static.Config.Interval)))
			// Leechers
			binary.Write(res, binary.BigEndian, uint32(file.Leechers()))
			// Seeders
			binary.Write(res, binary.BigEndian, uint32(file.Seeders()))
			// Peer list
			numwant, _ := strconv.Atoi(query["numwant"])
			binary.Write(res, binary.BigEndian, file.PeerList(query["ip"], numwant))

			rlen, err := l.WriteToUDP(res.Bytes(), addr)
			if err != nil {
				Static.LogChan <- err.Error()
				continue
			}

			Static.LogChan <- fmt.Sprintf("udp: Wrote %d bytes, %s", rlen, hex.EncodeToString(res.Bytes()))
		default:
			Static.LogChan <- "Invalid action"
			continue
		}
	}
}
