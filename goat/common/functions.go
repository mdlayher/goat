package common

import (
	crand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	mrand "math/rand"
	"net"
	"strconv"
	"time"
)

// B2IP converts a packed byte buffer into an IP and port
func B2IP(buf []byte) (string, uint16) {
	// IP address
	ip := net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()

	// Port
	port := binary.BigEndian.Uint16(buf[4:6])

	return ip, port
}

// IP2B converts an IP and port into byte buffer
func IP2B(ip string, port uint16) []byte {
	// Empty buffers
	ipBuf, portBuf := [4]byte{}, [2]byte{}

	// Write IP
	binary.BigEndian.PutUint32(ipBuf[:], binary.BigEndian.Uint32(net.ParseIP(ip).To4()))

	// Write port
	binary.BigEndian.PutUint16(portBuf[:], port)

	// Concatenate buffers
	return append(ipBuf[:], portBuf[:]...)
}

// RandRange generates a random announce interval in the specified range
func RandRange(min int, max int) int {
	mrand.Seed(time.Now().Unix())
	return min + mrand.Intn(max-min)
}

// RandString generates a random hex string, used for generating hashes
// Thanks: http://stackoverflow.com/questions/15130321/is-there-a-method-to-generate-a-uuid-with-go-language
func RandString() string {
	u := make([]byte, 16)
	_, err := crand.Read(u)
	if err != nil {
		// On failure, get a random number
		return strconv.Itoa(RandRange(0, 1000000))
	}

	u[8] = (u[8] | 0x80) & 0xBF
	u[6] = (u[6] | 0x40) & 0x4F

	return hex.EncodeToString(u)
}
