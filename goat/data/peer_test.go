package data

import (
	"log"
	"testing"
)

// TestPeer verifies that Peer binary marshal and unmarshal work properly
func TestPeer(t *testing.T) {
	log.Println("TestPeer()")

	// Generate mock peer
	peer := Peer{
		IP:   "127.0.0.1",
		Port: 8080,
	}

	// Marshal to binary representation
	out, err := peer.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal peer to binary: %s", err.Error())
	}

	// Unmarshal peer from binary representation
	peer2 := new(Peer)
	if err := peer2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal peer from binary: %s", err.Error())
	}

	// Verify peers are identical
	if peer.IP != peer2.IP || peer.Port != peer2.Port {
		t.Fatalf("Peer results do not match")
	}
}
