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
