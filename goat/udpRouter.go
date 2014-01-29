package goat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

// Handshake for UDP tracker protocol
const udpInitID = 4497486125440

// UDP errors
var (
	// errUDPAction is returned when a client requests an invalid tracker action
	errUDPAction = errors.New("udp: client did not send a valid UDP tracker action")
	// errUDPHandshake is returned when a client does not send the proper handshake ID
	errUDPHandshake = errors.New("udp: client did not send proper UDP tracker handshake")
	// errUDPInteger is returned when a client sends an invalid integer parameter
	errUDPInteger = errors.New("udp: client sent an invalid integer parameter")
	// errUDPWrite is returned when the tracker cannot generate a proper response
	errUDPWrite = errors.New("udp: tracker cannot generate UDP tracker response")
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
		atomic.AddInt64(&static.UDP.Current, 1)
		atomic.AddInt64(&static.UDP.Total, 1)

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

// Parse a UDP byte buffer, return response from tracker
func parseUDP(buf []byte, addr *net.UDPAddr) ([]byte, error) {
	// Attempt to create UDP packet from buffer
	packet, err := new(udpPacket).FromBytes(buf)
	if err != nil {
		// Because no transaction ID is present on failure, we must return nil
		return nil, errUDPHandshake
	}

	// Create a udpTracker to handle this client
	tracker := udpTracker{TransID: packet.TransID}

	// Action switch
	// Action 0: Connect
	if packet.Action == 0 {
		// Validate UDP tracker handshake
		if packet.ConnID != udpInitID {
			return tracker.Error("Invalid UDP tracker handshake"), errUDPHandshake
		}

		res := bytes.NewBuffer(make([]byte, 0))

		// Action
		err := binary.Write(res, binary.BigEndian, uint32(0))
		if err != nil {
			log.Println(err.Error())
			return tracker.Error("Could not generate UDP tracker response"), errUDPWrite
		}

		// Transaction ID
		err = binary.Write(res, binary.BigEndian, packet.TransID)
		if err != nil {
			log.Println(err.Error())
			return tracker.Error("Could not generate UDP tracker response"), errUDPWrite
		}

		// Generate a connection ID, which will be expected for this client next call
		expID := uint64(randRange(1, 1000000000))

		// Store this client's address and ID in map
		udpAddrToID[addr.String()] = expID

		// Connection ID, generated for this session
		err = binary.Write(res, binary.BigEndian, expID)
		if err != nil {
			log.Println(err.Error())
			return tracker.Error("Could not generate UDP tracker response"), errUDPWrite
		}

		return res.Bytes(), nil
	}

	// For all tracker actions other than connect, we must validate the connection ID for this
	// address, ensuring it matches the previously set value

	// Ensure connection ID map contains this IP address
	expID, ok := udpAddrToID[addr.String()]
	if !ok {
		return tracker.Error("Client must properly handshake before announce"), errUDPHandshake
	}

	// Validate expected connection ID using map
	if packet.ConnID != expID {
		return tracker.Error("Invalid UDP connection ID"), errUDPHandshake
	}

	// Clear this IP from the connection map after 2 minutes
	// note: this is done to conserve memory and prevent session fixation
	go func(addr *net.UDPAddr) {
		<-time.After(2 * time.Minute)
		delete(udpAddrToID, addr.String())
	}(addr)

	// Action 1: Announce
	if packet.Action == 1 {
		// Generate UDP announce packet from byte buffer
		announce, err := new(udpAnnouncePacket).FromBytes(buf)
		if err != nil {
			return tracker.Error("Malformed UDP announce"), errUDPInteger
		}

		// Convert UDP announce to query map
		query := announce.ToValues()

		// Check if a proper IP was set, and if not, use the UDP connection address
		if query.Get("ip") == "0" {
			query.Set("ip", strings.Split(addr.String(), ":")[0])
		}

		// Trigger an anonymous announce
		return trackerAnnounce(tracker, userRecord{}, query), nil
	}

	// Action 2: Scrape
	if packet.Action == 2 {
		// Generate UDP scrape packet from byte buffer
		scrape, err := new(udpScrapePacket).FromBytes(buf)
		if err != nil {
			return tracker.Error("Malformed UDP scrape"), errUDPHandshake
		}

		// Convert UDP scrape to query map
		query := scrape.ToValues()

		// Store IP in query map
		query.Set("ip", strings.Split(addr.String(), ":")[0])

		// Trigger a scrape
		return trackerScrape(tracker, query), nil
	}

	// No action matched
	return tracker.Error("Invalid action"), errUDPAction
}
