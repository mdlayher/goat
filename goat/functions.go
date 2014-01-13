package goat

import (
	"encoding/binary"
	"math/rand"
	"net"
	"time"
)

// Parse IP and port into byte buffer
func ip2b(ip_ string, port_ uint16) []byte {
	// Empty buffers
	ip, port := [4]byte{}, [2]byte{}

	// Write IP
	binary.BigEndian.PutUint32(ip[:], binary.BigEndian.Uint32(net.ParseIP(ip_).To4()))

	// Write port
	binary.BigEndian.PutUint16(port[:], port_)

	// Concatenate buffers
	return append(ip[:], port[:]...)
}

// RandRange generates a random announce interval in the specified range
func randRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
