package udp

import (
	"bytes"
	"log"
	"testing"
)

// TestAnnounceRequest verifies that AnnounceRequest binary marshal and unmarshal work properly
func TestAnnounceRequest(t *testing.T) {
	log.Println("TestAnnounceRequest()")

	// Generate mock AnnounceRequest
	announce := AnnounceRequest{
		ConnID: 1,
		Action: 1,
		TransID: 1,
		InfoHash: []byte("abcdef0123456789abcd"),
		PeerID: []byte("abcdef0123456789abcd"),
		Downloaded: 0,
		Left: 0,
		Uploaded: 0,
		Event: 0,
		IP: 1234,
		Key: 1234,
		Numwant: 50,
		Port: 8080,
	}

	// Marshal to binary representation
	out, err := announce.MarshalBinary()
	if err != nil {
		t.Fatalf("Failed to marshal AnnounceRequest to binary: %s", err.Error())
	}

	// Unmarshal announce from binary representation
	announce2 := new(AnnounceRequest)
	if err := announce2.UnmarshalBinary(out); err != nil {
		t.Fatalf("Failed to unmarshal AnnounceRequest from binary: %s", err.Error())
	}

	// Verify announces are identical
	if announce.ConnID != announce2.ConnID || !bytes.Equal(announce.InfoHash, announce2.InfoHash) {
		t.Fatalf("AnnounceRequest results do not match")
	}
}
