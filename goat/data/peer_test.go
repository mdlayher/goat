package data

import (
	"log"
	"testing"

	"github.com/mdlayher/goat/goat/common"
)

// TestPeer verifies that Peer creation, methods, save, load, and delete work properly
func TestPeer(t *testing.T) {
	log.Println("TestPeer()")

	// Load config
	config, err := common.LoadConfig()
	if err != nil {
		t.Fatalf("Could not load configuration: %s", err.Error())
	}
	common.Static.Config = config

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
