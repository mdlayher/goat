package goat

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

// Parse IP and port into byte buffer
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

// RandRange generates a random announce interval in the specified range
func randRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
