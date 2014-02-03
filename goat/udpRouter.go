package goat

import (
	"errors"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
	"github.com/mdlayher/goat/goat/data/udp"
	"github.com/mdlayher/goat/goat/tracker"
)

// Handshake for UDP udpTracker protocol
const udpInitID = 4497486125440

// UDP errors
var (
	// errUDPAction is returned when a client requests an invalid udpTracker action
	errUDPAction = errors.New("udp: client did not send a valid UDP udpTracker action")
	// errUDPHandshake is returned when a client does not send the proper handshake ID
	errUDPHandshake = errors.New("udp: client did not send proper UDP udpTracker handshake")
	// errUDPInteger is returned when a client sends an invalid integer parameter
	errUDPInteger = errors.New("udp: client sent an invalid integer parameter")
	// errUDPWrite is returned when the udpTracker cannot generate a proper response
	errUDPWrite = errors.New("udp: udpTracker cannot generate UDP udpTracker response")
)

// UDP address to connection ID map
var udpAddrToID = map[string]uint64{}

// Handle incoming UDP connections and return response
func handleUDP(l *net.UDPConn, sendChan chan bool, recvChan chan bool) {
	// Create shutdown function
	go func(l *net.UDPConn, sendChan chan bool, recvChan chan bool) {
		// Wait for done signal
		<-sendChan

		// Close listener
		if err := l.Close(); err != nil {
			log.Println(err.Error())
		}

		log.Println("UDP listener stopped")
		recvChan <- true
	}(l, sendChan, recvChan)

	// Loop and read connections
	for {
		buf := make([]byte, 2048)
		rlen, addr, err := l.ReadFromUDP(buf)

		// Count incoming connections
		atomic.AddInt64(&common.Static.UDP.Current, 1)
		atomic.AddInt64(&common.Static.UDP.Total, 1)

		// Triggered on graceful shutdown
		if err != nil {
			// Ignore connection closing error, caused by stopping network listener
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Println(err.Error())
				panic(err)
			}

			return
		}

		// Verify length is at least 16 bytes
		if rlen < 16 {
			log.Println("Invalid length")
			continue
		}

		// Spawn a goroutine to handle the connection and send back the response
		go func(l *net.UDPConn, buf []byte, addr *net.UDPAddr) {
			// Capture initial response from buffer
			res, err := parseUDP(buf, addr)
			if err != nil {
				// Client sent a malformed UDP handshake
				log.Println(err.Error())

				// If error, client did not handshake correctly, so boot them with error message
				_, err2 := l.WriteToUDP(res, addr)
				if err2 != nil {
					log.Println(err2.Error())
				}

				return
			}

			// Write response
			_, err = l.WriteToUDP(res, addr)
			if err != nil {
				log.Println(err.Error())
			}

			return
		}(l, buf, addr)
	}
}

// Parse a UDP byte buffer, return response from udpTracker
func parseUDP(buf []byte, addr *net.UDPAddr) ([]byte, error) {
	// Attempt to grab generic UDP connection fields
	packet := new(udp.Packet)
	err := packet.UnmarshalBinary(buf)
	if err != nil {
		// Because no transaction ID is present on failure, we must return nil
		return nil, errUDPHandshake
	}

	// Create a udpTracker to handle this client
	udpTracker := tracker.UDPTracker{TransID: packet.TransID}

	// Check for maintenance mode
	if common.Static.Maintenance {
		// Return tracker error with maintenance message
		return udpTracker.Error("Maintenance: " + common.Static.StatusMessage), nil
	}

	// Action switch
	// Action 0: Connect
	if packet.Action == 0 {
		// Validate UDP udpTracker handshake
		if packet.ConnID != udpInitID {
			return udpTracker.Error("Invalid UDP udpTracker handshake"), errUDPHandshake
		}

		// Generate a connection ID, which will be expected for this client next call
		expID := uint64(common.RandRange(1, 1000000000))

		// Store this client's address and ID in map
		udpAddrToID[addr.String()] = expID

		// Generate connect response
		connect := udp.ConnectResponse{
			Action:  0,
			TransID: packet.TransID,
			ConnID:  expID,
		}

		// Grab bytes from connect response
		connectBuf, err := connect.MarshalBinary()
		if err != nil {
			log.Println(err.Error())
			return udpTracker.Error("Could not generate UDP connect response"), errUDPWrite
		}

		return connectBuf, nil
	}

	// For all udpTracker actions other than connect, we must validate the connection ID for this
	// address, ensuring it matches the previously set value

	// Ensure connection ID map contains this IP address
	expID, ok := udpAddrToID[addr.String()]
	if !ok {
		return udpTracker.Error("Client must properly handshake before announce"), errUDPHandshake
	}

	// Validate expected connection ID using map
	if packet.ConnID != expID {
		return udpTracker.Error("Invalid UDP connection ID"), errUDPHandshake
	}

	// Clear this IP from the connection map after 2 minutes
	// note: this is done to conserve memory and prevent session fixation
	go func(addr *net.UDPAddr) {
		<-time.After(2 * time.Minute)
		delete(udpAddrToID, addr.String())
	}(addr)

	// Action 1: Announce
	if packet.Action == 1 {
		// Retrieve UDP announce request from byte buffer
		announce := new(udp.AnnounceRequest)
		err := announce.UnmarshalBinary(buf)
		log.Println(announce)
		if err != nil {
			return udpTracker.Error("Malformed UDP announce"), errUDPInteger
		}

		// Convert UDP announce to query map
		query := announce.ToValues()

		// Check if a proper IP was set, and if not, use the UDP connection address
		if query.Get("ip") == "0" {
			query.Set("ip", strings.Split(addr.String(), ":")[0])
		}

		// Trigger an anonymous announce
		return tracker.Announce(udpTracker, data.UserRecord{}, query), nil
	}

	// Action 2: Scrape
	if packet.Action == 2 {
		// Generate UDP scrape packet from byte buffer
		scrape := new(udp.ScrapeRequest)
		err := scrape.UnmarshalBinary(buf)
		if err != nil {
			return udpTracker.Error("Malformed UDP scrape"), errUDPHandshake
		}

		// Convert UDP scrape to query map
		query := scrape.ToValues()

		// Store IP in query map
		query.Set("ip", strings.Split(addr.String(), ":")[0])

		// Trigger a scrape
		return tracker.Scrape(udpTracker, query), nil
	}

	// No action matched
	return udpTracker.Error("Invalid action"), errUDPAction
}
