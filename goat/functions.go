package goat

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"net"
	"strconv"
	"time"
)

// b2ip converts a packed byte buffer into a compactPeer
func b2ip(buf []byte) compactPeer {
	// Collect peer info
	peer := compactPeer{}

	// IP address
	peer.IP = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()

	// Port
	peer.Port = binary.BigEndian.Uint16(buf[4:6])

	return peer
}

// ip2b converts an IP and port into byte buffer
func ip2b(ip string, port uint16) []byte {
	// Empty buffers
	ipBuf, portBuf := [4]byte{}, [2]byte{}

	// Write IP
	binary.BigEndian.PutUint32(ipBuf[:], binary.BigEndian.Uint32(net.ParseIP(ip).To4()))

	// Write port
	binary.BigEndian.PutUint16(portBuf[:], port)

	// Concatenate buffers
	return append(ipBuf[:], portBuf[:]...)
}

// randRange generates a random announce interval in the specified range
func randRange(min int, max int) int {
	mrand.Seed(time.Now().Unix())
	return min + mrand.Intn(max-min)
}

// randString generates a random hex string, used for generating hashes
// Thanks: http://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
func randString() string {
	u := make([]byte, 16)
	_, err := crand.Read(u)
	if err != nil {
		// On failure, get a random number
		return strconv.Itoa(randRange(0, 1000000))
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}
