package goat

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"strings"
)

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

			// Mark client as UDP
			query["udp"] = "1"

			// Connection ID
			connId := binary.BigEndian.Uint64(buf[0:8])

			// Action
			action := binary.BigEndian.Uint32(buf[8:12])

			// Transaction ID
			transId := binary.BigEndian.Uint32(buf[12:16])
			transIdBuf := buf[12:16]

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

			// TODO: temporary, load user
			user := new(UserRecord).Load(1, "id")

			// Trigger an announce
			resChan := make(chan []byte)
			go TrackerAnnounce(user, query, true, transIdBuf, resChan)

			rlen, err := l.WriteToUDP(<-resChan, addr)
			if err != nil {
				Static.LogChan <- err.Error()
				continue
			}

			Static.LogChan <- fmt.Sprintf("udp: Wrote %d bytes", rlen)
		default:
			Static.LogChan <- "Invalid action"
			continue
		}
	}
}
